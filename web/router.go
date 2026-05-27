package web

import (
	"embed"
	"html/template"

	"github.com/MarcelArt/gotel/internal/configs"
	"github.com/MarcelArt/gotel/internal/v1/repositories"
	"github.com/MarcelArt/gotel/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

//go:embed templates/*
var templatesFS embed.FS

func SetupRoutes(app *fiber.App) {
	// Initialize templates parsed in isolated template pools to avoid conflict on "content" template name
	views["login"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/login.html"))
	views["register"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/register.html"))
	views["dashboard"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/dashboard_tab.html"))
	views["roles"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/roles_tab.html"))
	views["users"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/users_tab.html"))
	views["categories"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/categories_tab.html"))
	views["locations"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/locations_tab.html"))
	views["items"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/items_tab.html"))
	views["transactions"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/transactions_tab.html"))
	views["asset_instances"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/asset_instances_tab.html"))
	views["asset_transactions"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/asset_transactions_tab.html"))
	views["settings"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/settings_tab.html"))
	views["licenses"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/licenses_tab.html"))
	views["unauthorized"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/unauthorized_tab.html"))

	// Instantiate services
	uRepo := repositories.NewUserRepo(configs.DB)
	rRepo := repositories.NewRoleRepo(configs.DB)
	urRepo := repositories.NewUserRoleRepo(configs.DB)
	cRepo := repositories.NewCategoryRepo(configs.DB)
	locRepo := repositories.NewLocationRepo(configs.DB)
	itemRepo := repositories.NewItemRepo(configs.DB)
	itRepo := repositories.NewInventoryTransactionRepo(configs.DB)
	aiRepo := repositories.NewAssetInstanceRepo(configs.DB)
	atRepo := repositories.NewAssetTransactionRepo(configs.DB)

	uService := services.NewUserService(configs.DB, uRepo, urRepo)
	rService := services.NewRoleService(rRepo)
	cService := services.NewCategoryService(cRepo)
	locService := services.NewLocationService(locRepo)
	itemService := services.NewItemService(itemRepo)
	itService := services.NewInventoryTransactionService(itRepo)
	aiService := services.NewAssetInstanceService(aiRepo)
	atService := services.NewAssetTransactionService(atRepo)

	h := NewWebHandler(uService, rService, cService, locService, itemService, itService, aiService, atService)

	// Register routes
	app.Get("/login", h.LoginGet)
	app.Post("/login", h.LoginPost)

	app.Get("/register", h.RegisterGet)
	app.Post("/register", h.RegisterPost)

	app.Post("/logout", h.LogoutPost)

	// Dashboard and authenticated routes
	authGroup := app.Group("", WebAuth(uService))
	authGroup.Get("/", h.DashboardGet)
	
	// Roles routes
	authGroup.Get("/roles", h.WebAuthz("roles#read"), h.RolesGet)
	authGroup.Post("/roles", h.WebAuthz("roles#create"), h.RolesPost)
	authGroup.Get("/roles/:id/edit", h.WebAuthz("roles#update"), h.RolesEditGet)
	authGroup.Put("/roles/:id", h.WebAuthz("roles#update"), h.RolesPut)
	authGroup.Delete("/roles/:id", h.WebAuthz("roles#delete"), h.RolesDelete)
	
	// Users routes
	authGroup.Get("/users", h.WebAuthz("users#read"), h.UsersGet)
	authGroup.Post("/users", h.WebAuthz("users#create"), h.UsersPost)
	authGroup.Get("/users/:id/edit", h.WebAuthz("users#update"), h.UsersEditGet)
	authGroup.Put("/users/:id", h.WebAuthz("users#update"), h.UsersPut)
	authGroup.Delete("/users/:id", h.WebAuthz("users#delete"), h.UsersDelete)
	authGroup.Get("/users/roles/list", h.WebAuthz("users#read"), h.UserRolesListGet)
	
	// Categories routes
	authGroup.Get("/categories", h.WebAuthz("categories#read"), h.CategoriesGet)
	authGroup.Post("/categories", h.WebAuthz("categories#create"), h.CategoriesPost)
	authGroup.Get("/categories/:id/edit", h.WebAuthz("categories#update"), h.CategoriesEditGet)
	authGroup.Put("/categories/:id", h.WebAuthz("categories#update"), h.CategoriesPut)
	authGroup.Delete("/categories/:id", h.WebAuthz("categories#delete"), h.CategoriesDelete)
	authGroup.Get("/categories/options", h.WebAuthz("categories#read"), h.CategoriesOptionsGet)
	
	// Locations routes
	authGroup.Get("/locations", h.WebAuthz("locations#read"), h.LocationsGet)
	authGroup.Post("/locations", h.WebAuthz("locations#create"), h.LocationsPost)
	authGroup.Get("/locations/:id/edit", h.WebAuthz("locations#update"), h.LocationsEditGet)
	authGroup.Put("/locations/:id", h.WebAuthz("locations#update"), h.LocationsPut)
	authGroup.Delete("/locations/:id", h.WebAuthz("locations#delete"), h.LocationsDelete)
	authGroup.Get("/locations/options", h.WebAuthz("locations#read"), h.LocationsOptionsGet)
	
	// Items routes
	authGroup.Get("/items", h.WebAuthz("items#read"), h.ItemsGet)
	authGroup.Post("/items", h.WebAuthz("items#create"), h.ItemsPost)
	authGroup.Get("/items/:id/edit", h.WebAuthz("items#update"), h.ItemsEditGet)
	authGroup.Put("/items/:id", h.WebAuthz("items#update"), h.ItemsPut)
	authGroup.Delete("/items/:id", h.WebAuthz("items#delete"), h.ItemsDelete)
	
	// Inventory Transactions routes
	authGroup.Get("/inventory-transactions", h.WebAuthz("inventoryTransactions#read"), h.InventoryTransactionsGet)
	authGroup.Post("/inventory-transactions", h.WebAuthz("inventoryTransactions#create"), h.InventoryTransactionsPost)

	// Asset Instances routes
	authGroup.Get("/asset-instances", h.WebAuthz("asset_instances#read"), h.AssetInstancesGet)
	authGroup.Post("/asset-instances", h.WebAuthz("asset_instances#create"), h.AssetInstancesPost)
	authGroup.Get("/asset-instances/:id/edit", h.WebAuthz("asset_instances#update"), h.AssetInstancesEditGet)
	authGroup.Put("/asset-instances/:id", h.WebAuthz("asset_instances#update"), h.AssetInstancesPut)
	authGroup.Delete("/asset-instances/:id", h.WebAuthz("asset_instances#delete"), h.AssetInstancesDelete)

	// Asset Transactions routes
	authGroup.Get("/asset-transactions", h.WebAuthz("asset_transactions#read"), h.AssetTransactionsGet)
	authGroup.Post("/asset-transactions", h.WebAuthz("asset_transactions#create"), h.AssetTransactionsPost)
	authGroup.Get("/asset-transactions/:id/edit", h.WebAuthz("asset_transactions#update"), h.AssetTransactionsEditGet)
	authGroup.Put("/asset-transactions/:id", h.WebAuthz("asset_transactions#update"), h.AssetTransactionsPut)
	authGroup.Delete("/asset-transactions/:id", h.WebAuthz("asset_transactions#delete"), h.AssetTransactionsDelete)

	authGroup.Get("/settings", h.SettingsGet)
	authGroup.Get("/licenses", h.LicensesGet)
}
