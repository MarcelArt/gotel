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

type IHousekeepingTaskRepo interface {
	common.IBaseCrudRepo[entities.HousekeepingTask, models.HousekeepingTaskInput, models.HousekeepingTaskPage]
}

type HousekeepingTaskRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IHousekeepingTaskRepo = &HousekeepingTaskRepo{}

func NewHousekeepingTaskRepo(db *gorm.DB) *HousekeepingTaskRepo {
	return &HousekeepingTaskRepo{
		db:        db,
		pageQuery: "SELECT id, priority, started_at, completed_at, note, assignee_id, assigner_id, room_id FROM housekeeping_tasks where deleted_at isnull",
	}
}

func (r *HousekeepingTaskRepo) Create(c common.Context, input models.HousekeepingTaskInput) (uint, error) {
	ctx := c.Context()

	task, err := common.Cast[entities.HousekeepingTask](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}

	err = gorm.G[entities.HousekeepingTask](r.db).Create(ctx, &task)

	return task.ID, err
}

func (r *HousekeepingTaskRepo) Read(c fiber.Ctx) (paginate.Page, []models.HousekeepingTaskPage) {
	tasks := make([]models.HousekeepingTaskPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request()).Response(&tasks)

	return page, tasks
}

func (r *HousekeepingTaskRepo) Update(c common.Context, id any, input models.HousekeepingTaskInput) error {
	ctx := c.Context()
	task, err := common.Cast[entities.HousekeepingTask](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}

	_, err = gorm.G[entities.HousekeepingTask](r.db).Where("id = ?", id).Updates(ctx, task)

	return err
}

func (r *HousekeepingTaskRepo) Delete(c common.Context, id any) error {
	ctx := c.Context()
	_, err := gorm.G[entities.HousekeepingTask](r.db).Where("id = ?", id).Delete(ctx)

	return err
}

func (r *HousekeepingTaskRepo) GetByID(c common.Context, id any) (entities.HousekeepingTask, error) {
	var task entities.HousekeepingTask
	ctx := c.Context()

	task, err := gorm.G[entities.HousekeepingTask](r.db).Where("id = ?", id).First(ctx)

	return task, err
}
