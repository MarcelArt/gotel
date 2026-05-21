package user_roles

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IUserRoleService interface {
	common.IBaseCrudService[entities.UserRole, UserRoleInput, UserRolePage]
}

type UserRoleService struct {
	repo IUserRoleRepo
}

var _ IUserRoleService = &UserRoleService{}

func NewUserRoleService(repo IUserRoleRepo) *UserRoleService {
	return &UserRoleService{
		repo: repo,
	}
}

func (s *UserRoleService) Create(c common.Context, input UserRoleInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *UserRoleService) Read(c fiber.Ctx) (paginate.Page, []UserRolePage) {
	return s.repo.Read(c)
}

func (s *UserRoleService) Update(c common.Context, id any, input UserRoleInput) error {
	return s.repo.Update(c, id, input)
}

func (s *UserRoleService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *UserRoleService) GetByID(c common.Context, id any) (entities.UserRole, error) {
	return s.repo.GetByID(c, id)
}
