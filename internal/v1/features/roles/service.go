package roles

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IRoleService interface {
	common.IBaseCrudService[entities.Role, RoleInput, RolePage]
}

type RoleService struct {
	repo IRoleRepo
}

var _ IRoleService = &RoleService{}

func NewRoleService(repo IRoleRepo) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (s *RoleService) Create(c common.Context, input RoleInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *RoleService) Read(c fiber.Ctx) (paginate.Page, []RolePage) {
	return s.repo.Read(c)
}

func (s *RoleService) Update(c common.Context, id any, input RoleInput) error {
	return s.repo.Update(c, id, input)
}

func (s *RoleService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *RoleService) GetByID(c common.Context, id any) (entities.Role, error) {
	return s.repo.GetByID(c, id)
}
