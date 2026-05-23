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
