package web

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/enums"
	"github.com/MarcelArt/gotel/internal/v1/features/roles"
	"github.com/MarcelArt/gotel/pkg/jsonb"
	"github.com/gofiber/fiber/v3"
	"gorm.io/datatypes"
)

type PermissionDefinition struct {
	Key         string
	Name        string
	Description string
}

var AvailablePermissions = []PermissionDefinition{
	{enums.PermFullAccess, "Full Access", "Grant all permissions to the operator"},
	{"roles#create", "Create Roles", "Allows creating new roles"},
	{"roles#read", "Read Roles", "Allows listing and viewing roles"},
	{"roles#update", "Update Roles", "Allows updating roles"},
	{"roles#delete", "Delete Roles", "Allows deleting roles"},
	{"users#create", "Create Users", "Allows creating new users"},
	{"users#read", "Read Users", "Allows listing and viewing users"},
	{"users#update", "Update Users", "Allows assigning roles to users"},
	{"users#delete", "Delete Users", "Allows deleting users"},
}

type RoleWebViewModel struct {
	ID          uint
	Name        string
	Description string
	Permissions []string
}

type PaginationInfo struct {
	Page        int64
	CurrentPage int64
	Size        int64
	TotalPages  int64
	Total       int64
	Last        bool
	First       bool
	NextPage    int64
	PrevPage    int64
}

type RolesViewModel struct {
	BaseViewModel
	Roles                []RoleWebViewModel
	AvailablePermissions []PermissionDefinition
	Error                string
	Success              string
	Pagination           PaginationInfo
}

type RoleEditViewModel struct {
	Role                 RoleWebViewModel
	AvailablePermissions []PermissionDefinition
}

func (h *WebHandler) getRolesViewModel(c fiber.Ctx, userID any) (RolesViewModel, error) {
	user, err := h.userService.GetByID(c, userID)
	if err != nil {
		return RolesViewModel{}, err
	}

	pageInfo, rolesList := h.roleService.Read(c)
	webRoles := make([]RoleWebViewModel, len(rolesList))
	for i, r := range rolesList {
		perms, _ := r.Permissions.Deserialize()
		webRoles[i] = RoleWebViewModel{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Permissions: perms,
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

	return RolesViewModel{
		BaseViewModel: BaseViewModel{
			Title:     "Roles Management - Gotel",
			ActiveTab: "roles",
			User:      user,
		},
		Roles:                webRoles,
		AvailablePermissions: AvailablePermissions,
		Pagination:           pagination,
	}, nil
}

// RolesGet handles GET /roles requests.
func (h *WebHandler) RolesGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	vm, err := h.getRolesViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	return h.renderTab(c, "roles", vm)
}

// RolesPost handles POST /roles requests.
func (h *WebHandler) RolesPost(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	name := c.FormValue("name")
	description := c.FormValue("description")

	var selectedPerms []string
	for _, p := range AvailablePermissions {
		if c.FormValue("permission_"+p.Key) == "on" {
			selectedPerms = append(selectedPerms, p.Key)
		}
	}

	permsBytes, err := json.Marshal(selectedPerms)
	var createErr error
	if err == nil {
		roleInput := roles.RoleInput{
			Name:        name,
			Description: description,
			Permissions: jsonb.JSONB[[]string]{
				JSON: datatypes.JSON(permsBytes),
			},
		}
		_, createErr = h.roleService.Create(c, roleInput)
	} else {
		createErr = err
	}

	vm, err := h.getRolesViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if createErr != nil {
		vm.Error = "Failed to create role: " + createErr.Error()
	} else {
		vm.Success = "Role created successfully!"
	}

	return h.renderTab(c, "roles", vm)
}

// RolesEditGet handles GET /roles/:id/edit requests.
func (h *WebHandler) RolesEditGet(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid role ID")
	}

	role, err := h.roleService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Role not found")
	}

	perms, _ := role.Permissions.Deserialize()
	roleWeb := RoleWebViewModel{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: perms,
	}

	vm := RoleEditViewModel{
		Role:                 roleWeb,
		AvailablePermissions: AvailablePermissions,
	}

	t, ok := views["roles"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found")
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "edit_role_modal", vm); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}

// RolesPut handles PUT /roles/:id requests.
func (h *WebHandler) RolesPut(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid role ID")
	}

	name := c.FormValue("name")
	description := c.FormValue("description")

	var selectedPerms []string
	for _, p := range AvailablePermissions {
		if c.FormValue("permission_"+p.Key) == "on" {
			selectedPerms = append(selectedPerms, p.Key)
		}
	}

	permsBytes, err := json.Marshal(selectedPerms)
	var updateErr error
	if err == nil {
		roleInput := roles.RoleInput{
			Name:        name,
			Description: description,
			Permissions: jsonb.JSONB[[]string]{
				JSON: datatypes.JSON(permsBytes),
			},
		}
		// Convert roleInput to entities.Role type check
		var _ entities.Role
		updateErr = h.roleService.Update(c, uint(id), roleInput)
	} else {
		updateErr = err
	}

	vm, err := h.getRolesViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if updateErr != nil {
		vm.Error = "Failed to update role: " + updateErr.Error()
	} else {
		vm.Success = "Role updated successfully!"
	}

	return h.renderTab(c, "roles", vm)
}

// RolesDelete handles DELETE /roles/:id requests.
func (h *WebHandler) RolesDelete(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid role ID")
	}

	deleteErr := h.roleService.Delete(c, uint(id))

	vm, err := h.getRolesViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if deleteErr != nil {
		vm.Error = "Failed to delete role: " + deleteErr.Error()
	} else {
		vm.Success = "Role deleted successfully!"
	}

	return h.renderTab(c, "roles", vm)
}
