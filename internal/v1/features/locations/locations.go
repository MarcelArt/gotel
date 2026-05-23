package locations

import (
	"github.com/gofiber/fiber/v3"
)

func Setup(v1 fiber.Router, s ILocationService) {
	h := NewLocationHandler(s)

	h.SetupRoutes(v1)
}
