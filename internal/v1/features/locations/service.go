package locations

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type ILocationService interface {
	common.IBaseCrudService[entities.Location, LocationInput, LocationPage]
}

type LocationService struct {
	repo ILocationRepo
}

var _ ILocationService = &LocationService{}

func NewLocationService(repo ILocationRepo) *LocationService {
	return &LocationService{
		repo: repo,
	}
}

func (s *LocationService) Create(c common.Context, input LocationInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *LocationService) Read(c fiber.Ctx) (paginate.Page, []LocationPage) {
	return s.repo.Read(c)
}

func (s *LocationService) Update(c common.Context, id any, input LocationInput) error {
	return s.repo.Update(c, id, input)
}

func (s *LocationService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *LocationService) GetByID(c common.Context, id any) (entities.Location, error) {
	return s.repo.GetByID(c, id)
}
