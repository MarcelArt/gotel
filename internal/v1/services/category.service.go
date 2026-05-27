package services

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type ICategoryService interface {
	common.IBaseCrudService[entities.Category, models.CategoryInput, models.CategoryPage]
}

type CategoryService struct {
	repo repositories.ICategoryRepo
}

var _ ICategoryService = &CategoryService{}

func NewCategoryService(repo repositories.ICategoryRepo) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) Create(c common.Context, input models.CategoryInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *CategoryService) Read(c fiber.Ctx) (paginate.Page, []models.CategoryPage) {
	return s.repo.Read(c)
}

func (s *CategoryService) Update(c common.Context, id any, input models.CategoryInput) error {
	return s.repo.Update(c, id, input)
}

func (s *CategoryService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *CategoryService) GetByID(c common.Context, id any) (entities.Category, error) {
	return s.repo.GetByID(c, id)
}
