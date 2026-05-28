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

type IRoomRepo interface {
	common.IBaseCrudRepo[entities.Room, models.RoomInput, models.RoomPage]
}

type RoomRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IRoomRepo = &RoomRepo{}

func NewRoomRepo(db *gorm.DB) *RoomRepo {
	return &RoomRepo{
		db:        db,
		pageQuery: "SELECT id, room_number, floor, status FROM rooms where deleted_at isnull",
	}
}

func (r *RoomRepo) Create(c common.Context, input models.RoomInput) (uint, error) {
	ctx := c.Context()

	room, err := common.Cast[entities.Room](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.Room](r.db).Create(ctx, &room)

	return room.ID, err
}

func (r *RoomRepo) Read(c fiber.Ctx) (paginate.Page, []models.RoomPage) {
	rooms := make([]models.RoomPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&rooms)

	return page, rooms
}

func (r *RoomRepo) Update(c common.Context, id any, input models.RoomInput) error {
	ctx := c.Context()
	room, err := common.Cast[entities.Room](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.Room](r.db).Where("id = ?", id).Updates(ctx, room)

	return err
}

func (r *RoomRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.Room](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *RoomRepo) GetByID(c common.Context, id any) (entities.Room, error) {
	var room entities.Room
	ctx := c.Context()

	room, err := gorm.G[entities.Room](r.db).Where("id = ?", id).First(ctx)

	return room, err
}
