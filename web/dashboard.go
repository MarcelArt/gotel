package web

import (
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/gofiber/fiber/v3"
)

// BaseViewModel represents the common view parameters shared by all authenticated layouts.
type BaseViewModel struct {
	Title       string
	ActiveTab   string
	User        entities.User
	Permissions []string
}

func (vm BaseViewModel) HasPermission(permissionKey string) bool {
	for _, p := range vm.Permissions {
		if p == "fullAccess" {
			return true
		}
	}
	for _, p := range vm.Permissions {
		if p == permissionKey {
			return true
		}
	}
	return false
}

func (vm BaseViewModel) PermissionName(key string) string {
	for _, p := range AvailablePermissions {
		if p.Key == key {
			return p.Name
		}
	}
	return key
}

func (vm BaseViewModel) GetTitle() string {
	return vm.Title
}


// DashboardViewModel represents the specific data required to render the dashboard view.
type DashboardViewModel struct {
	BaseViewModel
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

	vm := DashboardViewModel{
		BaseViewModel: BaseViewModel{
			Title:       "Dashboard - Gotel",
			ActiveTab:   "dashboard",
			User:        user,
			Permissions: getPermissions(c),
		},
	}

	return h.renderTab(c, "dashboard", vm)
}
