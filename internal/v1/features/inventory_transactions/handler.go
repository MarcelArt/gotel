package inventory_transactions

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/v1/middlewares"
	"github.com/gofiber/fiber/v3"
	_ "github.com/morkid/paginate"
)

type InventoryTransactionHandler struct {
	service IInventoryTransactionService
}

func NewInventoryTransactionHandler(service IInventoryTransactionService) *InventoryTransactionHandler {
	return &InventoryTransactionHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new inventory transaction
// @Description  Create a new inventory transaction with the provided details
// @Tags         inventory-transactions
// @Accept       json
// @Produce      json
// @Param        transaction  body      InventoryTransactionInput  true  "Inventory transaction details"
// @Success      201          {object}  common.JSONResponse{items=uint}
// @Failure      400          {object}  common.JSONResponse
// @Failure      500          {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/inventory-transactions [post]
func (h *InventoryTransactionHandler) Create(c fiber.Ctx) error {
	var tx InventoryTransactionInput
	if err := c.Bind().JSON(&tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.service.Create(c, tx)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed creating inventory transaction"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Inventory transaction created"))
}

// Read godoc
// @Summary      List inventory transactions
// @Description  Get a paginated list of inventory transactions
// @Tags         inventory-transactions
// @Accept       json
// @Produce      json
// @Param			page		query		int		false	"Page"
// @Param			size		query		int		false	"Size"
// @Param			sort		query		string	false	"Sort"
// @Param			filters		query		string	false	"Filter"
// @Success      200    {object}  paginate.Page{items=[]InventoryTransactionPage}
// @Failure      401    {object}  common.JSONResponse
// @Failure      500    {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/inventory-transactions [get]
func (h *InventoryTransactionHandler) Read(c fiber.Ctx) error {
	txs, _ := h.service.Read(c)

	return c.Status(fiber.StatusOK).JSON(txs)
}

// Update godoc
// @Summary      Update inventory transaction
// @Description  Update an existing inventory transaction's details
// @Tags         inventory-transactions
// @Accept       json
// @Produce      json
// @Param        id           path      string                     true  "Transaction ID"
// @Param        transaction  body      InventoryTransactionInput  true  "Updated transaction details"
// @Success      200          {object}  common.JSONResponse
// @Failure      400          {object}  common.JSONResponse
// @Failure      401          {object}  common.JSONResponse
// @Failure      500          {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/inventory-transactions/{id} [put]
func (h *InventoryTransactionHandler) Update(c fiber.Ctx) error {
	var tx InventoryTransactionInput
	if err := c.Bind().JSON(&tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.service.Update(c, c.Params("id"), tx); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed updating inventory transaction"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Inventory transaction updated"))
}

// Delete godoc
// @Summary      Delete inventory transaction
// @Description  Delete an inventory transaction by ID
// @Tags         inventory-transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/inventory-transactions/{id} [delete]
func (h *InventoryTransactionHandler) Delete(c fiber.Ctx) error {
	if err := h.service.Delete(c, c.Params("id")); err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed deleting inventory transaction"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Inventory transaction deleted"))
}

// GetByID godoc
// @Summary      Get inventory transaction by ID
// @Description  Get detailed information about an inventory transaction by its ID
// @Tags         inventory-transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  common.JSONResponse{items=entities.InventoryTransaction}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/inventory-transactions/{id} [get]
func (h *InventoryTransactionHandler) GetByID(c fiber.Ctx) error {
	tx, err := h.service.GetByID(c, c.Params("id"))
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting inventory transaction"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(tx, "Inventory transaction found"))
}

// GetItemCounts godoc
// @Summary      Get item counts
// @Description  Get aggregated transaction counts (by transaction type) for a specific item
// @Tags         inventory-transactions
// @Accept       json
// @Produce      json
// @Param        item_id  path      int  true  "Item ID"
// @Success      200      {object}  common.JSONResponse{items=[]ItemCount}
// @Failure      400      {object}  common.JSONResponse
// @Failure      401      {object}  common.JSONResponse
// @Failure      500      {object}  common.JSONResponse
// @Security     ApiKeyAuth
// @Router       /v1/inventory-transactions/item-counts/{item_id} [get]
func (h *InventoryTransactionHandler) GetItemCounts(c fiber.Ctx) error {
	itemID := fiber.Params[uint](c, "item_id")

	counts, err := h.service.GetItemCounts(itemID)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting item counts"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(counts, "Item counts found"))
}

func (h *InventoryTransactionHandler) SetupRoutes(v1 fiber.Router) {
	txs := v1.Group("/inventory-transactions")

	txs.Use(middlewares.Authn())

	txs.Post("/", middlewares.Authz("inventoryTransactions#create"), h.Create)

	txs.Get("/", middlewares.Authz("inventoryTransactions#read"), h.Read)
	txs.Get("/item-counts/:item_id", middlewares.Authz("inventoryTransactions#read"), h.GetItemCounts)
	txs.Get("/:id", middlewares.Authz("inventoryTransactions#read"), h.GetByID)

	txs.Put("/:id", middlewares.Authz("inventoryTransactions#update"), h.Update)
	txs.Delete("/:id", middlewares.Authz("inventoryTransactions#delete"), h.Delete)
}
