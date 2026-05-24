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
	views["licenses"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/licenses_tab.html"))
	views["unauthorized"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/unauthorized_tab.html"))

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

	authGroup.Get("/settings", h.SettingsGet)
	authGroup.Get("/licenses", h.LicensesGet)
}
