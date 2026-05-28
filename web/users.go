package web

import (
	"bytes"
	"strconv"

	"github.com/MarcelArt/gotel/internal/v1/models"
	"github.com/gofiber/fiber/v3"
)

type UserWebViewModel struct {
	ID       uint
	Username string
	Email    string
	Roles    []string
}

type UsersViewModel struct {
	BaseViewModel
	Users      []UserWebViewModel
	Pagination PaginationInfo
	Error      string
	Success    string
}

func (h *WebHandler) getUsersViewModel(c fiber.Ctx, userID any) (UsersViewModel, error) {
	currentUser, err := h.userService.GetByID(c, userID)
	if err != nil {
		return UsersViewModel{}, err
	}

	pageInfo, usersList := h.userService.Read(c)

	webUsers := make([]UserWebViewModel, len(usersList))
	for i, u := range usersList {
		rolesList, _ := u.Roles.Deserialize()
		webUsers[i] = UserWebViewModel{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Roles:    rolesList,
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

	return UsersViewModel{
		BaseViewModel: BaseViewModel{
			Title:       "Users Management - Gotel",
			ActiveTab:   "users",
			User:        currentUser,
			Permissions: getPermissions(c),
		},
		Users:      webUsers,
		Pagination: pagination,
	}, nil
}

// UsersGet handles GET /users requests.
func (h *WebHandler) UsersGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	vm, err := h.getUsersViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	return h.renderTab(c, "users", vm)
}

// UsersPost handles POST /users requests.
func (h *WebHandler) UsersPost(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	userInput := models.UserInput{
		Username: username,
		Email:    email,
		Password: password,
	}

	_, createErr := h.userService.Create(c, userInput)

	vm, err := h.getUsersViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if createErr != nil {
		vm.Error = "Failed to create user: " + createErr.Error()
	} else {
		vm.Success = "User created successfully!"
	}

	return h.renderTab(c, "users", vm)
}

// UsersEditGet handles GET /users/:id/edit requests.
func (h *WebHandler) UsersEditGet(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	user, err := h.userService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	userWeb := UserWebViewModel{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	t, ok := views["users"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found")
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "edit_user_modal", struct {
		User UserWebViewModel
	}{
		User: userWeb,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}

// UsersPut handles PUT /users/:id requests.
func (h *WebHandler) UsersPut(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	// 1. Assign roles
	// Read all role check values from the request
	// Temporarily override query size to get all roles to map checkboxes
	originalSize := string(c.Request().URI().QueryArgs().Peek("size"))
	c.Request().URI().QueryArgs().Set("size", "-1")
	_, allRoles := h.roleService.Read(c)
	if originalSize != "" {
		c.Request().URI().QueryArgs().Set("size", originalSize)
	} else {
		c.Request().URI().QueryArgs().Del("size")
	}

	var selectedRoleIDs []uint
	for _, r := range allRoles {
		if c.FormValue("role_"+strconv.FormatUint(uint64(r.ID), 10)) == "on" {
			selectedRoleIDs = append(selectedRoleIDs, r.ID)
		}
	}

	assignErr := h.userService.AssignRoles(c, uint(id), selectedRoleIDs)

	// 2. Update user profile (username, email)
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	var updateErr error
	if username != "" && email != "" {
		userInput := models.UserInput{
			Username: username,
			Email:    email,
			Password: password,
		}
		updateErr = h.userService.Update(c, uint(id), userInput)
	}

	vm, err := h.getUsersViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if assignErr != nil {
		vm.Error = "Failed to assign roles: " + assignErr.Error()
	} else if updateErr != nil {
		vm.Error = "Failed to update user details: " + updateErr.Error()
	} else {
		vm.Success = "User updated successfully!"
	}

	return h.renderTab(c, "users", vm)
}

// UsersDelete handles DELETE /users/:id requests.
func (h *WebHandler) UsersDelete(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	currentUser, _ := h.userService.GetByID(c, userID)
	if currentUser.ID == uint(id) {
		vm, err := h.getUsersViewModel(c, userID)
		if err != nil {
			return h.LogoutPost(c)
		}
		vm.Error = "You cannot delete your own account!"
		return h.renderTab(c, "users", vm)
	}

	deleteErr := h.userService.Delete(c, uint(id))

	vm, err := h.getUsersViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if deleteErr != nil {
		vm.Error = "Failed to delete user: " + deleteErr.Error()
	} else {
		vm.Success = "User deleted successfully!"
	}

	return h.renderTab(c, "users", vm)
}

// UserRolesListGet handles GET /users/roles/list requests.
func (h *WebHandler) UserRolesListGet(c fiber.Ctx) error {
	userIDStr := c.Query("userId")
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	search := c.Query("search")
	pageStr := c.Query("page")
	page, _ := strconv.ParseInt(pageStr, 10, 64)
	if page < 0 {
		page = 0
	}

	assignedMap := make(map[uint]bool)
	if userID > 0 {
		assignedRoles, err := h.userService.GetRoles(uint(userID))
		if err == nil {
			for _, ur := range assignedRoles {
				assignedMap[ur.ID] = true
			}
		}
	}

	originalFilters := string(c.Request().URI().QueryArgs().Peek("filters"))
	originalPage := string(c.Request().URI().QueryArgs().Peek("page"))
	originalSize := string(c.Request().URI().QueryArgs().Peek("size"))

	c.Request().URI().QueryArgs().Set("page", strconv.FormatInt(page, 10))
	c.Request().URI().QueryArgs().Set("size", "10")
	if search != "" {
		c.Request().URI().QueryArgs().Set("filters", `[["name", "like", "`+search+`"]]`)
	} else {
		c.Request().URI().QueryArgs().Del("filters")
	}

	pageInfo, rolesList := h.roleService.Read(c)

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

	type RoleCheckItem struct {
		ID          uint
		Name        string
		Description string
		Checked     bool
	}
	items := make([]RoleCheckItem, len(rolesList))
	for i, r := range rolesList {
		items[i] = RoleCheckItem{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Checked:     assignedMap[r.ID],
		}
	}

	prevPage := pageInfo.Page - 1
	if prevPage < 0 {
		prevPage = 0
	}

	viewModel := struct {
		Items      []RoleCheckItem
		UserID     uint
		Search     string
		Pagination PaginationInfo
	}{
		Items:  items,
		UserID: uint(userID),
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

	t, ok := views["users"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found")
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "roles_checklist_items", viewModel); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}
