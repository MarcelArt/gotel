package users

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/gofiber/fiber/v3"
	_ "github.com/morkid/paginate"
)

type UserHandler struct {
	service IUserService
}

func NewUserHandler(service IUserService) *UserHandler {
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
// @Param        user  body      UserInput  true  "User details"
// @Success      201   {object}  common.JSONResponse{items=uint}
// @Failure      400   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Router       /v1/users [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	var user UserInput
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
// @Param			page		query		int		false	"Page"
// @Param			size		query		int		false	"Size"
// @Param			sort		query		string	false	"Sort"
// @Param			filters		query		string	false	"Filter"
// @Success      200    {object}  paginate.Page{items=[]UserPage}
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
// @Param        id    path      string     true  "User ID"
// @Param        user  body      UserInput  true  "Updated user details"
// @Success      200   {object}  common.JSONResponse
// @Failure      400   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Router       /v1/users/{id} [put]
func (h *UserHandler) Update(c fiber.Ctx) error {
	var user UserInput
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
// @Failure      500  {object}  common.JSONResponse
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
// @Failure      500  {object}  common.JSONResponse
// @Router       /v1/users/{id} [get]
func (h *UserHandler) GetByID(c fiber.Ctx) error {
	user, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting user"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(user, "User found"))
}

func (h *UserHandler) SetupRoutes(v1 fiber.Router) {
	users := v1.Group("/users")

	users.Post("/", h.Create)
	users.Get("/", h.Read)
	users.Get("/:id", h.GetByID)
	users.Put("/:id", h.Update)
	users.Delete("/:id", h.Delete)
}
