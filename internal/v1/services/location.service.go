package services

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type ILocationService interface {
	common.IBaseCrudService[entities.Location, models.LocationInput, models.LocationPage]
}

type LocationService struct {
	repo repositories.ILocationRepo
}

var _ ILocationService = &LocationService{}

func NewLocationService(repo repositories.ILocationRepo) *LocationService {
	return &LocationService{
		repo: repo,
	}
}

func (s *LocationService) Create(c common.Context, input models.LocationInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *LocationService) Read(c fiber.Ctx) (paginate.Page, []models.LocationPage) {
	return s.repo.Read(c)
}

func (s *LocationService) Update(c common.Context, id any, input models.LocationInput) error {
	return s.repo.Update(c, id, input)
}

func (s *LocationService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *LocationService) GetByID(c common.Context, id any) (entities.Location, error) {
	return s.repo.GetByID(c, id)
}
