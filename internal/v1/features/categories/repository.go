package categories

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type ICategoryRepo interface {
	common.IBaseCrudRepo[entities.Category, CategoryInput, CategoryPage]
}

type CategoryRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ ICategoryRepo = &CategoryRepo{}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{
		db:        db,
		pageQuery: "SELECT id, value, description FROM categories where deleted_at isnull",
	}
}

func (r *CategoryRepo) Create(c common.Context, input CategoryInput) (uint, error) {
	ctx := c.Context()

	category, err := common.Cast[entities.Category](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.Category](r.db).Create(ctx, &category)

	return category.ID, err
}

func (r *CategoryRepo) Read(c fiber.Ctx) (paginate.Page, []CategoryPage) {
	categories := make([]CategoryPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&categories)

	return page, categories
}

func (r *CategoryRepo) Update(c common.Context, id any, input CategoryInput) error {
	ctx := c.Context()
	category, err := common.Cast[entities.Category](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.Category](r.db).Where("id = ?", id).Updates(ctx, category)

	return err
}

func (r *CategoryRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.Category](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *CategoryRepo) GetByID(c common.Context, id any) (entities.Category, error) {
	var category entities.Category
	ctx := c.Context()

	category, err := gorm.G[entities.Category](r.db).Where("id = ?", id).First(ctx)

	return category, err
}
