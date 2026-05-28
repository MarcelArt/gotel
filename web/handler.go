package web

import (
	"bytes"
	"html/template"

	"github.com/MarcelArt/gotel/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

var views = make(map[string]*template.Template)

type WebHandler struct {
	userService                 services.IUserService
	roleService                 services.IRoleService
	categoryService             services.ICategoryService
	locationService             services.ILocationService
	itemService                 services.IItemService
	inventoryTransactionService services.IInventoryTransactionService
	assetInstanceService        services.IAssetInstanceService
	assetTransactionService     services.IAssetTransactionService
	roomService                 services.IRoomService
	housekeepingTaskService     services.IHousekeepingTaskService
}

func NewWebHandler(
	userService services.IUserService,
	roleService services.IRoleService,
	categoryService services.ICategoryService,
	locationService services.ILocationService,
	itemService services.IItemService,
	inventoryTransactionService services.IInventoryTransactionService,
	assetInstanceService services.IAssetInstanceService,
	assetTransactionService services.IAssetTransactionService,
	roomService services.IRoomService,
	housekeepingTaskService services.IHousekeepingTaskService,
) *WebHandler {
	return &WebHandler{
		userService:                 userService,
		roleService:                 roleService,
		categoryService:             categoryService,
		locationService:             locationService,
		itemService:                 itemService,
		inventoryTransactionService: inventoryTransactionService,
		assetInstanceService:        assetInstanceService,
		assetTransactionService:     assetTransactionService,
		roomService:                 roomService,
		housekeepingTaskService:     housekeepingTaskService,
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
		if hasTitle, ok := data.(interface{ GetTitle() string }); ok {
			buf.WriteString("<title>" + template.HTMLEscapeString(hasTitle.GetTitle()) + "</title>\n")
		}
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
