package users

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/configs"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/shared"
	"github.com/MarcelArt/gotel/pkg/arrays"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IUserService interface {
	common.IBaseCrudService[entities.User, UserInput, UserPage]
	Login(c fiber.Ctx, input LoginInput) (LoginResponse, error)
	RegenerateTokenPair(c fiber.Ctx, userID any, isRemember bool) (LoginResponse, error)
	AssignRoles(c common.Context, userID uint, newRoleIDs []uint) error
	GetPermissions(userID any) ([]string, error)
	GetRoles(id any) ([]UserRole, error)
}

type UserService struct {
	repo   IUserRepo
	urRepo shared.IUserRoleRepoTx
}

var _ IUserService = &UserService{}

func NewUserService(repo IUserRepo, urRepo shared.IUserRoleRepoTx) *UserService {
	return &UserService{
		repo:   repo,
		urRepo: urRepo,
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
	if input.Password != "" {
		password, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		input.Password = password
	}

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

func (s *UserService) AssignRoles(c common.Context, userID uint, newRoleIDs []uint) error {
	oldRoles, err := s.urRepo.GetRoleIDsByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to retrieve old roles: %w", err)
	}

	idsToRemove, idsToAdd := arrays.DiffCheck(oldRoles, newRoleIDs)

	tx := configs.DB.Begin()
	defer tx.Rollback()

	urRepo := s.urRepo.BeginTx(tx)

	if len(idsToRemove) > 0 {
		if err := urRepo.DeleteByUserIDAndRoleIDs(c, userID, idsToRemove); err != nil {
			return fmt.Errorf("failed to delete old roles: %w", err)
		}
	}

	if len(idsToAdd) > 0 {
		userRoles := arrays.Map(idsToAdd, func(roleID uint) shared.UserRoleInput {
			return shared.UserRoleInput{
				UserID: userID,
				RoleID: roleID,
			}
		})

		if err := urRepo.BulkCreate(c, userRoles); err != nil {
			return fmt.Errorf("failed to add new roles: %w", err)
		}
	}

	tx.Commit()
	return nil
}

func (s *UserService) GetPermissions(userID any) ([]string, error) {
	return s.repo.GetPermissions(userID)
}

func (s *UserService) generateTokenPair(user entities.User, isRemember bool, iss string) (LoginResponse, error) {
	var res LoginResponse
	claims := map[string]any{
		"sub":    user.Username,
		"userId": user.ID,
		"iss":    iss,
		"aud":    iss,
	}

	permissions, err := s.repo.GetPermissions(user.ID)
	if err != nil {
		return res, fmt.Errorf("failed retrieving permissions: %w", err)
	}

	at, rt, err := common.GenerateJWTPair(claims, permissions, isRemember)
	if err != nil {
		return res, fmt.Errorf("failed generating token pair: %w", err)
	}

	res.AccessToken = at
	res.RefreshToken = rt
	res.User = user

	return res, nil
}

func (s *UserService) GetRoles(id any) ([]UserRole, error) {
	return s.repo.GetRoles(id)
}
