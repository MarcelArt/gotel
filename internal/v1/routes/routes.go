package routes

import (
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router) {
	v1 := api.Group("/v1")

	SetupUserRoutes(v1)
	SetupRoleRoutes(v1)
	SetupUserRoleRoutes(v1)
	SetupCategoryRoutes(v1)
	SetupItemRoutes(v1)
	SetupLocationRoutes(v1)
	SetupInventoryTransactionRoutes(v1)
	SetupAssetInstanceRoutes(v1)
	SetupAssetTransactionRoutes(v1)
}
