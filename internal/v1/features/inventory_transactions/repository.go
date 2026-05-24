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
	GetItemCounts(itemID any) ([]ItemCount, error)
}

type InventoryTransactionRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IInventoryTransactionRepo = &InventoryTransactionRepo{}

func NewInventoryTransactionRepo(db *gorm.DB) *InventoryTransactionRepo {
	return &InventoryTransactionRepo{
		db: db,
		pageQuery: `
			select 
				it.id as id,
				it.created_at as created_at,
				it.transaction_type as transaction_type,
				it.quantity as quantity,
				it.note as note,
				it.item_id as item_id,
				i."name" as item,
				i.unit as unit,
				u.username as actor,
				lf.value as from,
				lt.value as to
			from inventory_transactions it 
			join items i on it.item_id = i.id 
			join users u on it.actor_id = u.id 
			left join locations lf on it.from_id = lf.id 
			left join locations lt on it.to_id = lt.id
			where it.deleted_at isnull
		`,
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

func (r *InventoryTransactionRepo) GetItemCounts(itemID any) ([]ItemCount, error) {
	itemCounts := make([]ItemCount, 0)
	query := `
		select
			it.transaction_type as transaction_type,
			SUM(it.quantity) as quantity
		from inventory_transactions it 
		where it.item_id = ?
		and it.transaction_type in ('RECEIVE', 'DISPOSE', 'CONSUME', 'LOST')
		and it.deleted_at isnull
		group by
			it.transaction_type 
	`

	err := r.db.Raw(query, itemID).Scan(&itemCounts).Error

	return itemCounts, err
}
