package items

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IItemService interface {
	common.IBaseCrudService[entities.Item, ItemInput, ItemPage]
}

type ItemService struct {
	repo IItemRepo
}

var _ IItemService = &ItemService{}

func NewItemService(repo IItemRepo) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

func (s *ItemService) Create(c common.Context, input ItemInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *ItemService) Read(c fiber.Ctx) (paginate.Page, []ItemPage) {
	return s.repo.Read(c)
}

func (s *ItemService) Update(c common.Context, id any, input ItemInput) error {
	return s.repo.Update(c, id, input)
}

func (s *ItemService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *ItemService) GetByID(c common.Context, id any) (entities.Item, error) {
	return s.repo.GetByID(c, id)
}
