package items

import (
	"github.com/gofiber/fiber/v3"
)

func Setup(v1 fiber.Router, s IItemService) {
	h := NewItemHandler(s)

	h.SetupRoutes(v1)
}
