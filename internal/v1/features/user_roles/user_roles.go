package user_roles

import (
	"github.com/gofiber/fiber/v3"
)

func Setup(v1 fiber.Router, s IUserRoleService) {
	h := NewUserRoleHandler(s)

	h.SetupRoutes(v1)
}
