package services

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IAssetTransactionService interface {
	common.IBaseCrudService[entities.AssetTransaction, models.AssetTransactionInput, models.AssetTransactionPage]
}

type AssetTransactionService struct {
	repo repositories.IAssetTransactionRepo
}

var _ IAssetTransactionService = &AssetTransactionService{}

func NewAssetTransactionService(repo repositories.IAssetTransactionRepo) *AssetTransactionService {
	return &AssetTransactionService{
		repo: repo,
	}
}

func (s *AssetTransactionService) Create(c common.Context, input models.AssetTransactionInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *AssetTransactionService) Read(c fiber.Ctx) (paginate.Page, []models.AssetTransactionPage) {
	return s.repo.Read(c)
}

func (s *AssetTransactionService) Update(c common.Context, id any, input models.AssetTransactionInput) error {
	return s.repo.Update(c, id, input)
}

func (s *AssetTransactionService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *AssetTransactionService) GetByID(c common.Context, id any) (entities.AssetTransaction, error) {
	return s.repo.GetByID(c, id)
}
