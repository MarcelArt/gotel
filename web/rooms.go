package web

import (
	"bytes"
	"strconv"
	"time"

	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/gofiber/fiber/v3"
)

type RoomWebViewModel struct {
	ID            uint
	RoomNumber    string
	Floor         string
	Status        string
	TaskID        uint
	TaskStartedAt *time.Time
	AssigneeID    uint
	Assignee      string
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
			ID:            r.ID,
			RoomNumber:    r.RoomNumber,
			Floor:         r.Floor,
			Status:        r.Status,
			TaskID:        r.TaskID,
			TaskStartedAt: r.TaskStartedAt,
			AssigneeID:    r.AssigneeID,
			Assignee:      r.Assignee,
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

// RoomsAssignCleaningPost handles POST /rooms/assign-cleaning requests.
func (h *WebHandler) RoomsAssignCleaningPost(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	var assignerID uint
	if idVal, ok := userID.(float64); ok {
		assignerID = uint(idVal)
	} else if idVal, ok := userID.(uint); ok {
		assignerID = idVal
	}

	roomIdStr := c.FormValue("roomId")
	roomId, err := strconv.ParseUint(roomIdStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Room ID")
	}

	priorityStr := c.FormValue("priority")
	priority, err := strconv.ParseUint(priorityStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Priority")
	}

	assigneeIdStr := c.FormValue("assigneeId")
	assigneeId, err := strconv.ParseUint(assigneeIdStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Assignee ID")
	}

	note := c.FormValue("note")

	input := models.HousekeepingTaskInput{
		RoomID:     uint(roomId),
		Priority:   uint(priority),
		AssigneeID: uint(assigneeId),
		AssignerID: assignerID,
		Note:       note,
	}

	assignErr := h.roomService.AssignCleaning(c, input)

	c.Request().Header.SetMethod(fiber.MethodGet)
	vm, err := h.getRoomsViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if assignErr != nil {
		vm.Error = "Failed to assign cleaning task: " + assignErr.Error()
	} else {
		vm.Success = "Cleaning task assigned successfully!"
	}

	return h.renderTab(c, "rooms", vm)
}

// UserDropdownListGet handles GET /users/dropdown/list requests.
func (h *WebHandler) UserDropdownListGet(c fiber.Ctx) error {
	search := c.Query("search")
	pageStr := c.Query("page")
	page, _ := strconv.ParseInt(pageStr, 10, 64)
	if page < 0 {
		page = 0
	}

	originalFilters := string(c.Request().URI().QueryArgs().Peek("filters"))
	originalPage := string(c.Request().URI().QueryArgs().Peek("page"))
	originalSize := string(c.Request().URI().QueryArgs().Peek("size"))

	c.Request().URI().QueryArgs().Set("page", strconv.FormatInt(page, 10))
	c.Request().URI().QueryArgs().Set("size", "10")
	if search != "" {
		c.Request().URI().QueryArgs().Set("filters", `[["username", "like", "`+search+`"]]`)
	} else {
		c.Request().URI().QueryArgs().Del("filters")
	}

	pageInfo, usersList := h.userService.Read(c)

	// Restore original query args
	if originalFilters != "" {
		c.Request().URI().QueryArgs().Set("filters", originalFilters)
	} else {
		c.Request().URI().QueryArgs().Del("filters")
	}
	if originalPage != "" {
		c.Request().URI().QueryArgs().Set("page", originalPage)
	} else {
		c.Request().URI().QueryArgs().Del("page")
	}
	if originalSize != "" {
		c.Request().URI().QueryArgs().Set("size", originalSize)
	} else {
		c.Request().URI().QueryArgs().Del("size")
	}

	type UserItem struct {
		ID       uint
		Username string
		Email    string
	}
	items := make([]UserItem, len(usersList))
	for i, u := range usersList {
		items[i] = UserItem{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		}
	}

	prevPage := pageInfo.Page - 1
	if prevPage < 0 {
		prevPage = 0
	}

	viewModel := struct {
		Items      []UserItem
		Search     string
		Pagination PaginationInfo
	}{
		Items:  items,
		Search: search,
		Pagination: PaginationInfo{
			Page:       pageInfo.Page,
			Size:       pageInfo.Size,
			TotalPages: pageInfo.TotalPages,
			Total:      pageInfo.Total,
			Last:       pageInfo.Last,
			First:      pageInfo.First,
			NextPage:   pageInfo.Page + 1,
			PrevPage:   prevPage,
		},
	}

	t, ok := views["rooms"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found")
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "users_dropdown_items", viewModel); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}
