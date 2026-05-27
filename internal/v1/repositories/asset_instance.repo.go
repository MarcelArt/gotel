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

type IAssetInstanceRepo interface {
	common.IBaseCrudRepo[entities.AssetInstance, models.AssetInstanceInput, models.AssetInstancePage]
}

type AssetInstanceRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IAssetInstanceRepo = &AssetInstanceRepo{}

func NewAssetInstanceRepo(db *gorm.DB) *AssetInstanceRepo {
	return &AssetInstanceRepo{
		db:        db,
		pageQuery: "SELECT id, code, item_id FROM asset_instances where deleted_at isnull",
	}
}

func (r *AssetInstanceRepo) Create(c common.Context, input models.AssetInstanceInput) (uint, error) {
	ctx := c.Context()

	instance, err := common.Cast[entities.AssetInstance](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.AssetInstance](r.db).Create(ctx, &instance)

	return instance.ID, err
}

func (r *AssetInstanceRepo) Read(c fiber.Ctx) (paginate.Page, []models.AssetInstancePage) {
	instances := make([]models.AssetInstancePage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&instances)

	return page, instances
}

func (r *AssetInstanceRepo) Update(c common.Context, id any, input models.AssetInstanceInput) error {
	ctx := c.Context()
	instance, err := common.Cast[entities.AssetInstance](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.AssetInstance](r.db).Where("id = ?", id).Updates(ctx, instance)

	return err
}

func (r *AssetInstanceRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.AssetInstance](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *AssetInstanceRepo) GetByID(c common.Context, id any) (entities.AssetInstance, error) {
	var instance entities.AssetInstance
	ctx := c.Context()

	instance, err := gorm.G[entities.AssetInstance](r.db).Where("id = ?", id).First(ctx)

	return instance, err
}
