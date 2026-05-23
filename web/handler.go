package web

import (
	"bytes"
	"html/template"

	"github.com/MarcelArt/gotel/internal/v1/features/categories"
	"github.com/MarcelArt/gotel/internal/v1/features/inventory_transactions"
	"github.com/MarcelArt/gotel/internal/v1/features/items"
	"github.com/MarcelArt/gotel/internal/v1/features/locations"
	"github.com/MarcelArt/gotel/internal/v1/features/roles"
	"github.com/MarcelArt/gotel/internal/v1/features/users"
	"github.com/gofiber/fiber/v3"
)

var views = make(map[string]*template.Template)

type WebHandler struct {
	userService                 users.IUserService
	roleService                 roles.IRoleService
	categoryService             categories.ICategoryService
	locationService             locations.ILocationService
	itemService                 items.IItemService
	inventoryTransactionService inventory_transactions.IInventoryTransactionService
}

func NewWebHandler(
	userService users.IUserService,
	roleService roles.IRoleService,
	categoryService categories.ICategoryService,
	locationService locations.ILocationService,
	itemService items.IItemService,
	inventoryTransactionService inventory_transactions.IInventoryTransactionService,
) *WebHandler {
	return &WebHandler{
		userService:                 userService,
		roleService:                 roleService,
		categoryService:             categoryService,
		locationService:             locationService,
		itemService:                 itemService,
		inventoryTransactionService: inventoryTransactionService,
	}
}

// Helper to render templates with a full layout
func (h *WebHandler) render(c fiber.Ctx, page string, data any) error {
	t, ok := views[page]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found: " + page)
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "layout", data); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}

// Helper to render layout tabs or direct inner tab fragments (HTMX swaps)
func (h *WebHandler) renderTab(c fiber.Ctx, page string, data any) error {
	t, ok := views[page]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found: " + page)
	}

	var buf bytes.Buffer
	var err error
	if c.Get("HX-Request") == "true" {
		err = t.ExecuteTemplate(&buf, "outlet", data)
	} else {
		err = t.ExecuteTemplate(&buf, "layout", data)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}
