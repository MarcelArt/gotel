package repositories

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IUserRoleRepo interface {
	common.IBaseCrudRepo[entities.UserRole, models.UserRoleInput, models.UserRolePage]
	GetRoleIDsByUserID(userID any) ([]uint, error)
	DeleteByUserIDAndRoleIDs(c common.Context, userID any, roleIDs []uint) error
	BulkCreate(c common.Context, input []models.UserRoleInput) error
}

type UserRoleRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IUserRoleRepo = &UserRoleRepo{}

func NewUserRoleRepo(db *gorm.DB) *UserRoleRepo {
	return &UserRoleRepo{
		db: db,
		pageQuery: `
			SELECT 
				ur.*,
				u.username as username,
				r."name" as role
			from user_roles ur 
			join users u on ur.user_id = u.id 
			join roles r on ur.role_id = r.id 
			where ur.deleted_at isnull
		`,
	}
}

func (r *UserRoleRepo) Create(c common.Context, input models.UserRoleInput) (uint, error) {
	ctx := c.Context()

	userRole, err := common.Cast[entities.UserRole](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.UserRole](r.db).Create(ctx, &userRole)

	return userRole.ID, err
}

func (r *UserRoleRepo) Read(c fiber.Ctx) (paginate.Page, []models.UserRolePage) {
	userRoles := make([]models.UserRolePage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&userRoles)

	return page, userRoles
}

func (r *UserRoleRepo) Update(c common.Context, id any, input models.UserRoleInput) error {
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

func (r *UserRoleRepo) GetRoleIDsByUserID(userID any) ([]uint, error) {
	var roleIDs []uint
	err := r.db.Model(entities.UserRole{}).Where("user_id = ?", userID).Distinct("role_id").Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

func (r *UserRoleRepo) DeleteByUserIDAndRoleIDs(c common.Context, userID any, roleIDs []uint) error {
	ctx := c.Context()
	_, err := gorm.G[entities.UserRole](r.db).Where("user_id = ? and role_id in ?", userID, roleIDs).Delete(ctx)

	return err
}

func (r *UserRoleRepo) BulkCreate(c common.Context, input []models.UserRoleInput) error {
	ctx := c.Context()

	userRoles, err := common.Cast[[]entities.UserRole](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.UserRole](r.db).CreateInBatches(ctx, &userRoles, 100)

	return err
}
