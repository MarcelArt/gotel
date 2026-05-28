package services

import (
	"errors"
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/configs"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IRoomService interface {
	common.IBaseCrudService[entities.Room, models.RoomInput, models.RoomPage]
	AssignCleaning(c common.Context, taskInput models.HousekeepingTaskInput) error
}

type RoomService struct {
	repo repositories.IRoomRepo
}

var _ IRoomService = &RoomService{}

func NewRoomService(repo repositories.IRoomRepo) *RoomService {
	return &RoomService{
		repo: repo,
	}
}

func (s *RoomService) Create(c common.Context, input models.RoomInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *RoomService) Read(c fiber.Ctx) (paginate.Page, []models.RoomPage) {
	return s.repo.Read(c)
}

func (s *RoomService) Update(c common.Context, id any, input models.RoomInput) error {
	return s.repo.Update(c, id, input)
}

func (s *RoomService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *RoomService) GetByID(c common.Context, id any) (entities.Room, error) {
	return s.repo.GetByID(c, id)
}

func (s *RoomService) AssignCleaning(c common.Context, taskInput models.HousekeepingTaskInput) error {
	tx := configs.DB.Begin()
	defer tx.Rollback()

	htRepo := repositories.NewHousekeepingTaskRepo(tx)
	repo := repositories.NewRoomRepo(tx)

	task, err := htRepo.GetActiveTaskByRoomID(c, taskInput.RoomID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed retrieving active task: %w", err)
	}

	if task.ID != 0 {
		return fmt.Errorf("room is already assigned to a cleaning task")
	}

	room := models.RoomInput{Status: "DIRTY"}
	if err := repo.Update(c, taskInput.RoomID, room); err != nil {
		return fmt.Errorf("failed updating room: %w", err)
	}

	if _, err := htRepo.Create(c, taskInput); err != nil {
		return fmt.Errorf("failed creating housekeeping task: %w", err)
	}

	return tx.Commit().Error
}
