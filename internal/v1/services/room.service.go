package services

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IRoomService interface {
	common.IBaseCrudService[entities.Room, models.RoomInput, models.RoomPage]
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
