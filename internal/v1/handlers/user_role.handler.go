package handlers

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/middlewares"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/services"
	"github.com/gofiber/fiber/v3"
	_ "github.com/morkid/paginate"
)

type UserRoleHandler struct {
	service services.IUserRoleService
}

var _ = entities.UserRole{}

func NewUserRoleHandler(service services.IUserRoleService) *UserRoleHandler {
	return &UserRoleHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Assign a role to a user
// @Description  Map a role to a user with the provided details
// @Tags         user-roles
// @Accept       json
// @Produce      json
// @Param        user_role  body      models.UserRoleInput  true  "User role assignment details"
// @Success      201        {object}  common.JSONResponse{items=uint}
// @Failure      400        {object}  common.JSONResponse
// @Failure      500        {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/user-roles [post]
func (h *UserRoleHandler) Create(c fiber.Ctx) error {
	var userRole models.UserRoleInput
	if err := c.Bind().JSON(&userRole); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, userRole)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed assigning role to user"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "User role assigned"))
}

// Read godoc
// @Summary      List user roles
// @Description  Get a paginated list of user role assignments
// @Tags         user-roles
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.UserRolePage}
// @Failure      401      {object}  common.JSONResponse
// @Failure      500      {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/user-roles [get]
func (h *UserRoleHandler) Read(c fiber.Ctx) error {
	userRoles, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(userRoles)
}

// Update godoc
// @Summary      Update user role assignment
// @Description  Update an existing user role mapping
// @Tags         user-roles
// @Accept       json
// @Produce      json
// @Param        id         path      string                true  "User Role ID"
// @Param        user_role  body      models.UserRoleInput  true  "Updated user role mapping details"
// @Success      200        {object}  common.JSONResponse
// @Failure      400        {object}  common.JSONResponse
// @Failure      401        {object}  common.JSONResponse
// @Failure      500        {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/user-roles/{id} [put]
func (h *UserRoleHandler) Update(c fiber.Ctx) error {
	var userRole models.UserRoleInput
	if err := c.Bind().JSON(&userRole); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), userRole); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating user role"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "User role updated"))
}

// Delete godoc
// @Summary      Delete user role assignment
// @Description  Delete a user role assignment by ID
// @Tags         user-roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User Role ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/user-roles/{id} [delete]
func (h *UserRoleHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting user role"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "User role deleted"))
}

// GetByID godoc
// @Summary      Get user role by ID
// @Description  Get detailed information about a user role assignment by ID
// @Tags         user-roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User Role ID"
// @Success      200  {object}  common.JSONResponse{items=entities.UserRole}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/user-roles/{id} [get]
func (h *UserRoleHandler) GetByID(c fiber.Ctx) error {
	userRole, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting user role"))
	}
	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(userRole, "User role found"))
}

func (h *UserRoleHandler) SetupRoutes(v1 fiber.Router) {
	userRoles := v1.Group("/user-roles")

	// Apply authentication middleware (Authn) to all user roles routes
	userRoles.Use(middlewares.Authn())

	userRoles.Post("/", middlewares.Authz("userRoles#create"), h.Create)
	userRoles.Get("/", middlewares.Authz("userRoles#read"), h.Read)
	userRoles.Get("/:id", middlewares.Authz("userRoles#read"), h.GetByID)
	userRoles.Put("/:id", middlewares.Authz("userRoles#update"), h.Update)
	userRoles.Delete("/:id", middlewares.Authz("userRoles#delete"), h.Delete)
}
