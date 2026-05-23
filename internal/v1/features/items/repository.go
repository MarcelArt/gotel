package items

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IItemRepo interface {
	common.IBaseCrudRepo[entities.Item, ItemInput, ItemPage]
}

type ItemRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IItemRepo = &ItemRepo{}

func NewItemRepo(db *gorm.DB) *ItemRepo {
	return &ItemRepo{
		db:        db,
		pageQuery: "SELECT id, code, name, picture, tracking_mode, unit, category_id FROM items",
	}
}

func (r *ItemRepo) Create(c common.Context, input ItemInput) (uint, error) {
	ctx := c.Context()

	item, err := common.Cast[entities.Item](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.Item](r.db).Create(ctx, &item)

	return item.ID, err
}

func (r *ItemRepo) Read(c fiber.Ctx) (paginate.Page, []ItemPage) {
	items := make([]ItemPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&items)

	return page, items
}

func (r *ItemRepo) Update(c common.Context, id any, input ItemInput) error {
	ctx := c.Context()
	item, err := common.Cast[entities.Item](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.Item](r.db).Where("id = ?", id).Updates(ctx, item)

	return err
}

func (r *ItemRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.Item](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *ItemRepo) GetByID(c common.Context, id any) (entities.Item, error) {
	var item entities.Item
	ctx := c.Context()

	item, err := gorm.G[entities.Item](r.db).Where("id = ?", id).First(ctx)

	return item, err
}
