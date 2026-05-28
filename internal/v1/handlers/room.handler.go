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

type RoomHandler struct {
	service services.IRoomService
}

var _ = entities.Room{}

func NewRoomHandler(service services.IRoomService) *RoomHandler {
	return &RoomHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new room
// @Description  Create a new room with the provided details
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        room  body      models.RoomInput  true  "Room details"
// @Success      201   {object}  common.JSONResponse{items=uint}
// @Failure      400   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/rooms [post]
func (h *RoomHandler) Create(c fiber.Ctx) error {
	var room models.RoomInput
	if err := c.Bind().JSON(&room); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, room)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed creating room"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Room created"))
}

// Read godoc
// @Summary      List rooms
// @Description  Get a paginated list of rooms
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.RoomPage}
// @Failure      401      {object}  common.JSONResponse
// @Failure      500      {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/rooms [get]
func (h *RoomHandler) Read(c fiber.Ctx) error {
	rooms, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(rooms)
}

// Update godoc
// @Summary      Update room
// @Description  Update an existing room's details
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id    path      string            true  "Room ID"
// @Param        room  body      models.RoomInput  true  "Updated room details"
// @Success      200   {object}  common.JSONResponse
// @Failure      400   {object}  common.JSONResponse
// @Failure      401   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/rooms/{id} [put]
func (h *RoomHandler) Update(c fiber.Ctx) error {
	var room models.RoomInput
	if err := c.Bind().JSON(&room); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), room); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating room"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Room updated"))
}

// Delete godoc
// @Summary      Delete room
// @Description  Delete a room by ID
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Room ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/rooms/{id} [delete]
func (h *RoomHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting room"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Room deleted"))
}

// GetByID godoc
// @Summary      Get room by ID
// @Description  Get detailed information about a room by its ID
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Room ID"
// @Success      200  {object}  common.JSONResponse{items=entities.Room}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/rooms/{id} [get]
func (h *RoomHandler) GetByID(c fiber.Ctx) error {
	room, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting room"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(room, "Room found"))
}

func (h *RoomHandler) SetupRoutes(v1 fiber.Router) {
	rooms := v1.Group("/rooms")

	rooms.Use(middlewares.Authn())

	rooms.Post("/", middlewares.Authz("rooms#create"), h.Create)
	rooms.Get("/", middlewares.Authz("rooms#read"), h.Read)
	rooms.Get("/:id", middlewares.Authz("rooms#read"), h.GetByID)
	rooms.Put("/:id", middlewares.Authz("rooms#update"), h.Update)
	rooms.Delete("/:id", middlewares.Authz("rooms#delete"), h.Delete)
}
