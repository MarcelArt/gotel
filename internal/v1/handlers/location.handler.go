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

type LocationHandler struct {
	service services.ILocationService
}

var _ = entities.Location{}

func NewLocationHandler(service services.ILocationService) *LocationHandler {
	return &LocationHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new location
// @Description  Create a new location with the provided details
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        location  body      models.LocationInput  true  "Location details"
// @Success      201       {object}  common.JSONResponse{items=uint}
// @Failure      400       {object}  common.JSONResponse
// @Failure      500       {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/locations [post]
func (h *LocationHandler) Create(c fiber.Ctx) error {
	var location models.LocationInput
	if err := c.Bind().JSON(&location); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, location)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed creating location"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Location created"))
}

// Read godoc
// @Summary      List locations
// @Description  Get a paginated list of locations
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.LocationPage}
// @Failure      401      {object}  common.JSONResponse
// @Failure      500      {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/locations [get]
func (h *LocationHandler) Read(c fiber.Ctx) error {
	locations, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(locations)
}

// Update godoc
// @Summary      Update location
// @Description  Update an existing location's details
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        id        path      string                true  "Location ID"
// @Param        location  body      models.LocationInput  true  "Updated location details"
// @Success      200       {object}  common.JSONResponse
// @Failure      400       {object}  common.JSONResponse
// @Failure      401       {object}  common.JSONResponse
// @Failure      500       {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/locations/{id} [put]
func (h *LocationHandler) Update(c fiber.Ctx) error {
	var location models.LocationInput
	if err := c.Bind().JSON(&location); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), location); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating location"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Location updated"))
}

// Delete godoc
// @Summary      Delete location
// @Description  Delete a location by ID
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Location ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/locations/{id} [delete]
func (h *LocationHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting location"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Location deleted"))
}

// GetByID godoc
// @Summary      Get location by ID
// @Description  Get detailed information about a location by its ID
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Location ID"
// @Success      200  {object}  common.JSONResponse{items=entities.Location}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/locations/{id} [get]
func (h *LocationHandler) GetByID(c fiber.Ctx) error {
	location, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting location"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(location, "Location found"))
}

func (h *LocationHandler) SetupRoutes(v1 fiber.Router) {
	locations := v1.Group("/locations")

	locations.Use(middlewares.Authn())

	locations.Post("/", middlewares.Authz("locations#create"), h.Create)
	locations.Get("/", middlewares.Authz("locations#read"), h.Read)
	locations.Get("/:id", middlewares.Authz("locations#read"), h.GetByID)
	locations.Put("/:id", middlewares.Authz("locations#update"), h.Update)
	locations.Delete("/:id", middlewares.Authz("locations#delete"), h.Delete)
}
