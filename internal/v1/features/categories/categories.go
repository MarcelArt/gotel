package categories

import (
	"github.com/gofiber/fiber/v3"
)

func Setup(v1 fiber.Router, s ICategoryService) {
	h := NewCategoryHandler(s)

	h.SetupRoutes(v1)
}
