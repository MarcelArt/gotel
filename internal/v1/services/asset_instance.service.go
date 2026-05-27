package services

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IAssetInstanceService interface {
	common.IBaseCrudService[entities.AssetInstance, models.AssetInstanceInput, models.AssetInstancePage]
}

type AssetInstanceService struct {
	repo repositories.IAssetInstanceRepo
}

var _ IAssetInstanceService = &AssetInstanceService{}

func NewAssetInstanceService(repo repositories.IAssetInstanceRepo) *AssetInstanceService {
	return &AssetInstanceService{
		repo: repo,
	}
}

func (s *AssetInstanceService) Create(c common.Context, input models.AssetInstanceInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *AssetInstanceService) Read(c fiber.Ctx) (paginate.Page, []models.AssetInstancePage) {
	return s.repo.Read(c)
}

func (s *AssetInstanceService) Update(c common.Context, id any, input models.AssetInstanceInput) error {
	return s.repo.Update(c, id, input)
}

func (s *AssetInstanceService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *AssetInstanceService) GetByID(c common.Context, id any) (entities.AssetInstance, error) {
	return s.repo.GetByID(c, id)
}
