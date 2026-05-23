package locations

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type ILocationRepo interface {
	common.IBaseCrudRepo[entities.Location, LocationInput, LocationPage]
}

type LocationRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ ILocationRepo = &LocationRepo{}

func NewLocationRepo(db *gorm.DB) *LocationRepo {
	return &LocationRepo{
		db:        db,
		pageQuery: "SELECT id, value, is_virtual FROM locations",
	}
}

func (r *LocationRepo) Create(c common.Context, input LocationInput) (uint, error) {
	ctx := c.Context()

	location, err := common.Cast[entities.Location](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.Location](r.db).Create(ctx, &location)

	return location.ID, err
}

func (r *LocationRepo) Read(c fiber.Ctx) (paginate.Page, []LocationPage) {
	locations := make([]LocationPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&locations)

	return page, locations
}

func (r *LocationRepo) Update(c common.Context, id any, input LocationInput) error {
	ctx := c.Context()
	location, err := common.Cast[entities.Location](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.Location](r.db).Where("id = ?", id).Updates(ctx, location)

	return err
}

func (r *LocationRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.Location](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *LocationRepo) GetByID(c common.Context, id any) (entities.Location, error) {
	var location entities.Location
	ctx := c.Context()

	location, err := gorm.G[entities.Location](r.db).Where("id = ?", id).First(ctx)

	return location, err
}
