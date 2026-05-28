package routes

import (
	"github.com/MarcelArt/gotel/internal/configs"
	"github.com/MarcelArt/gotel/internal/v1/handlers"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/MarcelArt/gotel/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupHousekeepingTaskRoutes(v1 fiber.Router) {
	h := handlers.NewHousekeepingTaskHandler(
		services.NewHousekeepingTaskService(
			repositories.NewHousekeepingTaskRepo(configs.DB),
		),
	)
	h.SetupRoutes(v1)
}
