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
	rRepo := roles.NewRoleRepo(configs.DB)
	urRepo := user_roles.NewUserRoleRepo(configs.DB)

	uService := users.NewUserService(uRepo, urRepo)
	rService := roles.NewRoleService(rRepo)
	urService := user_roles.NewUserRoleService(urRepo)

	users.Setup(v1, uService)
	roles.Setup(v1, rService)
	user_roles.Setup(v1, urService)
}
