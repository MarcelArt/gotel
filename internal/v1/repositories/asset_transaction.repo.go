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

type IAssetTransactionRepo interface {
	common.IBaseCrudRepo[entities.AssetTransaction, models.AssetTransactionInput, models.AssetTransactionPage]
}

type AssetTransactionRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IAssetTransactionRepo = &AssetTransactionRepo{}

func NewAssetTransactionRepo(db *gorm.DB) *AssetTransactionRepo {
	return &AssetTransactionRepo{
		db: db,
		pageQuery: `
			SELECT 
				t.*,
				l.value location,
				u.username actor
			FROM asset_transactions t
			join locations l on t.location_id = l.id
			join users u on t.actor_id = u.id
			where t.deleted_at isnull
		`,
	}
}

func (r *AssetTransactionRepo) Create(c common.Context, input models.AssetTransactionInput) (uint, error) {
	ctx := c.Context()

	tx, err := common.Cast[entities.AssetTransaction](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.AssetTransaction](r.db).Create(ctx, &tx)

	return tx.ID, err
}

func (r *AssetTransactionRepo) Read(c fiber.Ctx) (paginate.Page, []models.AssetTransactionPage) {
	txs := make([]models.AssetTransactionPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&txs)

	return page, txs
}

func (r *AssetTransactionRepo) Update(c common.Context, id any, input models.AssetTransactionInput) error {
	ctx := c.Context()
	tx, err := common.Cast[entities.AssetTransaction](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.AssetTransaction](r.db).Where("id = ?", id).Updates(ctx, tx)

	return err
}

func (r *AssetTransactionRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.AssetTransaction](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *AssetTransactionRepo) GetByID(c common.Context, id any) (entities.AssetTransaction, error) {
	var tx entities.AssetTransaction
	ctx := c.Context()

	tx, err := gorm.G[entities.AssetTransaction](r.db).Where("id = ?", id).First(ctx)

	return tx, err
}
