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
