package middlewares

import (
	"github.com/MarcelArt/gotel/internal/common"
	"github.com/MarcelArt/gotel/internal/configs"
	"github.com/MarcelArt/gotel/internal/enums"
	"github.com/MarcelArt/gotel/pkg/arrays"
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
)

func Refresh() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(configs.Env.JwtSecret)},
		Extractor:  extractors.FromHeader("X-Refresh-Token"),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(common.NewJSONResponse(err, "unauthorized"))
		},
	})
}

func Authn() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(configs.Env.JwtSecret)},
		Extractor:  extractors.FromAuthHeader("Bearer"),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(common.NewJSONResponse(err, "unauthorized"))
		},
	})
}

func Authz(permissionKey string) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims := common.FiberCtxToClaims(c)

		permissions, err := common.ParseClaimsToStringSlice(claims["permissions"])
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(common.NewJSONResponse(err, "unauthorized"))
		}

		permission := arrays.Find(permissions, func(p string) bool {
			return p == enums.PermFullAccess || p == permissionKey
		})

		if permission == nil {
			return c.Status(fiber.StatusForbidden).JSON(common.NewJSONResponse(nil, "unauthorized"))
		}

		return c.Next()
	}
}
