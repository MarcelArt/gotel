package v1

import (
	"github.com/MarcelArt/gotel/internal/configs"
	"github.com/MarcelArt/gotel/internal/v1/features/categories"
	"github.com/MarcelArt/gotel/internal/v1/features/inventory_transactions"
	"github.com/MarcelArt/gotel/internal/v1/features/items"
	"github.com/MarcelArt/gotel/internal/v1/features/locations"
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

	catRepo := categories.NewCategoryRepo(configs.DB)
	catService := categories.NewCategoryService(catRepo)

	categories.Setup(v1, catService)

	itemRepo := items.NewItemRepo(configs.DB)
	itemService := items.NewItemService(itemRepo)

	items.Setup(v1, itemService)

	locRepo := locations.NewLocationRepo(configs.DB)
	locService := locations.NewLocationService(locRepo)

	locations.Setup(v1, locService)

	txRepo := inventory_transactions.NewInventoryTransactionRepo(configs.DB)
	txService := inventory_transactions.NewInventoryTransactionService(txRepo)

	inventory_transactions.Setup(v1, txService)
}
