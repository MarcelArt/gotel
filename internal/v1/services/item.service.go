package services

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IItemService interface {
	common.IBaseCrudService[entities.Item, models.ItemInput, models.ItemPage]
}

type ItemService struct {
	repo repositories.IItemRepo
}

var _ IItemService = &ItemService{}

func NewItemService(repo repositories.IItemRepo) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

func (s *ItemService) Create(c common.Context, input models.ItemInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *ItemService) Read(c fiber.Ctx) (paginate.Page, []models.ItemPage) {
	return s.repo.Read(c)
}

func (s *ItemService) Update(c common.Context, id any, input models.ItemInput) error {
	return s.repo.Update(c, id, input)
}

func (s *ItemService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *ItemService) GetByID(c common.Context, id any) (entities.Item, error) {
	return s.repo.GetByID(c, id)
}
