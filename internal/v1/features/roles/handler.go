package roles

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/v1/middlewares"
	"github.com/gofiber/fiber/v3"
	_ "github.com/morkid/paginate"
)

type RoleHandler struct {
	service IRoleService
}

func NewRoleHandler(service IRoleService) *RoleHandler {
	return &RoleHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new role
// @Description  Create a new role with the provided details
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        role  body      RoleInput  true  "Role details"
// @Success      201   {object}  common.JSONResponse{items=uint}
// @Failure      400   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/roles [post]
func (h *RoleHandler) Create(c fiber.Ctx) error {
	var role RoleInput
	if err := c.Bind().JSON(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, role)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed creating role"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Role created"))
}

// Read godoc
// @Summary      List roles
// @Description  Get a paginated list of roles
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param			page		query		int		false	"Page"
// @Param			size		query		int		false	"Size"
// @Param			sort		query		string	false	"Sort"
// @Param			filters		query		string	false	"Filter"
// @Success      200    {object}  paginate.Page{items=[]RolePage}
// @Failure      401    {object}  common.JSONResponse
// @Failure      500    {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/roles [get]
func (h *RoleHandler) Read(c fiber.Ctx) error {
	roles, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(roles)
}

// Update godoc
// @Summary      Update role
// @Description  Update an existing role's details
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id    path      string     true  "Role ID"
// @Param        role  body      RoleInput  true  "Updated role details"
// @Success      200   {object}  common.JSONResponse
// @Failure      400   {object}  common.JSONResponse
// @Failure      401   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/roles/{id} [put]
func (h *RoleHandler) Update(c fiber.Ctx) error {
	var role RoleInput
	if err := c.Bind().JSON(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), role); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating role"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Role updated"))
}

// Delete godoc
// @Summary      Delete role
// @Description  Delete a role by ID
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/roles/{id} [delete]
func (h *RoleHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting role"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Role deleted"))
}

// GetByID godoc
// @Summary      Get role by ID
// @Description  Get detailed information about a role by its ID
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID"
// @Success      200  {object}  common.JSONResponse{items=entities.Role}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/roles/{id} [get]
func (h *RoleHandler) GetByID(c fiber.Ctx) error {
	role, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting role"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(role, "Role found"))
}

func (h *RoleHandler) SetupRoutes(v1 fiber.Router) {
	roles := v1.Group("/roles")

	// Apply authentication middleware (Authn) to all roles routes
	roles.Use(middlewares.Authn())

	roles.Post("/", middlewares.Authz("roles#create"), h.Create)
	roles.Get("/", middlewares.Authz("roles#read"), h.Read)
	roles.Get("/:id", middlewares.Authz("roles#read"), h.GetByID)
	roles.Put("/:id", middlewares.Authz("roles#update"), h.Update)
	roles.Delete("/:id", middlewares.Authz("roles#delete"), h.Delete)
}
