package users

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IUserService interface {
	common.IBaseCrudService[entities.User, UserInput, UserPage]
	Login(c fiber.Ctx, input LoginInput) (LoginResponse, error)
	RegenerateTokenPair(c fiber.Ctx, userID any, isRemember bool) (LoginResponse, error)
}

type UserService struct {
	repo IUserRepo
}

var _ IUserService = &UserService{}

func NewUserService(repo IUserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(c common.Context, input UserInput) (uint, error) {
	password, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	input.Password = password

	return s.repo.Create(c, input)
}

func (s *UserService) Read(c fiber.Ctx) (paginate.Page, []UserPage) {
	return s.repo.Read(c)
}

func (s *UserService) Update(c common.Context, id any, input UserInput) error {
	password, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	input.Password = password

	return s.repo.Update(c, id, input)
}

func (s *UserService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *UserService) GetByID(c common.Context, id any) (entities.User, error) {
	return s.repo.GetByID(c, id)
}

func (s *UserService) Login(c fiber.Ctx, input LoginInput) (LoginResponse, error) {
	var res LoginResponse
	user, err := s.repo.GetByUsernameOrEmail(c, input.Username)
	if err != nil {
		return res, err
	}

	if ok, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password); !ok {
		return res, fiber.ErrUnauthorized
	}

	res, err = s.generateTokenPair(user, input.IsRemember, c.BaseURL())
	if err != nil {
		return res, fmt.Errorf("failed to generate token pair: %w", err)
	}
	common.GenerateCookies(c, res.AccessToken, res.RefreshToken, input.IsRemember)

	return res, nil
}

func (s *UserService) RegenerateTokenPair(c fiber.Ctx, userID any, isRemember bool) (LoginResponse, error) {
	var res LoginResponse
	user, err := s.repo.GetByID(c, userID)
	if err != nil {
		return res, err
	}

	res, err = s.generateTokenPair(user, isRemember, c.BaseURL())
	if err != nil {
		return res, fmt.Errorf("failed to generate token pair: %w", err)
	}
	common.GenerateCookies(c, res.AccessToken, res.RefreshToken, isRemember)

	return res, nil
}

func (s *UserService) generateTokenPair(user entities.User, isRemember bool, iss string) (LoginResponse, error) {
	var res LoginResponse
	claims := map[string]any{
		"sub":    user.Username,
		"userId": user.ID,
		"iss":    iss,
		"aud":    iss,
	}
	at, rt, err := common.GenerateJWTPair(claims, nil, isRemember)
	if err != nil {
		return res, fmt.Errorf("failed generating token pair: %w", err)
	}

	res.AccessToken = at
	res.RefreshToken = rt
	res.User = user

	return res, nil
}
