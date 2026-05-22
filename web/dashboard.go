package web

import (
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
)

// BaseViewModel represents the common view parameters shared by all authenticated layouts.
type BaseViewModel struct {
	Title     string
	ActiveTab string
	User      entities.User
}

// DashboardViewModel represents the specific data required to render the dashboard view.
type DashboardViewModel struct {
	BaseViewModel
	Permissions []string
}

// DashboardGet handles GET / requests (main dashboard layout / dashboard tab).
func (h *WebHandler) DashboardGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	user, err := h.userService.GetByID(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	permissions, err := h.userService.GetPermissions(userID)
	if err != nil {
		permissions = []string{}
	}

	vm := DashboardViewModel{
		BaseViewModel: BaseViewModel{
			Title:     "Dashboard - Gotel",
			ActiveTab: "dashboard",
			User:      user,
		},
		Permissions: permissions,
	}

	return h.renderTab(c, "dashboard", vm)
}
