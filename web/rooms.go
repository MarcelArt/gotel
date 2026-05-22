package web

import "github.com/gofiber/fiber/v3"

// RoomsViewModel represents the data required to render the rooms tab.
type RoomsViewModel struct {
	BaseViewModel
}

// RoomsGet handles GET /rooms requests.
func (h *WebHandler) RoomsGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	user, err := h.userService.GetByID(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	vm := RoomsViewModel{
		BaseViewModel: BaseViewModel{
			Title:     "Rooms - Gotel",
			ActiveTab: "rooms",
			User:      user,
		},
	}

	return h.renderTab(c, "rooms", vm)
}
