package web

import (
	"embed"
	"html/template"

	"github.com/MarcelArt/gotel/internal/configs"
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
	views["settings"] = template.Must(template.New("").ParseFS(templatesFS, "templates/layout.html", "templates/dashboard.html", "templates/settings_tab.html"))

	// Instantiate services
	uRepo := users.NewUserRepo(configs.DB)
	rRepo := roles.NewRoleRepo(configs.DB)
	urRepo := user_roles.NewUserRoleRepo(configs.DB)
	uService := users.NewUserService(uRepo, urRepo)
	rService := roles.NewRoleService(rRepo)

	h := NewWebHandler(uService, rService)

	// Register routes
	app.Get("/login", h.LoginGet)
	app.Post("/login", h.LoginPost)

	app.Get("/register", h.RegisterGet)
	app.Post("/register", h.RegisterPost)

	app.Post("/logout", h.LogoutPost)

	// Dashboard and authenticated routes
	authGroup := app.Group("", WebAuth(uService))
	authGroup.Get("/", h.DashboardGet)
	authGroup.Get("/roles", h.RolesGet)
	authGroup.Post("/roles", h.RolesPost)
	authGroup.Get("/roles/:id/edit", h.RolesEditGet)
	authGroup.Put("/roles/:id", h.RolesPut)
	authGroup.Delete("/roles/:id", h.RolesDelete)
	authGroup.Get("/settings", h.SettingsGet)
}
