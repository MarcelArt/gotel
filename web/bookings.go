package web

import "github.com/gofiber/fiber/v3"

// BookingsViewModel represents the data required to render the bookings tab.
type BookingsViewModel struct {
	BaseViewModel
}

// BookingsGet handles GET /bookings requests.
func (h *WebHandler) BookingsGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	user, err := h.userService.GetByID(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	vm := BookingsViewModel{
		BaseViewModel: BaseViewModel{
			Title:     "Bookings - Gotel",
			ActiveTab: "bookings",
			User:      user,
		},
	}

	return h.renderTab(c, "bookings", vm)
}
