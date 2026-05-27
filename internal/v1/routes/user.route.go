package routes

import (
	"github.com/MarcelArt/gotel/internal/configs"
	"github.com/MarcelArt/gotel/internal/v1/handlers"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/MarcelArt/gotel/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupUserRoutes(v1 fiber.Router) {
	uRepo := repositories.NewUserRepo(configs.DB)
	urRepo := repositories.NewUserRoleRepo(configs.DB)
	h := handlers.NewUserHandler(
		services.NewUserService(configs.DB, uRepo, urRepo),
	)
	h.SetupRoutes(v1)
}
