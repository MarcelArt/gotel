package v1

import (
	"github.com/MarcelArt/gotel/internal/configs"
	"github.com/MarcelArt/gotel/internal/v1/features/roles"
	"github.com/MarcelArt/gotel/internal/v1/features/user_roles"
	"github.com/MarcelArt/gotel/internal/v1/features/users"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router) {
	v1 := api.Group("/v1")

	uRepo := users.NewUserRepo(configs.DB)
	uService := users.NewUserService(uRepo)

	users.Setup(v1, uService)

	rRepo := roles.NewRoleRepo(configs.DB)
	rService := roles.NewRoleService(rRepo)

	roles.Setup(v1, rService)

	urRepo := user_roles.NewUserRoleRepo(configs.DB)
	urService := user_roles.NewUserRoleService(urRepo)

	user_roles.Setup(v1, urService)
}
