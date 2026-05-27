package services

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IUserRoleService interface {
	common.IBaseCrudService[entities.UserRole, models.UserRoleInput, models.UserRolePage]
}

type UserRoleService struct {
	repo repositories.IUserRoleRepo
}

var _ IUserRoleService = &UserRoleService{}

func NewUserRoleService(repo repositories.IUserRoleRepo) *UserRoleService {
	return &UserRoleService{
		repo: repo,
	}
}

func (s *UserRoleService) Create(c common.Context, input models.UserRoleInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *UserRoleService) Read(c fiber.Ctx) (paginate.Page, []models.UserRolePage) {
	return s.repo.Read(c)
}

func (s *UserRoleService) Update(c common.Context, id any, input models.UserRoleInput) error {
	return s.repo.Update(c, id, input)
}

func (s *UserRoleService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *UserRoleService) GetByID(c common.Context, id any) (entities.UserRole, error) {
	return s.repo.GetByID(c, id)
}
