package users

import (
	"github.com/gofiber/fiber/v3"
)

func Setup(v1 fiber.Router, s IUserService) {
	h := NewUserHandler(s)

	h.SetupRoutes(v1)
}
