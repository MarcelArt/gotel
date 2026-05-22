package web

import "github.com/gofiber/fiber/v3"

// StaffViewModel represents the data required to render the staff tab.
type StaffViewModel struct {
	BaseViewModel
}

// StaffGet handles GET /staff requests.
func (h *WebHandler) StaffGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	user, err := h.userService.GetByID(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	vm := StaffViewModel{
		BaseViewModel: BaseViewModel{
			Title:     "Staff - Gotel",
			ActiveTab: "staff",
			User:      user,
		},
	}

	return h.renderTab(c, "staff", vm)
}
