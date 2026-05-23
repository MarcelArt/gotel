package items

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/v1/middlewares"
	"github.com/gofiber/fiber/v3"
	_ "github.com/morkid/paginate"
)

type ItemHandler struct {
	service IItemService
}

func NewItemHandler(service IItemService) *ItemHandler {
	return &ItemHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new item
// @Description  Create a new item with the provided details
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        item  body      ItemInput  true  "Item details"
// @Success      201   {object}  common.JSONResponse{items=uint}
// @Failure      400   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/items [post]
func (h *ItemHandler) Create(c fiber.Ctx) error {
	var item ItemInput
	if err := c.Bind().JSON(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, item)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed creating item"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Item created"))
}

// Read godoc
// @Summary      List items
// @Description  Get a paginated list of items
// @Tags         items
// @Accept       json
// @Produce      json
// @Param			page		query		int		false	"Page"
// @Param			size		query		int		false	"Size"
// @Param			sort		query		string	false	"Sort"
// @Param			filters		query		string	false	"Filter"
// @Success      200    {object}  paginate.Page{items=[]ItemPage}
// @Failure      401    {object}  common.JSONResponse
// @Failure      500    {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/items [get]
func (h *ItemHandler) Read(c fiber.Ctx) error {
	items, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(items)
}

// Update godoc
// @Summary      Update item
// @Description  Update an existing item's details
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id    path      string     true  "Item ID"
// @Param        item  body      ItemInput  true  "Updated item details"
// @Success      200   {object}  common.JSONResponse
// @Failure      400   {object}  common.JSONResponse
// @Failure      401   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/items/{id} [put]
func (h *ItemHandler) Update(c fiber.Ctx) error {
	var item ItemInput
	if err := c.Bind().JSON(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), item); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating item"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Item updated"))
}

// Delete godoc
// @Summary      Delete item
// @Description  Delete an item by ID
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Item ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/items/{id} [delete]
func (h *ItemHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting item"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Item deleted"))
}

// GetByID godoc
// @Summary      Get item by ID
// @Description  Get detailed information about an item by its ID
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Item ID"
// @Success      200  {object}  common.JSONResponse{items=entities.Item}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/items/{id} [get]
func (h *ItemHandler) GetByID(c fiber.Ctx) error {
	item, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting item"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(item, "Item found"))
}

func (h *ItemHandler) SetupRoutes(v1 fiber.Router) {
	items := v1.Group("/items")

	items.Use(middlewares.Authn())

	items.Post("/", middlewares.Authz("items#create"), h.Create)
	items.Get("/", middlewares.Authz("items#read"), h.Read)
	items.Get("/:id", middlewares.Authz("items#read"), h.GetByID)
	items.Put("/:id", middlewares.Authz("items#update"), h.Update)
	items.Delete("/:id", middlewares.Authz("items#delete"), h.Delete)
}
