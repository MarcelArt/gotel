package web

import (
	"bytes"
	"strconv"

	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/gofiber/fiber/v3"
)

type RoomWebViewModel struct {
	ID         uint
	RoomNumber string
	Floor      string
	Status     string
}

type RoomsViewModel struct {
	BaseViewModel
	Rooms      []RoomWebViewModel
	Pagination PaginationInfo
	Error      string
	Success    string
}

func (h *WebHandler) getRoomsViewModel(c fiber.Ctx, userID any) (RoomsViewModel, error) {
	currentUser, err := h.userService.GetByID(c, userID)
	if err != nil {
		return RoomsViewModel{}, err
	}

	pageInfo, roomsList := h.roomService.Read(c)

	webRooms := make([]RoomWebViewModel, len(roomsList))
	for i, r := range roomsList {
		webRooms[i] = RoomWebViewModel{
			ID:         r.ID,
			RoomNumber: r.RoomNumber,
			Floor:      r.Floor,
			Status:     r.Status,
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

	return RoomsViewModel{
		BaseViewModel: BaseViewModel{
			Title:       "Rooms Directory - Gotel",
			ActiveTab:   "rooms",
			User:        currentUser,
			Permissions: getPermissions(c),
		},
		Rooms:      webRooms,
		Pagination: pagination,
	}, nil
}

// RoomsGet handles GET /rooms requests.
func (h *WebHandler) RoomsGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	vm, err := h.getRoomsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	return h.renderTab(c, "rooms", vm)
}

// RoomsPost handles POST /rooms requests.
func (h *WebHandler) RoomsPost(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	roomNumber := c.FormValue("roomNumber")
	floor := c.FormValue("floor")
	status := c.FormValue("status")

	input := models.RoomInput{
		RoomNumber: roomNumber,
		Floor:      floor,
		Status:     status,
	}

	_, createErr := h.roomService.Create(c, input)

	c.Request().Header.SetMethod(fiber.MethodGet)
	vm, err := h.getRoomsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if createErr != nil {
		vm.Error = "Failed to create room: " + createErr.Error()
	} else {
		vm.Success = "Room created successfully!"
	}

	return h.renderTab(c, "rooms", vm)
}

// RoomsEditGet handles GET /rooms/:id/edit requests.
func (h *WebHandler) RoomsEditGet(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid room ID")
	}

	room, err := h.roomService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Room not found")
	}

	vm := struct {
		Room RoomWebViewModel
	}{
		Room: RoomWebViewModel{
			ID:         room.ID,
			RoomNumber: room.RoomNumber,
			Floor:      room.Floor,
			Status:     room.Status,
		},
	}

	t, ok := views["rooms"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found")
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "edit_room_modal", vm); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}

// RoomsPut handles PUT /rooms/:id requests.
func (h *WebHandler) RoomsPut(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid room ID")
	}

	_, err = h.roomService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Room not found")
	}

	roomNumber := c.FormValue("roomNumber")
	floor := c.FormValue("floor")
	status := c.FormValue("status")

	input := models.RoomInput{
		RoomNumber: roomNumber,
		Floor:      floor,
		Status:     status,
	}

	updateErr := h.roomService.Update(c, uint(id), input)

	c.Request().Header.SetMethod(fiber.MethodGet)
	vm, err := h.getRoomsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if updateErr != nil {
		vm.Error = "Failed to update room: " + updateErr.Error()
	} else {
		vm.Success = "Room updated successfully!"
	}

	return h.renderTab(c, "rooms", vm)
}

// RoomsDelete handles DELETE /rooms/:id requests.
func (h *WebHandler) RoomsDelete(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid room ID")
	}

	_, err = h.roomService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Room not found")
	}

	deleteErr := h.roomService.Delete(c, uint(id))

	c.Request().Header.SetMethod(fiber.MethodGet)
	vm, err := h.getRoomsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if deleteErr != nil {
		vm.Error = "Failed to delete room: " + deleteErr.Error()
	} else {
		vm.Success = "Room deleted successfully!"
	}

	return h.renderTab(c, "rooms", vm)
}
