package web

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/gofiber/fiber/v3"
)

type AssetInstanceWebViewModel struct {
	ID       uint
	Code     string
	ItemID   uint
	Status   string
	Location string
	Note     string
}

type AssetInstancesViewModel struct {
	BaseViewModel
	Instances   []AssetInstanceWebViewModel
	Pagination  PaginationInfo
	ItemID      uint
	ItemName    string
	ItemCode    string
	ItemUnit    string
	ItemPicture string
	Error       string
	Success     string
}

func (h *WebHandler) getAssetInstancesViewModel(c fiber.Ctx, userID any, itemID uint) (AssetInstancesViewModel, error) {
	currentUser, err := h.userService.GetByID(c, userID)
	if err != nil {
		return AssetInstancesViewModel{}, err
	}

	item, err := h.itemService.GetByID(c, itemID)
	if err != nil {
		return AssetInstancesViewModel{}, err
	}

	itemIDStr := strconv.FormatUint(uint64(itemID), 10)
	filtersJSON := fmt.Sprintf(`[["item_id", "=", "%s"]]`, itemIDStr)
	c.Request().URI().QueryArgs().Set("filters", filtersJSON)

	pageInfo, instancesList := h.assetInstanceService.Read(c)

	webInstances := make([]AssetInstanceWebViewModel, len(instancesList))
	for i, inst := range instancesList {
		webInstances[i] = AssetInstanceWebViewModel{
			ID:       inst.ID,
			Code:     inst.Code,
			ItemID:   inst.ItemID,
			Status:   inst.Status,
			Location: inst.Location,
			Note:     inst.Note,
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

	return AssetInstancesViewModel{
		BaseViewModel: BaseViewModel{
			Title:       "Asset Instances: " + item.Name + " - Gotel",
			ActiveTab:   "items",
			User:        currentUser,
			Permissions: getPermissions(c),
		},
		Instances:   webInstances,
		Pagination:  pagination,
		ItemID:      item.ID,
		ItemName:    item.Name,
		ItemCode:    item.Code,
		ItemUnit:    item.Unit,
		ItemPicture: item.Picture,
	}, nil
}

// AssetInstancesGet handles GET /asset-instances requests.
func (h *WebHandler) AssetInstancesGet(c fiber.Ctx) error {
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

	vm, err := h.getAssetInstancesViewModel(c, userID, uint(itemID))
	if err != nil {
		return h.LogoutPost(c)
	}

	return h.renderTab(c, "asset_instances", vm)
}

// AssetInstancesPost handles POST /asset-instances requests.
func (h *WebHandler) AssetInstancesPost(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	itemIDStr := c.FormValue("itemId")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Item ID")
	}

	code := c.FormValue("code")

	input := models.AssetInstanceInput{
		Code:   code,
		ItemID: uint(itemID),
	}

	_, createErr := h.assetInstanceService.Create(c, input)

	c.Request().Header.SetMethod(fiber.MethodGet)
	vm, err := h.getAssetInstancesViewModel(c, userID, uint(itemID))
	if err != nil {
		return h.LogoutPost(c)
	}

	if createErr != nil {
		vm.Error = "Failed to create asset instance: " + createErr.Error()
	} else {
		vm.Success = "Asset instance created successfully!"
	}

	return h.renderTab(c, "asset_instances", vm)
}

// AssetInstancesEditGet handles GET /asset-instances/:id/edit requests.
func (h *WebHandler) AssetInstancesEditGet(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid asset instance ID")
	}

	instance, err := h.assetInstanceService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Asset instance not found")
	}

	vm := struct {
		Instance AssetInstanceWebViewModel
	}{
		Instance: AssetInstanceWebViewModel{
			ID:     instance.ID,
			Code:   instance.Code,
			ItemID: instance.ItemID,
		},
	}

	t, ok := views["asset_instances"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found")
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "edit_instance_modal", vm); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}

// AssetInstancesPut handles PUT /asset-instances/:id requests.
func (h *WebHandler) AssetInstancesPut(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid asset instance ID")
	}

	existingInstance, err := h.assetInstanceService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Asset instance not found")
	}

	code := c.FormValue("code")

	input := models.AssetInstanceInput{
		Code:   code,
		ItemID: existingInstance.ItemID,
	}

	updateErr := h.assetInstanceService.Update(c, uint(id), input)

	c.Request().Header.SetMethod(fiber.MethodGet)
	vm, err := h.getAssetInstancesViewModel(c, userID, existingInstance.ItemID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if updateErr != nil {
		vm.Error = "Failed to update asset instance: " + updateErr.Error()
	} else {
		vm.Success = "Asset instance updated successfully!"
	}

	return h.renderTab(c, "asset_instances", vm)
}

// AssetInstancesDelete handles DELETE /asset-instances/:id requests.
func (h *WebHandler) AssetInstancesDelete(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid asset instance ID")
	}

	existingInstance, err := h.assetInstanceService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Asset instance not found")
	}

	deleteErr := h.assetInstanceService.Delete(c, uint(id))

	c.Request().Header.SetMethod(fiber.MethodGet)
	vm, err := h.getAssetInstancesViewModel(c, userID, existingInstance.ItemID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if deleteErr != nil {
		vm.Error = "Failed to delete asset instance: " + deleteErr.Error()
	} else {
		vm.Success = "Asset instance deleted successfully!"
	}

	return h.renderTab(c, "asset_instances", vm)
}
