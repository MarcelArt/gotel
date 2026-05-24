package web

import (
	"strconv"

	"github.com/MarcelArt/gotel/internal/v1/features/inventory_transactions"
	"github.com/gofiber/fiber/v3"
)

type InventoryTransactionWebViewModel struct {
	ID              uint
	CreatedAt       string
	TransactionType string
	Quantity        float64
	Note            string
	ItemID          uint
	Item            string
	Unit            string
	Actor           string
	From            string
	To              string
	Route           string
}

type InventoryTransactionsViewModel struct {
	BaseViewModel
	Transactions []InventoryTransactionWebViewModel
	Pagination   PaginationInfo
	ItemID       uint
	ItemName     string
	ItemCode     string
	ItemUnit     string
	ItemPicture  string
	Error        string
	Success      string
}

func (h *WebHandler) getInventoryTransactionsViewModel(c fiber.Ctx, userID any, itemID uint) (InventoryTransactionsViewModel, error) {
	currentUser, err := h.userService.GetByID(c, userID)
	if err != nil {
		return InventoryTransactionsViewModel{}, err
	}

	item, err := h.itemService.GetByID(c, itemID)
	if err != nil {
		return InventoryTransactionsViewModel{}, err
	}

	pageInfo, txsList := h.inventoryTransactionService.Read(c)

	webTxs := make([]InventoryTransactionWebViewModel, len(txsList))
	for i, tx := range txsList {
		route := ""
		if tx.TransactionType == "RECEIVE" || tx.TransactionType == "LAUNDRY_IN" {
			route = "To: " + tx.To
		} else if tx.TransactionType == "ISSUE" || tx.TransactionType == "LAUNDRY_OUT" || tx.TransactionType == "DAMAGE" || tx.TransactionType == "LOST" || tx.TransactionType == "CONSUME" || tx.TransactionType == "DISPOSE" {
			route = "From: " + tx.From
		} else if tx.TransactionType == "TRANSFER" {
			route = tx.From + " → " + tx.To
		} else { // OTHER or unspecified
			if tx.From != "" && tx.To != "" {
				route = tx.From + " → " + tx.To
			} else if tx.From != "" {
				route = "From: " + tx.From
			} else if tx.To != "" {
				route = "To: " + tx.To
			} else {
				route = "-"
			}
		}

		webTxs[i] = InventoryTransactionWebViewModel{
			ID:              tx.ID,
			CreatedAt:       tx.CreatedAt.Format("02 Jan 2006, 15:04:05"),
			TransactionType: tx.TransactionType,
			Quantity:        tx.Quantity,
			Note:            tx.Note,
			ItemID:          tx.ItemID,
			Item:            tx.Item,
			Unit:            tx.Unit,
			Actor:           tx.Actor,
			From:            tx.From,
			To:              tx.To,
			Route:           route,
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

	return InventoryTransactionsViewModel{
		BaseViewModel: BaseViewModel{
			Title:       "Stock Ledger: " + item.Name + " - Gotel",
			ActiveTab:   "items", // Keeps Items sidebar highlighted
			User:        currentUser,
			Permissions: getPermissions(c),
		},
		Transactions: webTxs,
		Pagination:   pagination,
		ItemID:       item.ID,
		ItemName:     item.Name,
		ItemCode:     item.Code,
		ItemUnit:     item.Unit,
		ItemPicture:  item.Picture,
	}, nil
}

// InventoryTransactionsGet handles GET /inventory-transactions requests.
func (h *WebHandler) InventoryTransactionsGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	itemIDStr := c.Query("itemId")
	if itemIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Item ID is required")
	}

	itemID, err := strconv.ParseUint(itemIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Item ID")
	}

	// Set GORM filter and sorting programmatically before calling Read
	c.Request().URI().QueryArgs().Set("filters", `[["item_id", "=", "`+itemIDStr+`"]]`)
	c.Request().URI().QueryArgs().Set("sort", "-id")

	vm, err := h.getInventoryTransactionsViewModel(c, userID, uint(itemID))
	if err != nil {
		return h.LogoutPost(c)
	}

	return h.renderTab(c, "transactions", vm)
}

// InventoryTransactionsPost handles POST /inventory-transactions requests.
func (h *WebHandler) InventoryTransactionsPost(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	itemIDStr := c.FormValue("itemId")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Item ID")
	}

	transactionType := c.FormValue("transactionType")
	quantityStr := c.FormValue("quantity")
	note := c.FormValue("note")
	fromIdStr := c.FormValue("fromId")
	toIdStr := c.FormValue("toId")

	quantity, err := strconv.ParseFloat(quantityStr, 64)
	if err != nil {
		c.Request().Header.SetMethod(fiber.MethodGet)
		c.Request().URI().QueryArgs().Set("filters", `[["item_id", "=", "`+itemIDStr+`"]]`)
		vm, _ := h.getInventoryTransactionsViewModel(c, userID, uint(itemID))
		vm.Error = "Invalid quantity value"
		return h.renderTab(c, "transactions", vm)
	}

	var fromID *uint
	if fromIdStr != "" && fromIdStr != "0" {
		if id, err := strconv.ParseUint(fromIdStr, 10, 64); err == nil {
			uid := uint(id)
			fromID = &uid
		}
	}

	var toID *uint
	if toIdStr != "" && toIdStr != "0" {
		if id, err := strconv.ParseUint(toIdStr, 10, 64); err == nil {
			uid := uint(id)
			toID = &uid
		}
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

	input := inventory_transactions.InventoryTransactionInput{
		TransactionType: transactionType,
		Quantity:        quantity,
		Note:            note,
		ItemID:          uint(itemID),
		FromID:          fromID,
		ToID:            toID,
		ActorID:         actorID,
	}

	_, createErr := h.inventoryTransactionService.Create(c, input)

	c.Request().Header.SetMethod(fiber.MethodGet)
	c.Request().URI().QueryArgs().Set("filters", `[["item_id", "=", "`+itemIDStr+`"]]`)
	c.Request().URI().QueryArgs().Set("sort", "-id")
	vm, err := h.getInventoryTransactionsViewModel(c, userID, uint(itemID))
	if err != nil {
		return h.LogoutPost(c)
	}

	if createErr != nil {
		vm.Error = "Failed to record transaction: " + createErr.Error()
	} else {
		vm.Success = "Inventory transaction recorded successfully!"
	}

	return h.renderTab(c, "transactions", vm)
}
