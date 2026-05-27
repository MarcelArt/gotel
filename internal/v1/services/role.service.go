package services

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IRoleService interface {
	common.IBaseCrudService[entities.Role, models.RoleInput, models.RolePage]
}

type RoleService struct {
	repo repositories.IRoleRepo
}

var _ IRoleService = &RoleService{}

func NewRoleService(repo repositories.IRoleRepo) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (s *RoleService) Create(c common.Context, input models.RoleInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *RoleService) Read(c fiber.Ctx) (paginate.Page, []models.RolePage) {
	return s.repo.Read(c)
}

func (s *RoleService) Update(c common.Context, id any, input models.RoleInput) error {
	return s.repo.Update(c, id, input)
}

func (s *RoleService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *RoleService) GetByID(c common.Context, id any) (entities.Role, error) {
	return s.repo.GetByID(c, id)
}
