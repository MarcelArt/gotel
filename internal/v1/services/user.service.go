package services

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/MarcelArt/gotel/internal/v1/usecases"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IUserService interface {
	common.IBaseCrudService[entities.User, models.UserInput, models.UserPage]
	Login(c fiber.Ctx, input models.LoginInput) (models.LoginResponse, error)
	RegenerateTokenPair(c fiber.Ctx, userID any, isRemember bool) (models.LoginResponse, error)
	AssignRoles(c common.Context, userID uint, newRoleIDs []uint) error
	GetPermissions(userID any) ([]string, error)
	GetRoles(id any) ([]models.UserRole, error)
}

type UserService struct {
	db     *gorm.DB
	repo   repositories.IUserRepo
	urRepo repositories.IUserRoleRepo
}

var _ IUserService = &UserService{}

func NewUserService(db *gorm.DB, repo repositories.IUserRepo, urRepo repositories.IUserRoleRepo) *UserService {
	return &UserService{
		db:     db,
		repo:   repo,
		urRepo: urRepo,
	}
}

func (s *UserService) Create(c common.Context, input models.UserInput) (uint, error) {
	password, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	input.Password = password

	return s.repo.Create(c, input)
}

func (s *UserService) Read(c fiber.Ctx) (paginate.Page, []models.UserPage) {
	return s.repo.Read(c)
}

func (s *UserService) Update(c common.Context, id any, input models.UserInput) error {
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

func (s *UserService) Login(c fiber.Ctx, input models.LoginInput) (models.LoginResponse, error) {
	var res models.LoginResponse
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

func (s *UserService) RegenerateTokenPair(c fiber.Ctx, userID any, isRemember bool) (models.LoginResponse, error) {
	var res models.LoginResponse
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
	tx := s.db.Begin()
	defer tx.Rollback()

	usecase := usecases.InitAssignRolesUsecase(tx)
	usecase.UserID = userID
	usecase.NewRoleIDs = newRoleIDs

	if err := usecase.Execute(c); err != nil {
		return err
	}

	return tx.Commit().Error
}

func (s *UserService) GetPermissions(userID any) ([]string, error) {
	return s.repo.GetPermissions(userID)
}

func (s *UserService) generateTokenPair(user entities.User, isRemember bool, iss string) (models.LoginResponse, error) {
	var res models.LoginResponse
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

func (s *UserService) GetRoles(id any) ([]models.UserRole, error) {
	return s.repo.GetRoles(id)
}
