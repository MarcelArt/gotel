package inventory_transactions

import (
	"github.com/gofiber/fiber/v3"
)

func Setup(v1 fiber.Router, s IInventoryTransactionService) {
	h := NewInventoryTransactionHandler(s)

	h.SetupRoutes(v1)
}
