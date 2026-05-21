package user_roles

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IUserRoleRepo interface {
	common.IBaseCrudRepo[entities.UserRole, UserRoleInput, UserRolePage]
}

type UserRoleRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IUserRoleRepo = &UserRoleRepo{}

func NewUserRoleRepo(db *gorm.DB) *UserRoleRepo {
	return &UserRoleRepo{
		db:        db,
		pageQuery: "SELECT id, user_id, role_id FROM user_roles",
	}
}

func (r *UserRoleRepo) Create(c common.Context, input UserRoleInput) (uint, error) {
	ctx := c.Context()

	userRole, err := common.Cast[entities.UserRole](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.UserRole](r.db).Create(ctx, &userRole)

	return userRole.ID, err
}

func (r *UserRoleRepo) Read(c fiber.Ctx) (paginate.Page, []UserRolePage) {
	userRoles := make([]UserRolePage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&userRoles)

	return page, userRoles
}

func (r *UserRoleRepo) Update(c common.Context, id any, input UserRoleInput) error {
	ctx := c.Context()
	userRole, err := common.Cast[entities.UserRole](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.UserRole](r.db).Where("id = ?", id).Updates(ctx, userRole)

	return err
}

func (r *UserRoleRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.UserRole](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *UserRoleRepo) GetByID(c common.Context, id any) (entities.UserRole, error) {
	var userRole entities.UserRole
	ctx := c.Context()

	userRole, err := gorm.G[entities.UserRole](r.db).Where("id = ?", id).First(ctx)

	return userRole, err
}
