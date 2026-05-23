package web

import (
	"embed"
	"html/template"

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
	views["settings"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/settings_tab.html"))

	// Instantiate services
	uRepo := users.NewUserRepo(configs.DB)
	rRepo := roles.NewRoleRepo(configs.DB)
	urRepo := user_roles.NewUserRoleRepo(configs.DB)
	cRepo := categories.NewCategoryRepo(configs.DB)
	locRepo := locations.NewLocationRepo(configs.DB)
	itemRepo := items.NewItemRepo(configs.DB)
	itRepo := inventory_transactions.NewInventoryTransactionRepo(configs.DB)

	uService := users.NewUserService(uRepo, urRepo)
	rService := roles.NewRoleService(rRepo)
	cService := categories.NewCategoryService(cRepo)
	locService := locations.NewLocationService(locRepo)
	itemService := items.NewItemService(itemRepo)
	itService := inventory_transactions.NewInventoryTransactionService(itRepo)

	h := NewWebHandler(uService, rService, cService, locService, itemService, itService)

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
	authGroup.Get("/roles", h.RolesGet)
	authGroup.Post("/roles", h.RolesPost)
	authGroup.Get("/roles/:id/edit", h.RolesEditGet)
	authGroup.Put("/roles/:id", h.RolesPut)
	authGroup.Delete("/roles/:id", h.RolesDelete)

	// Users routes
	authGroup.Get("/users", h.UsersGet)
	authGroup.Post("/users", h.UsersPost)
	authGroup.Get("/users/:id/edit", h.UsersEditGet)
	authGroup.Put("/users/:id", h.UsersPut)
	authGroup.Delete("/users/:id", h.UsersDelete)
	authGroup.Get("/users/roles/list", h.UserRolesListGet)

	// Categories routes
	authGroup.Get("/categories", h.CategoriesGet)
	authGroup.Post("/categories", h.CategoriesPost)
	authGroup.Get("/categories/:id/edit", h.CategoriesEditGet)
	authGroup.Put("/categories/:id", h.CategoriesPut)
	authGroup.Delete("/categories/:id", h.CategoriesDelete)
	authGroup.Get("/categories/options", h.CategoriesOptionsGet)

	// Locations routes
	authGroup.Get("/locations", h.LocationsGet)
	authGroup.Post("/locations", h.LocationsPost)
	authGroup.Get("/locations/:id/edit", h.LocationsEditGet)
	authGroup.Put("/locations/:id", h.LocationsPut)
	authGroup.Delete("/locations/:id", h.LocationsDelete)
	authGroup.Get("/locations/options", h.LocationsOptionsGet)

	// Items routes
	authGroup.Get("/items", h.ItemsGet)
	authGroup.Post("/items", h.ItemsPost)
	authGroup.Get("/items/:id/edit", h.ItemsEditGet)
	authGroup.Put("/items/:id", h.ItemsPut)
	authGroup.Delete("/items/:id", h.ItemsDelete)

	// Inventory Transactions routes
	authGroup.Get("/inventory-transactions", h.InventoryTransactionsGet)
	authGroup.Post("/inventory-transactions", h.InventoryTransactionsPost)

	authGroup.Get("/settings", h.SettingsGet)
}
