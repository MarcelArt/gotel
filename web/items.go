package web

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MarcelArt/gotel/internal/v1/features/items"
	"github.com/gofiber/fiber/v3"
)

type ItemWebViewModel struct {
	ID           uint
	Code         string
	Name         string
	Picture      string
	TrackingMode string
	Unit         string
	CategoryID   uint
	CategoryName string
}

type ItemsViewModel struct {
	BaseViewModel
	Items      []ItemWebViewModel
	Pagination PaginationInfo
	Error      string
	Success    string
}

func (h *WebHandler) getItemsViewModel(c fiber.Ctx, userID any) (ItemsViewModel, error) {
	currentUser, err := h.userService.GetByID(c, userID)
	if err != nil {
		return ItemsViewModel{}, err
	}

	pageInfo, itemsList := h.itemService.Read(c)

	categoryMap := make(map[uint]string)
	webItems := make([]ItemWebViewModel, len(itemsList))
	for i, it := range itemsList {
		categoryName := ""
		if it.CategoryID != 0 {
			if name, exists := categoryMap[it.CategoryID]; exists {
				categoryName = name
			} else {
				if cat, err := h.categoryService.GetByID(c, it.CategoryID); err == nil {
					categoryMap[it.CategoryID] = cat.Value
					categoryName = cat.Value
				}
			}
		}

		webItems[i] = ItemWebViewModel{
			ID:           it.ID,
			Code:         it.Code,
			Name:         it.Name,
			Picture:      it.Picture,
			TrackingMode: it.TrackingMode,
			Unit:         it.Unit,
			CategoryID:   it.CategoryID,
			CategoryName: categoryName,
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

	return ItemsViewModel{
		BaseViewModel: BaseViewModel{
			Title:     "Items Catalog - Gotel",
			ActiveTab: "items",
			User:      currentUser,
		},
		Items:      webItems,
		Pagination: pagination,
	}, nil
}

// ItemsGet handles GET /items requests.
func (h *WebHandler) ItemsGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	vm, err := h.getItemsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	return h.renderTab(c, "items", vm)
}

// ItemsPost handles POST /items requests.
func (h *WebHandler) ItemsPost(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	code := c.FormValue("code")
	name := c.FormValue("name")
	trackingMode := c.FormValue("trackingMode")
	unit := c.FormValue("unit")
	categoryIDStr := c.FormValue("categoryId")

	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		vm, _ := h.getItemsViewModel(c, userID)
		vm.Error = "Invalid Category Selected"
		return h.renderTab(c, "items", vm)
	}

	input := items.ItemInput{
		Code:         code,
		Name:         name,
		TrackingMode: trackingMode,
		Unit:         unit,
		CategoryID:   uint(categoryID),
	}

	// File upload logic matching items/handler.go (clean path convention)
	file, _ := c.FormFile("file")
	if file != nil {
		today := time.Now().Unix()
		basePath := "public/uploads"
		if err := os.MkdirAll(basePath, 0755); err != nil {
			vm, _ := h.getItemsViewModel(c, userID)
			vm.Error = "Failed to create upload directory: " + err.Error()
			return h.renderTab(c, "items", vm)
		}

		filename := fmt.Sprintf("/%s/item-%d-%s", basePath, today, file.Filename)
		if err := c.SaveFile(file, fmt.Sprintf(".%s", filename)); err != nil {
			vm, _ := h.getItemsViewModel(c, userID)
			vm.Error = "Failed to save uploaded file: " + err.Error()
			return h.renderTab(c, "items", vm)
		}

		input.Picture = filename
	}

	_, createErr := h.itemService.Create(c, input)

	vm, err := h.getItemsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if createErr != nil {
		vm.Error = "Failed to create item: " + createErr.Error()
	} else {
		vm.Success = "Item created successfully!"
	}

	return h.renderTab(c, "items", vm)
}

// ItemsEditGet handles GET /items/:id/edit requests.
func (h *WebHandler) ItemsEditGet(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid item ID")
	}

	item, err := h.itemService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Item not found")
	}

	categoryName := ""
	if item.CategoryID != 0 {
		if cat, err := h.categoryService.GetByID(c, item.CategoryID); err == nil {
			categoryName = cat.Value
		}
	}

	vm := struct {
		Item ItemWebViewModel
	}{
		Item: ItemWebViewModel{
			ID:           item.ID,
			Code:         item.Code,
			Name:         item.Name,
			Picture:      item.Picture,
			TrackingMode: item.TrackingMode,
			Unit:         item.Unit,
			CategoryID:   item.CategoryID,
			CategoryName: categoryName,
		},
	}

	t, ok := views["items"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found")
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "edit_item_modal", vm); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}

// ItemsPut handles PUT /items/:id requests.
func (h *WebHandler) ItemsPut(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid item ID")
	}

	existingItem, err := h.itemService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Item not found")
	}

	code := c.FormValue("code")
	name := c.FormValue("name")
	trackingMode := c.FormValue("trackingMode")
	unit := c.FormValue("unit")
	categoryIDStr := c.FormValue("categoryId")

	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		vm, _ := h.getItemsViewModel(c, userID)
		vm.Error = "Invalid Category Selected"
		return h.renderTab(c, "items", vm)
	}

	input := items.ItemInput{
		Code:         code,
		Name:         name,
		TrackingMode: trackingMode,
		Unit:         unit,
		CategoryID:   uint(categoryID),
		Picture:      existingItem.Picture, // retain original picture if no new file is uploaded
	}

	// File upload logic matching items/handler.go
	file, _ := c.FormFile("file")
	if file != nil {
		today := time.Now().Unix()
		basePath := "public/uploads"
		if err := os.MkdirAll(basePath, 0755); err != nil {
			vm, _ := h.getItemsViewModel(c, userID)
			vm.Error = "Failed to create upload directory: " + err.Error()
			return h.renderTab(c, "items", vm)
		}

		filename := fmt.Sprintf("/%s/item-%d-%s", basePath, today, file.Filename)
		if err := c.SaveFile(file, fmt.Sprintf(".%s", filename)); err != nil {
			vm, _ := h.getItemsViewModel(c, userID)
			vm.Error = "Failed to save uploaded file: " + err.Error()
			return h.renderTab(c, "items", vm)
		}

		input.Picture = filename
	}

	updateErr := h.itemService.Update(c, uint(id), input)

	vm, err := h.getItemsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if updateErr != nil {
		vm.Error = "Failed to update item: " + updateErr.Error()
	} else {
		vm.Success = "Item updated successfully!"
	}

	return h.renderTab(c, "items", vm)
}

// ItemsDelete handles DELETE /items/:id requests.
func (h *WebHandler) ItemsDelete(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid item ID")
	}

	deleteErr := h.itemService.Delete(c, uint(id))

	vm, err := h.getItemsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if deleteErr != nil {
		vm.Error = "Failed to delete item: " + deleteErr.Error()
	} else {
		vm.Success = "Item deleted successfully!"
	}

	return h.renderTab(c, "items", vm)
}
