package web

import (
	"bytes"
	"fmt"
	"html"
	"strconv"

	"github.com/MarcelArt/gotel/internal/v1/features/locations"
	"github.com/gofiber/fiber/v3"
)

type LocationWebViewModel struct {
	ID          uint
	Value       string
	Description string
	IsVirtual   bool
}

type LocationsViewModel struct {
	BaseViewModel
	Locations  []LocationWebViewModel
	Pagination PaginationInfo
	Error      string
	Success    string
}

func (h *WebHandler) getLocationsViewModel(c fiber.Ctx, userID any) (LocationsViewModel, error) {
	currentUser, err := h.userService.GetByID(c, userID)
	if err != nil {
		return LocationsViewModel{}, err
	}

	pageInfo, locationsList := h.locationService.Read(c)

	webLocations := make([]LocationWebViewModel, len(locationsList))
	for i, loc := range locationsList {
		webLocations[i] = LocationWebViewModel{
			ID:          loc.ID,
			Value:       loc.Value,
			Description: loc.Description,
			IsVirtual:   loc.IsVirtual,
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

	return LocationsViewModel{
		BaseViewModel: BaseViewModel{
			Title:       "Locations Directory - Gotel",
			ActiveTab:   "locations",
			User:        currentUser,
			Permissions: getPermissions(c),
		},
		Locations:  webLocations,
		Pagination: pagination,
	}, nil
}

// LocationsGet handles GET /locations requests.
func (h *WebHandler) LocationsGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	vm, err := h.getLocationsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	return h.renderTab(c, "locations", vm)
}

// LocationsPost handles POST /locations requests.
func (h *WebHandler) LocationsPost(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	value := c.FormValue("value")
	description := c.FormValue("description")
	isVirtual := c.FormValue("isVirtual") == "on"

	input := locations.LocationInput{
		Value:       value,
		Description: description,
		IsVirtual:   isVirtual,
	}

	_, createErr := h.locationService.Create(c, input)

	vm, err := h.getLocationsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if createErr != nil {
		vm.Error = "Failed to create location: " + createErr.Error()
	} else {
		vm.Success = "Location created successfully!"
	}

	return h.renderTab(c, "locations", vm)
}

// LocationsEditGet handles GET /locations/:id/edit requests.
func (h *WebHandler) LocationsEditGet(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid location ID")
	}

	location, err := h.locationService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Location not found")
	}

	vm := struct {
		Location LocationWebViewModel
	}{
		Location: LocationWebViewModel{
			ID:          location.ID,
			Value:       location.Value,
			Description: location.Description,
			IsVirtual:   location.IsVirtual,
		},
	}

	t, ok := views["locations"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found")
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "edit_location_modal", vm); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}

// LocationsPut handles PUT /locations/:id requests.
func (h *WebHandler) LocationsPut(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid location ID")
	}

	value := c.FormValue("value")
	description := c.FormValue("description")
	isVirtual := c.FormValue("isVirtual") == "on"

	input := locations.LocationInput{
		Value:       value,
		Description: description,
		IsVirtual:   isVirtual,
	}

	updateErr := h.locationService.Update(c, uint(id), input)

	vm, err := h.getLocationsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if updateErr != nil {
		vm.Error = "Failed to update location: " + updateErr.Error()
	} else {
		vm.Success = "Location updated successfully!"
	}

	return h.renderTab(c, "locations", vm)
}

// LocationsDelete handles DELETE /locations/:id requests.
func (h *WebHandler) LocationsDelete(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid location ID")
	}

	deleteErr := h.locationService.Delete(c, uint(id))

	vm, err := h.getLocationsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if deleteErr != nil {
		vm.Error = "Failed to delete location: " + deleteErr.Error()
	} else {
		vm.Success = "Location deleted successfully!"
	}

	return h.renderTab(c, "locations", vm)
}

// LocationsOptionsGet handles GET /locations/options requests.
// It returns dropdown item buttons for infinite scrolling of locations.
func (h *WebHandler) LocationsOptionsGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	pageInfo, locationsList := h.locationService.Read(c)

	var buf bytes.Buffer
	for _, loc := range locationsList {
		escapedVal := html.EscapeString(loc.Value)
		buf.WriteString(fmt.Sprintf(`<button type="button" class="dropdown-item" @click="selectedId = '%d'; selectedName = '%s'; open = false;">%s</button>`, loc.ID, escapedVal, escapedVal))
	}

	if !pageInfo.Last {
		nextPage := pageInfo.Page + 1
		buf.WriteString(fmt.Sprintf(`
		<div hx-get="/locations/options?page=%d&size=%d"
			 hx-trigger="intersect once"
			 hx-target="this"
			 hx-swap="outerHTML"
			 style="padding: 8px; text-align: center; color: var(--color-secondary); font-size: 12px;">
			 Loading more...
		</div>`, nextPage, pageInfo.Size))
	}

	c.Type("html")
	return c.Send(buf.Bytes())
}
