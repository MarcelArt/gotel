package web

import (
	"github.com/gofiber/fiber/v3"
)

type PlaceholderViewModel struct {
	BaseViewModel
}

func (h *WebHandler) getPlaceholderViewModel(c fiber.Ctx, activeTab string, title string) (PlaceholderViewModel, error) {
	userID := c.Locals("userId")
	if userID == nil {
		return PlaceholderViewModel{}, fiber.ErrUnauthorized
	}

	currentUser, err := h.userService.GetByID(c, userID)
	if err != nil {
		return PlaceholderViewModel{}, err
	}

	return PlaceholderViewModel{
		BaseViewModel: BaseViewModel{
			Title:     title + " - Gotel",
			ActiveTab: activeTab,
			User:      currentUser,
		},
	}, nil
}

// LocationsGet handles GET /locations requests.
func (h *WebHandler) LocationsGet(c fiber.Ctx) error {
	vm, err := h.getPlaceholderViewModel(c, "locations", "Locations Directory")
	if err != nil {
		return h.LogoutPost(c)
	}
	return h.renderTab(c, "locations", vm)
}

// ItemsGet handles GET /items requests.
func (h *WebHandler) ItemsGet(c fiber.Ctx) error {
	vm, err := h.getPlaceholderViewModel(c, "items", "Inventory Items")
	if err != nil {
		return h.LogoutPost(c)
	}
	return h.renderTab(c, "items", vm)
}

// TransactionsGet handles GET /transactions requests.
func (h *WebHandler) TransactionsGet(c fiber.Ctx) error {
	vm, err := h.getPlaceholderViewModel(c, "transactions", "Inventory Transactions")
	if err != nil {
		return h.LogoutPost(c)
	}
	return h.renderTab(c, "transactions", vm)
}
