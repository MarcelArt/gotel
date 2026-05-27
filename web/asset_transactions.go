package web

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/gofiber/fiber/v3"
)

type AssetTransactionWebViewModel struct {
	ID              uint
	CreatedAt       string
	TransactionType string
	Status          string
	Note            string
	LocationID      uint
	Location        string
	Actor           string
}

type AssetTransactionsViewModel struct {
	BaseViewModel
	Transactions []AssetTransactionWebViewModel
	Pagination   PaginationInfo
	ItemID       uint
	ItemName     string
	ItemCode     string
	ItemPicture  string
	InstanceID   uint
	InstanceCode string
	Error        string
	Success      string
}

func (h *WebHandler) getAssetTransactionsViewModel(c fiber.Ctx, userID any, instanceID uint) (AssetTransactionsViewModel, error) {
	currentUser, err := h.userService.GetByID(c, userID)
	if err != nil {
		return AssetTransactionsViewModel{}, err
	}

	instance, err := h.assetInstanceService.GetByID(c, instanceID)
	if err != nil {
		return AssetTransactionsViewModel{}, err
	}

	item, err := h.itemService.GetByID(c, instance.ItemID)
	if err != nil {
		return AssetTransactionsViewModel{}, err
	}

	instanceIDStr := strconv.FormatUint(uint64(instanceID), 10)
	filtersJSON := fmt.Sprintf(`[["instance_id", "=", "%s"]]`, instanceIDStr)
	c.Request().URI().QueryArgs().Set("filters", filtersJSON)
	c.Request().URI().QueryArgs().Set("sort", "-id")

	pageInfo, txsList := h.assetTransactionService.Read(c)

	webTxs := make([]AssetTransactionWebViewModel, len(txsList))
	for i, tx := range txsList {
		webTxs[i] = AssetTransactionWebViewModel{
			ID:              tx.ID,
			CreatedAt:       tx.CreatedAt.Format("02 Jan 2006, 15:04:05"),
			TransactionType: tx.TransactionType,
			Status:          tx.Status,
			Note:            tx.Note,
			LocationID:      tx.LocationID,
			Location:        tx.Location,
			Actor:           tx.Actor,
		}
	}

	prevPage := pageInfo.Page - 1
	if prevPage < 0 {
		prevPage = 0
	}

	pagination := PaginationInfo{
		Page:        pageInfo.Page,
		CurrentPage: pageInfo.Page + 1,
		Size:        pageInfo.Size,
		TotalPages:  pageInfo.TotalPages,
		Total:       pageInfo.Total,
		Last:        pageInfo.Last,
		First:       pageInfo.First,
		NextPage:    pageInfo.Page + 1,
		PrevPage:    prevPage,
	}

	return AssetTransactionsViewModel{
		BaseViewModel: BaseViewModel{
			Title:       "Asset Ledger: " + instance.Code + " - Gotel",
			ActiveTab:   "items",
			User:        currentUser,
			Permissions: getPermissions(c),
		},
		Transactions: webTxs,
		Pagination:   pagination,
		ItemID:       item.ID,
		ItemName:     item.Name,
		ItemCode:     item.Code,
		ItemPicture:  item.Picture,
		InstanceID:   instance.ID,
		InstanceCode: instance.Code,
	}, nil
}

// AssetTransactionsGet handles GET /asset-transactions requests.
func (h *WebHandler) AssetTransactionsGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	instanceIDStr := c.Query("instanceId")
	if instanceIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Instance ID is required")
	}

	instanceID, err := strconv.ParseUint(instanceIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Instance ID")
	}

	vm, err := h.getAssetTransactionsViewModel(c, userID, uint(instanceID))
	if err != nil {
		return h.LogoutPost(c)
	}

	return h.renderTab(c, "asset_transactions", vm)
}

// AssetTransactionsPost handles POST /asset-transactions requests.
func (h *WebHandler) AssetTransactionsPost(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	instanceIDStr := c.FormValue("instanceId")
	instanceID, err := strconv.ParseUint(instanceIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Instance ID")
	}

	transactionType := c.FormValue("transactionType")
	status := c.FormValue("status")
	note := c.FormValue("note")
	locationIDStr := c.FormValue("locationId")

	locationID, err := strconv.ParseUint(locationIDStr, 10, 64)
	if err != nil {
		c.Request().Header.SetMethod(fiber.MethodGet)
		vm, _ := h.getAssetTransactionsViewModel(c, userID, uint(instanceID))
		vm.Error = "Invalid location selected"
		return h.renderTab(c, "asset_transactions", vm)
	}

	var actorID uint
	if val, ok := userID.(float64); ok {
		actorID = uint(val)
	} else if val, ok := userID.(uint); ok {
		actorID = val
	} else if val, ok := userID.(int); ok {
		actorID = uint(val)
	} else if valStr, ok := userID.(string); ok {
		if id, err := strconv.ParseUint(valStr, 10, 64); err == nil {
			actorID = uint(id)
		}
	}

	input := models.AssetTransactionInput{
		TransactionType: transactionType,
		Status:          status,
		Note:            note,
		LocationID:      uint(locationID),
		InstanceID:      uint(instanceID),
		ActorID:         actorID,
	}

	_, createErr := h.assetTransactionService.Create(c, input)

	c.Request().Header.SetMethod(fiber.MethodGet)
	vm, err := h.getAssetTransactionsViewModel(c, userID, uint(instanceID))
	if err != nil {
		return h.LogoutPost(c)
	}

	if createErr != nil {
		vm.Error = "Failed to record transaction: " + createErr.Error()
	} else {
		vm.Success = "Asset transaction recorded successfully!"
	}

	return h.renderTab(c, "asset_transactions", vm)
}

// AssetTransactionsEditGet handles GET /asset-transactions/:id/edit requests.
func (h *WebHandler) AssetTransactionsEditGet(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid transaction ID")
	}

	tx, err := h.assetTransactionService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Transaction not found")
	}

	location, err := h.locationService.GetByID(c, tx.LocationID)
	locationValue := ""
	if err == nil {
		locationValue = location.Value
	}

	vm := struct {
		Transaction                  entities.AssetTransaction
		LocationValue                string
	}{
		Transaction:                  tx,
		LocationValue:                locationValue,
	}

	t, ok := views["asset_transactions"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found")
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "edit_transaction_modal", vm); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}

// AssetTransactionsPut handles PUT /asset-transactions/:id requests.
func (h *WebHandler) AssetTransactionsPut(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid transaction ID")
	}

	existingTx, err := h.assetTransactionService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Transaction not found")
	}

	transactionType := c.FormValue("transactionType")
	status := c.FormValue("status")
	note := c.FormValue("note")
	locationIDStr := c.FormValue("locationId")

	locationID, err := strconv.ParseUint(locationIDStr, 10, 64)
	if err != nil {
		c.Request().Header.SetMethod(fiber.MethodGet)
		vm, _ := h.getAssetTransactionsViewModel(c, userID, existingTx.InstanceID)
		vm.Error = "Invalid location selected"
		return h.renderTab(c, "asset_transactions", vm)
	}

	input := models.AssetTransactionInput{
		TransactionType: transactionType,
		Status:          status,
		Note:            note,
		LocationID:      uint(locationID),
		InstanceID:      existingTx.InstanceID,
		ActorID:         existingTx.ActorID,
	}

	updateErr := h.assetTransactionService.Update(c, uint(id), input)

	c.Request().Header.SetMethod(fiber.MethodGet)
	vm, err := h.getAssetTransactionsViewModel(c, userID, existingTx.InstanceID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if updateErr != nil {
		vm.Error = "Failed to update transaction: " + updateErr.Error()
	} else {
		vm.Success = "Asset transaction updated successfully!"
	}

	return h.renderTab(c, "asset_transactions", vm)
}

// AssetTransactionsDelete handles DELETE /asset-transactions/:id requests.
func (h *WebHandler) AssetTransactionsDelete(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid transaction ID")
	}

	existingTx, err := h.assetTransactionService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Transaction not found")
	}

	deleteErr := h.assetTransactionService.Delete(c, uint(id))

	c.Request().Header.SetMethod(fiber.MethodGet)
	vm, err := h.getAssetTransactionsViewModel(c, userID, existingTx.InstanceID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if deleteErr != nil {
		vm.Error = "Failed to delete transaction: " + deleteErr.Error()
	} else {
		vm.Success = "Asset transaction deleted successfully!"
	}

	return h.renderTab(c, "asset_transactions", vm)
}
