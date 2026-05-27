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

type AssetInstanceHandler struct {
	service services.IAssetInstanceService
}

var _ = entities.AssetInstance{}

func NewAssetInstanceHandler(service services.IAssetInstanceService) *AssetInstanceHandler {
	return &AssetInstanceHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new asset instance
// @Description  Create a new asset instance with the provided details
// @Tags         asset-instances
// @Accept       json
// @Produce      json
// @Param        assetInstance  body      models.AssetInstanceInput  true  "Asset instance details"
// @Success      201            {object}  common.JSONResponse{items=uint}
// @Failure      400            {object}  common.JSONResponse
// @Failure      500            {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/asset-instances [post]
func (h *AssetInstanceHandler) Create(c fiber.Ctx) error {
	var instance models.AssetInstanceInput
	if err := c.Bind().JSON(&instance); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, instance)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed creating asset instance"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Asset instance created"))
}

// Read godoc
// @Summary      List asset instances
// @Description  Get a paginated list of asset instances
// @Tags         asset-instances
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.AssetInstancePage}
// @Failure      401      {object}  common.JSONResponse
// @Failure      500      {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/asset-instances [get]
func (h *AssetInstanceHandler) Read(c fiber.Ctx) error {
	instances, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(instances)
}

// Update godoc
// @Summary      Update asset instance
// @Description  Update an existing asset instance's details
// @Tags         asset-instances
// @Accept       json
// @Produce      json
// @Param        id             path      string                     true  "Asset Instance ID"
// @Param        assetInstance  body      models.AssetInstanceInput  true  "Updated asset instance details"
// @Success      200            {object}  common.JSONResponse
// @Failure      400            {object}  common.JSONResponse
// @Failure      401            {object}  common.JSONResponse
// @Failure      500            {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/asset-instances/{id} [put]
func (h *AssetInstanceHandler) Update(c fiber.Ctx) error {
	var instance models.AssetInstanceInput
	if err := c.Bind().JSON(&instance); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), instance); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating asset instance"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Asset instance updated"))
}

// Delete godoc
// @Summary      Delete asset instance
// @Description  Delete an asset instance by ID
// @Tags         asset-instances
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Asset Instance ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/asset-instances/{id} [delete]
func (h *AssetInstanceHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting asset instance"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Asset instance deleted"))
}

// GetByID godoc
// @Summary      Get asset instance by ID
// @Description  Get detailed information about an asset instance by its ID
// @Tags         asset-instances
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Asset Instance ID"
// @Success      200  {object}  common.JSONResponse{items=entities.AssetInstance}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/asset-instances/{id} [get]
func (h *AssetInstanceHandler) GetByID(c fiber.Ctx) error {
	instance, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting asset instance"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(instance, "Asset instance found"))
}

func (h *AssetInstanceHandler) SetupRoutes(v1 fiber.Router) {
	instances := v1.Group("/asset-instances")

	instances.Use(middlewares.Authn())

	instances.Post("/", middlewares.Authz("asset_instances#create"), h.Create)
	instances.Get("/", middlewares.Authz("asset_instances#read"), h.Read)
	instances.Get("/:id", middlewares.Authz("asset_instances#read"), h.GetByID)
	instances.Put("/:id", middlewares.Authz("asset_instances#update"), h.Update)
	instances.Delete("/:id", middlewares.Authz("asset_instances#delete"), h.Delete)
}
