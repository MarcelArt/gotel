package roles

import (
	"github.com/gofiber/fiber/v3"
)

func Setup(v1 fiber.Router, s IRoleService) {
	h := NewRoleHandler(s)

	h.SetupRoutes(v1)
}
