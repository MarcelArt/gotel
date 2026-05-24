package inventory_transactions

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type IInventoryTransactionService interface {
	common.IBaseCrudService[entities.InventoryTransaction, InventoryTransactionInput, InventoryTransactionPage]
	GetItemCounts(itemID any) ([]ItemCount, error)
}

type InventoryTransactionService struct {
	repo IInventoryTransactionRepo
}

var _ IInventoryTransactionService = &InventoryTransactionService{}

func NewInventoryTransactionService(repo IInventoryTransactionRepo) *InventoryTransactionService {
	return &InventoryTransactionService{
		repo: repo,
	}
}

func (s *InventoryTransactionService) Create(c common.Context, input InventoryTransactionInput) (uint, error) {
	return s.repo.Create(c, input)
}

func (s *InventoryTransactionService) Read(c fiber.Ctx) (paginate.Page, []InventoryTransactionPage) {
	return s.repo.Read(c)
}

func (s *InventoryTransactionService) Update(c common.Context, id any, input InventoryTransactionInput) error {
	return s.repo.Update(c, id, input)
}

func (s *InventoryTransactionService) Delete(c common.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *InventoryTransactionService) GetByID(c common.Context, id any) (entities.InventoryTransaction, error) {
	return s.repo.GetByID(c, id)
}

func (s *InventoryTransactionService) GetItemCounts(itemID any) ([]ItemCount, error) {
	itemCounts, err := s.repo.GetItemCounts(itemID)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving item counts: %w", err)
	}

	actualQuantity := ItemCount{TransactionType: "ACTUAL"}
	for _, itemCount := range itemCounts {
		if itemCount.TransactionType == "RECEIVE" {
			actualQuantity.Quantity += itemCount.Quantity
		} else {
			actualQuantity.Quantity -= itemCount.Quantity
		}
	}

	itemCounts = append(itemCounts, actualQuantity)

	return itemCounts, nil
}
