package v1

import (
	"github.com/MarcelArt/gotel/internal/configs"
	"github.com/MarcelArt/gotel/internal/v1/features/roles"
	"github.com/MarcelArt/gotel/internal/v1/features/users"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router) {
	v1 := api.Group("/v1")

	uRepo := users.NewUserRepo(configs.DB)
	uService := users.NewUserService(uRepo)

	rRepo := roles.NewRoleRepo(configs.DB)
	rService := roles.NewRoleService(rRepo)

	users.Setup(v1, uService)
	roles.Setup(v1, rService)
}
