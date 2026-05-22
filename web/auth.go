package web

import (
	"fmt"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/v1/features/users"
	"github.com/gofiber/fiber/v3"
)

// LoginViewModel represents the data required to render the login view.
type LoginViewModel struct {
	Title    string
	Error    string
	Success  string
	Username string
}

// RegisterViewModel represents the data required to render the register view.
type RegisterViewModel struct {
	Title    string
	Error    string
	Username string
	Email    string
}

// LoginGet handles GET /login requests.
func (h *WebHandler) LoginGet(c fiber.Ctx) error {
	// If already logged in, redirect to dashboard
	atCookie := c.Cookies("at")
	if atCookie != "" {
		if _, err := common.ParseToken(atCookie); err == nil {
			return c.Redirect().To("/")
		}
	}

	successMsg := ""
	if c.Query("registered") == "true" {
		successMsg = "Account created successfully! Please sign in below."
	}

	vm := LoginViewModel{
		Title:   "Sign In - Gotel",
		Success: successMsg,
	}

	return h.render(c, "login", vm)
}

// LoginPost handles POST /login requests.
func (h *WebHandler) LoginPost(c fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	isRemember := c.FormValue("isRemember") == "on"

	_, err := h.userService.Login(c, users.LoginInput{
		Username:   username,
		Password:   password,
		IsRemember: isRemember,
	})

	if err != nil {
		errorHTML := `
		<div class="alert">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
			<span>Invalid username or password</span>
		</div>`

		if c.Get("HX-Request") == "true" {
			return c.Status(fiber.StatusUnauthorized).SendString(errorHTML)
		}

		vm := LoginViewModel{
			Title:    "Sign In - Gotel",
			Error:    "Invalid username or password",
			Username: username,
		}

		return h.render(c, "login", vm)
	}

	if c.Get("HX-Request") == "true" {
		c.Set("HX-Redirect", "/")
		return c.SendString("Success")
	}

	return c.Redirect().To("/")
}

// RegisterGet handles GET /register requests.
func (h *WebHandler) RegisterGet(c fiber.Ctx) error {
	// If already logged in, redirect to dashboard
	atCookie := c.Cookies("at")
	if atCookie != "" {
		if _, err := common.ParseToken(atCookie); err == nil {
			return c.Redirect().To("/")
		}
	}

	vm := RegisterViewModel{
		Title: "Sign Up - Gotel",
	}

	return h.render(c, "register", vm)
}

// RegisterPost handles POST /register requests.
func (h *WebHandler) RegisterPost(c fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	_, err := h.userService.Create(c, users.UserInput{
		Username: username,
		Email:    email,
		Password: password,
	})

	if err != nil {
		errorHTML := fmt.Sprintf(`
		<div class="alert">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
			<span>Registration failed: %s</span>
		</div>`, err.Error())

		if c.Get("HX-Request") == "true" {
			return c.Status(fiber.StatusBadRequest).SendString(errorHTML)
		}

		vm := RegisterViewModel{
			Title:    "Sign Up - Gotel",
			Error:    err.Error(),
			Username: username,
			Email:    email,
		}

		return h.render(c, "register", vm)
	}

	if c.Get("HX-Request") == "true" {
		c.Set("HX-Redirect", "/login?registered=true")
		return c.SendString("Success")
	}

	return c.Redirect().To("/login?registered=true")
}

// LogoutPost handles POST /logout requests.
func (h *WebHandler) LogoutPost(c fiber.Ctx) error {
	// Set cookies max age to -1 to clear them
	c.Cookie(&fiber.Cookie{
		Name:     "at",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "rt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
	})

	return c.Redirect().To("/login")
}
