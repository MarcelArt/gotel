package routes

import (
	"github.com/MarcelArt/gotel/internal/configs"
	"github.com/MarcelArt/gotel/internal/v1/handlers"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/MarcelArt/gotel/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupItemRoutes(v1 fiber.Router) {
	h := handlers.NewItemHandler(
		services.NewItemService(
			repositories.NewItemRepo(configs.DB),
		),
	)
	h.SetupRoutes(v1)
}
