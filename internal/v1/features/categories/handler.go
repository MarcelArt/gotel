package categories

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/v1/middlewares"
	"github.com/gofiber/fiber/v3"
	_ "github.com/morkid/paginate"
)

type CategoryHandler struct {
	service ICategoryService
}

func NewCategoryHandler(service ICategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new category
// @Description  Create a new category with the provided details
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      CategoryInput  true  "Category details"
// @Success      201       {object}  common.JSONResponse{items=uint}
// @Failure      400       {object}  common.JSONResponse
// @Failure      500       {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/categories [post]
func (h *CategoryHandler) Create(c fiber.Ctx) error {
	var category CategoryInput
	if err := c.Bind().JSON(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, category)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed creating category"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Category created"))
}

// Read godoc
// @Summary      List categories
// @Description  Get a paginated list of categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param			page		query		int		false	"Page"
// @Param			size		query		int		false	"Size"
// @Param			sort		query		string	false	"Sort"
// @Param			filters		query		string	false	"Filter"
// @Success      200    {object}  paginate.Page{items=[]CategoryPage}
// @Failure      401    {object}  common.JSONResponse
// @Failure      500    {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/categories [get]
func (h *CategoryHandler) Read(c fiber.Ctx) error {
	categories, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(categories)
}

// Update godoc
// @Summary      Update category
// @Description  Update an existing category's details
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id        path      string         true  "Category ID"
// @Param        category  body      CategoryInput  true  "Updated category details"
// @Success      200       {object}  common.JSONResponse
// @Failure      400       {object}  common.JSONResponse
// @Failure      401       {object}  common.JSONResponse
// @Failure      500       {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/categories/{id} [put]
func (h *CategoryHandler) Update(c fiber.Ctx) error {
	var category CategoryInput
	if err := c.Bind().JSON(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), category); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating category"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Category updated"))
}

// Delete godoc
// @Summary      Delete category
// @Description  Delete a category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/categories/{id} [delete]
func (h *CategoryHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting category"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Category deleted"))
}

// GetByID godoc
// @Summary      Get category by ID
// @Description  Get detailed information about a category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  common.JSONResponse{items=entities.Category}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/categories/{id} [get]
func (h *CategoryHandler) GetByID(c fiber.Ctx) error {
	category, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting category"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(category, "Category found"))
}

func (h *CategoryHandler) SetupRoutes(v1 fiber.Router) {
	categories := v1.Group("/categories")

	categories.Use(middlewares.Authn())

	categories.Post("/", middlewares.Authz("categories#create"), h.Create)
	categories.Get("/", middlewares.Authz("categories#read"), h.Read)
	categories.Get("/:id", middlewares.Authz("categories#read"), h.GetByID)
	categories.Put("/:id", middlewares.Authz("categories#update"), h.Update)
	categories.Delete("/:id", middlewares.Authz("categories#delete"), h.Delete)
}
