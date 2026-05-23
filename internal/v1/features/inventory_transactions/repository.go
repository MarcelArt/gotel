package inventory_transactions

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IInventoryTransactionRepo interface {
	common.IBaseCrudRepo[entities.InventoryTransaction, InventoryTransactionInput, InventoryTransactionPage]
}

type InventoryTransactionRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IInventoryTransactionRepo = &InventoryTransactionRepo{}

func NewInventoryTransactionRepo(db *gorm.DB) *InventoryTransactionRepo {
	return &InventoryTransactionRepo{
		db:        db,
		pageQuery: "SELECT id, transaction_type, quantity, note, item_id, from_id, to_id, actor_id FROM inventory_transactions",
	}
}

func (r *InventoryTransactionRepo) Create(c common.Context, input InventoryTransactionInput) (uint, error) {
	ctx := c.Context()

	tx, err := common.Cast[entities.InventoryTransaction](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.InventoryTransaction](r.db).Create(ctx, &tx)

	return tx.ID, err
}

func (r *InventoryTransactionRepo) Read(c fiber.Ctx) (paginate.Page, []InventoryTransactionPage) {
	txs := make([]InventoryTransactionPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&txs)

	return page, txs
}

func (r *InventoryTransactionRepo) Update(c common.Context, id any, input InventoryTransactionInput) error {
	ctx := c.Context()
	tx, err := common.Cast[entities.InventoryTransaction](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.InventoryTransaction](r.db).Where("id = ?", id).Updates(ctx, tx)

	return err
}

func (r *InventoryTransactionRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.InventoryTransaction](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *InventoryTransactionRepo) GetByID(c common.Context, id any) (entities.InventoryTransaction, error) {
	var tx entities.InventoryTransaction
	ctx := c.Context()

	tx, err := gorm.G[entities.InventoryTransaction](r.db).Where("id = ?", id).First(ctx)

	return tx, err
}
