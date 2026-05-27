package handlers

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/middlewares"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/MarcelArt/gotel/internal/v1/services"
	"github.com/gofiber/fiber/v3"
	_ "github.com/morkid/paginate"
)

type UserHandler struct {
	service services.IUserService
}

var _ = entities.User{}

func NewUserHandler(service services.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new user
// @Description  Create a new user with the provided details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserInput  true  "User details"
// @Success      201   {object}  common.JSONResponse{items=uint}
// @Failure      400   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Router       /v1/users [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	var user models.UserInput
	if err := c.Bind().JSON(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, user)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed creating user"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "User created"))
}

// Read godoc
// @Summary      List users
// @Description  Get a paginated list of users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.UserPage}
// @Failure      401      {object}  common.JSONResponse
// @Failure      500      {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/users [get]
func (h *UserHandler) Read(c fiber.Ctx) error {
	users, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(users)
}

// Update godoc
// @Summary      Update user
// @Description  Update an existing user's details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string            true  "User ID"
// @Param        user  body      models.UserInput  true  "Updated user details"
// @Success      200   {object}  common.JSONResponse
// @Failure      400   {object}  common.JSONResponse
// @Failure      401   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/users/{id} [put]
func (h *UserHandler) Update(c fiber.Ctx) error {
	var user models.UserInput
	if err := c.Bind().JSON(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), user); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating user"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "User updated"))
}

// Delete godoc
// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/users/{id} [delete]
func (h *UserHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting user"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "User deleted"))
}

// GetByID godoc
// @Summary      Get user by ID
// @Description  Get detailed information about a user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  common.JSONResponse{items=entities.User}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/users/{id} [get]
func (h *UserHandler) GetByID(c fiber.Ctx) error {
	user, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting user"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(user, "User found"))
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate a user with username/email and password, returning tokens
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.LoginInput  true  "Login credentials"
// @Success      200          {object}  common.JSONResponse{items=models.LoginResponse}
// @Failure      400          {object}  common.JSONResponse
// @Failure      401          {object}  common.JSONResponse
// @Router       /v1/users/login [post]
func (h *UserHandler) Login(c fiber.Ctx) error {
	var input models.LoginInput
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	res, err := h.service.Login(c, input)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(common.NewJSONResponse(err, "invalid username or password"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(res, "User logged in"))
}

// Refresh godoc
// @Summary      Refresh token
// @Description  Regenerate access and refresh token pair using the refresh token from header
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        X-Refresh-Token  header    string  true  "Refresh token"
// @Success      200              {object}  common.JSONResponse{items=models.LoginResponse}
// @Failure      400              {object}  common.JSONResponse
// @Failure      401              {object}  common.JSONResponse
// @Failure      500              {object}  common.JSONResponse
// @Router       /v1/users/refresh [post]
func (h *UserHandler) Refresh(c fiber.Ctx) error {
	claims := common.FiberCtxToClaims(c)
	id := claims["userId"]

	isRemember, ok := claims["isRemember"].(bool)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(fmt.Errorf("missing isRemember claim"), "invalid refresh token"))
	}

	res, err := h.service.RegenerateTokenPair(c, id, isRemember)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed generating tokens"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(res, "Authenticated"))
}

// AssignRoles godoc
// @Summary      Assign roles to user
// @Description  Assign a list of role IDs to a user by user ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "User ID"
// @Param        roleIDs  body      []uint  true  "List of Role IDs"
// @Success      200      {object}  common.JSONResponse
// @Failure      400      {object}  common.JSONResponse
// @Failure      401      {object}  common.JSONResponse
// @Failure      500      {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/users/{id}/roles [patch]
func (h *UserHandler) AssignRoles(c fiber.Ctx) error {
	id := fiber.Params[uint](c, "id")
	var roleIDs []uint
	if err := c.Bind().JSON(&roleIDs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.AssignRoles(c, id, roleIDs); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed assigning roles"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Roles assigned"))
}

// GetPermissions godoc
// @Summary      Get user permissions
// @Description  Get list of permissions for the authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.JSONResponse{items=[]string}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/users/permissions [get]
func (h *UserHandler) GetPermissions(c fiber.Ctx) error {
	claims := common.FiberCtxToClaims(c)
	id := claims["userId"]

	permissions, err := h.service.GetPermissions(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed retrieving permissions"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(permissions, "Permissions found"))
}

func (h *UserHandler) SetupRoutes(v1 fiber.Router) {
	users := v1.Group("/users")

	users.Post("/", h.Create)
	users.Post("/login", h.Login)
	users.Post("/refresh", middlewares.Refresh(), h.Refresh)

	users.Get("/", middlewares.Authn(), h.Read)
	users.Get("/permissions", middlewares.Authn(), h.GetPermissions)
	users.Get("/:id", middlewares.Authn(), h.GetByID)

	users.Put("/:id", middlewares.Authn(), h.Update)

	users.Patch("/:id/roles", middlewares.Authn(), middlewares.Authz("users#update"), h.AssignRoles)

	users.Delete("/:id", middlewares.Authn(), middlewares.Authz("users#delete"), h.Delete)
}
