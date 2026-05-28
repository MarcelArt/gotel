package services

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IHousekeepingTaskService interface {
	common.IBaseCrudService[entities.HousekeepingTask, models.HousekeepingTaskInput, models.HousekeepingTaskPage]
}

type HousekeepingTaskService struct {
	repo repositories.IHousekeepingTaskRepo
}

var _ IHousekeepingTaskService = &HousekeepingTaskService{}

func NewHousekeepingTaskService(repo repositories.IHousekeepingTaskRepo) *HousekeepingTaskService {
	return &HousekeepingTaskService{
		repo: repo,
	}
}

func (s *HousekeepingTaskService) Create(c common.Context, input models.HousekeepingTaskInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *HousekeepingTaskService) Read(c fiber.Ctx) (paginate.Page, []models.HousekeepingTaskPage) {
	return s.repo.Read(c)
}

func (s *HousekeepingTaskService) Update(c common.Context, id any, input models.HousekeepingTaskInput) error {
	return s.repo.Update(c, id, input)
}

func (s *HousekeepingTaskService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *HousekeepingTaskService) GetByID(c common.Context, id any) (entities.HousekeepingTask, error) {
	return s.repo.GetByID(c, id)
}
