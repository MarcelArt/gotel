package roles

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IRoleRepo interface {
	common.IBaseCrudRepo[entities.Role, RoleInput, RolePage]
}

type RoleRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IRoleRepo = &RoleRepo{}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		db:        db,
		pageQuery: "SELECT id, name, description, permissions FROM roles",
	}
}

func (r *RoleRepo) Create(c common.Context, input RoleInput) (uint, error) {
	ctx := c.Context()

	role, err := common.Cast[entities.Role](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.Role](r.db).Create(ctx, &role)

	return role.ID, err
}

func (r *RoleRepo) Read(c fiber.Ctx) (paginate.Page, []RolePage) {
	roles := make([]RolePage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&roles)

	return page, roles
}

func (r *RoleRepo) Update(c common.Context, id any, input RoleInput) error {
	ctx := c.Context()
	role, err := common.Cast[entities.Role](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.Role](r.db).Where("id = ?", id).Updates(ctx, role)

	return err
}

func (r *RoleRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.Role](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *RoleRepo) GetByID(c common.Context, id any) (entities.Role, error) {
	var role entities.Role
	ctx := c.Context()

	role, err := gorm.G[entities.Role](r.db).Where("id = ?", id).First(ctx)

	return role, err
}
