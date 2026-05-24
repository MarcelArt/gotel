package web

import (
	"bytes"
	"fmt"
	"html"
	"strconv"

	"github.com/MarcelArt/gotel/internal/v1/features/categories"
	"github.com/gofiber/fiber/v3"
)

type CategoryWebViewModel struct {
	ID          uint
	Value       string
	Description string
}

type CategoriesViewModel struct {
	BaseViewModel
	Categories []CategoryWebViewModel
	Pagination PaginationInfo
	Error      string
	Success    string
}

func (h *WebHandler) getCategoriesViewModel(c fiber.Ctx, userID any) (CategoriesViewModel, error) {
	currentUser, err := h.userService.GetByID(c, userID)
	if err != nil {
		return CategoriesViewModel{}, err
	}

	pageInfo, categoriesList := h.categoryService.Read(c)

	webCategories := make([]CategoryWebViewModel, len(categoriesList))
	for i, cat := range categoriesList {
		webCategories[i] = CategoryWebViewModel{
			ID:          cat.ID,
			Value:       cat.Value,
			Description: cat.Description,
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

	return CategoriesViewModel{
		BaseViewModel: BaseViewModel{
			Title:       "Categories Management - Gotel",
			ActiveTab:   "categories",
			User:        currentUser,
			Permissions: getPermissions(c),
		},
		Categories: webCategories,
		Pagination: pagination,
	}, nil
}

// CategoriesGet handles GET /categories requests.
func (h *WebHandler) CategoriesGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	vm, err := h.getCategoriesViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	return h.renderTab(c, "categories", vm)
}

// CategoriesPost handles POST /categories requests.
func (h *WebHandler) CategoriesPost(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	value := c.FormValue("value")
	description := c.FormValue("description")

	input := categories.CategoryInput{
		Value:       value,
		Description: description,
	}

	_, createErr := h.categoryService.Create(c, input)

	vm, err := h.getCategoriesViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if createErr != nil {
		vm.Error = "Failed to create category: " + createErr.Error()
	} else {
		vm.Success = "Category created successfully!"
	}

	return h.renderTab(c, "categories", vm)
}

// CategoriesEditGet handles GET /categories/:id/edit requests.
func (h *WebHandler) CategoriesEditGet(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid category ID")
	}

	category, err := h.categoryService.GetByID(c, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Category not found")
	}

	vm := struct {
		Category CategoryWebViewModel
	}{
		Category: CategoryWebViewModel{
			ID:          category.ID,
			Value:       category.Value,
			Description: category.Description,
		},
	}

	t, ok := views["categories"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Template not found")
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "edit_category_modal", vm); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Type("html")
	return c.Send(buf.Bytes())
}

// CategoriesPut handles PUT /categories/:id requests.
func (h *WebHandler) CategoriesPut(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid category ID")
	}

	value := c.FormValue("value")
	description := c.FormValue("description")

	input := categories.CategoryInput{
		Value:       value,
		Description: description,
	}

	updateErr := h.categoryService.Update(c, uint(id), input)

	vm, err := h.getCategoriesViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if updateErr != nil {
		vm.Error = "Failed to update category: " + updateErr.Error()
	} else {
		vm.Success = "Category updated successfully!"
	}

	return h.renderTab(c, "categories", vm)
}

// CategoriesDelete handles DELETE /categories/:id requests.
func (h *WebHandler) CategoriesDelete(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid category ID")
	}

	deleteErr := h.categoryService.Delete(c, uint(id))

	vm, err := h.getCategoriesViewModel(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	if deleteErr != nil {
		vm.Error = "Failed to delete category: " + deleteErr.Error()
	} else {
		vm.Success = "Category deleted successfully!"
	}

	return h.renderTab(c, "categories", vm)
}

// CategoriesOptionsGet handles GET /categories/options requests.
// It returns dropdown item buttons for infinite scrolling of categories.
func (h *WebHandler) CategoriesOptionsGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	pageInfo, categoriesList := h.categoryService.Read(c)

	var buf bytes.Buffer
	for _, cat := range categoriesList {
		escapedVal := html.EscapeString(cat.Value)
		buf.WriteString(fmt.Sprintf(`<button type="button" class="dropdown-item" @click="selectedCategoryId = '%d'; selectedCategoryName = '%s'; open = false;">%s</button>`, cat.ID, escapedVal, escapedVal))
	}

	if !pageInfo.Last {
		nextPage := pageInfo.Page + 1
		buf.WriteString(fmt.Sprintf(`
		<div hx-get="/categories/options?page=%d&size=%d"
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
