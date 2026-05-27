package services

import (
	"time"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IInventoryTransactionService interface {
	common.IBaseCrudService[entities.InventoryTransaction, models.InventoryTransactionInput, models.InventoryTransactionPage]
	GetItemCounts(itemID any, timeRanges ...time.Time) ([]models.ItemCount, error)
	GetItemActualQuantity(itemID any) (float64, error)
}

type InventoryTransactionService struct {
	repo repositories.IInventoryTransactionRepo
}

var _ IInventoryTransactionService = &InventoryTransactionService{}

func NewInventoryTransactionService(repo repositories.IInventoryTransactionRepo) *InventoryTransactionService {
	return &InventoryTransactionService{
		repo: repo,
	}
}

func (s *InventoryTransactionService) Create(c common.Context, input models.InventoryTransactionInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *InventoryTransactionService) Read(c fiber.Ctx) (paginate.Page, []models.InventoryTransactionPage) {
	return s.repo.Read(c)
}

func (s *InventoryTransactionService) Update(c common.Context, id any, input models.InventoryTransactionInput) error {
	return s.repo.Update(c, id, input)
}

func (s *InventoryTransactionService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *InventoryTransactionService) GetByID(c common.Context, id any) (entities.InventoryTransaction, error) {
	return s.repo.GetByID(c, id)
}

func (s *InventoryTransactionService) GetItemCounts(itemID any, timeRanges ...time.Time) ([]models.ItemCount, error) {
	return s.repo.GetItemCounts(itemID, timeRanges...)
}

func (s *InventoryTransactionService) GetItemActualQuantity(itemID any) (float64, error) {
	return s.repo.GetItemActualQuantity(itemID)
}
