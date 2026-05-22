package web

import (
	"strings"

	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/v1/features/users"
	"github.com/gofiber/fiber/v3"
)

func WebAuth(userService users.IUserService) fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()

		// Skip authentication checks for static assets
		if strings.HasPrefix(path, "/public/") {
			return c.Next()
		}

		// Skip authentication checks for auth endpoints themselves
		if path == "/login" || path == "/register" {
			return c.Next()
		}

		atCookie := c.Cookies("at")
		rtCookie := c.Cookies("rt")

		if atCookie != "" {
			claims, err := common.ParseToken(atCookie)
			if err == nil {
				// Access token is valid
				c.Locals("userId", claims["userId"])
				c.Locals("username", claims["sub"])
				return c.Next()
			}
		}

		// Access token expired/invalid, try refresh token
		if rtCookie != "" {
			claims, err := common.ParseToken(rtCookie)
			if err == nil {
				// Refresh token is valid, regenerate the token pair
				userID := claims["userId"]
				isRemember := false
				if ir, ok := claims["isRemember"].(bool); ok {
					isRemember = ir
				}

				// This service method automatically writes the new cookies to the response
				res, err := userService.RegenerateTokenPair(c, userID, isRemember)
				if err == nil {
					// Parse the new access token to set context locals
					newClaims, err := common.ParseToken(res.AccessToken)
					if err == nil {
						c.Locals("userId", newClaims["userId"])
						c.Locals("username", newClaims["sub"])
						return c.Next()
					}
				}
			}
		}

		// Not authenticated. Redirect to login.
		// If request is HTMX, return a custom header to trigger a full redirect in the browser.
		if c.Get("HX-Request") == "true" {
			c.Set("HX-Redirect", "/login")
			return c.Status(fiber.StatusUnauthorized).SendString("Session expired. Redirecting...")
		}

		return c.Redirect().To("/login")
	}
}
