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

type AssetTransactionHandler struct {
	service services.IAssetTransactionService
}

var _ = entities.AssetTransaction{}

func NewAssetTransactionHandler(service services.IAssetTransactionService) *AssetTransactionHandler {
	return &AssetTransactionHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new asset transaction
// @Description  Create a new asset transaction with the provided details
// @Tags         asset-transactions
// @Accept       json
// @Produce      json
// @Param        assetTransaction  body      models.AssetTransactionInput  true  "Asset transaction details"
// @Success      201               {object}  common.JSONResponse{items=uint}
// @Failure      400               {object}  common.JSONResponse
// @Failure      500               {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/asset-transactions [post]
func (h *AssetTransactionHandler) Create(c fiber.Ctx) error {
	var tx models.AssetTransactionInput
	if err := c.Bind().JSON(&tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, tx)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed creating asset transaction"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Asset transaction created"))
}

// Read godoc
// @Summary      List asset transactions
// @Description  Get a paginated list of asset transactions
// @Tags         asset-transactions
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.AssetTransactionPage}
// @Failure      401      {object}  common.JSONResponse
// @Failure      500      {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/asset-transactions [get]
func (h *AssetTransactionHandler) Read(c fiber.Ctx) error {
	txs, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(txs)
}

// Update godoc
// @Summary      Update asset transaction
// @Description  Update an existing asset transaction's details
// @Tags         asset-transactions
// @Accept       json
// @Produce      json
// @Param        id                path      string                        true  "Asset Transaction ID"
// @Param        assetTransaction  body      models.AssetTransactionInput  true  "Updated asset transaction details"
// @Success      200               {object}  common.JSONResponse
// @Failure      400               {object}  common.JSONResponse
// @Failure      401               {object}  common.JSONResponse
// @Failure      500               {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/asset-transactions/{id} [put]
func (h *AssetTransactionHandler) Update(c fiber.Ctx) error {
	var tx models.AssetTransactionInput
	if err := c.Bind().JSON(&tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), tx); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating asset transaction"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Asset transaction updated"))
}

// Delete godoc
// @Summary      Delete asset transaction
// @Description  Delete an asset transaction by ID
// @Tags         asset-transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Asset Transaction ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/asset-transactions/{id} [delete]
func (h *AssetTransactionHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting asset transaction"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Asset transaction deleted"))
}

// GetByID godoc
// @Summary      Get asset transaction by ID
// @Description  Get detailed information about an asset transaction by its ID
// @Tags         asset-transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Asset Transaction ID"
// @Success      200  {object}  common.JSONResponse{items=entities.AssetTransaction}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/asset-transactions/{id} [get]
func (h *AssetTransactionHandler) GetByID(c fiber.Ctx) error {
	tx, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting asset transaction"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(tx, "Asset transaction found"))
}

func (h *AssetTransactionHandler) SetupRoutes(v1 fiber.Router) {
	txs := v1.Group("/asset-transactions")

	txs.Use(middlewares.Authn())

	txs.Post("/", middlewares.Authz("assetTransactions#create"), h.Create)
	txs.Get("/", middlewares.Authz("assetTransactions#read"), h.Read)
	txs.Get("/:id", middlewares.Authz("assetTransactions#read"), h.GetByID)
	txs.Put("/:id", middlewares.Authz("assetTransactions#update"), h.Update)
	txs.Delete("/:id", middlewares.Authz("assetTransactions#delete"), h.Delete)
}
