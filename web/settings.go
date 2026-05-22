package web

import "github.com/gofiber/fiber/v3"

// SettingsViewModel represents the data required to render the settings tab.
type SettingsViewModel struct {
	BaseViewModel
}

// SettingsGet handles GET /settings requests.
func (h *WebHandler) SettingsGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	user, err := h.userService.GetByID(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	vm := SettingsViewModel{
		BaseViewModel: BaseViewModel{
			Title:     "Settings - Gotel",
			ActiveTab: "settings",
			User:      user,
		},
	}

	return h.renderTab(c, "settings", vm)
}
