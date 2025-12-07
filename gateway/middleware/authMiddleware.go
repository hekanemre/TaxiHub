package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hekanemre/taxihub/gateway/helpers"
)

// Authenticate returns a Fiber middleware that checks JWT tokens
func Authenticate(tokenHelper *helpers.TokenHelper) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientToken := c.Get("token")
		if clientToken == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "No Authorization token provided",
			})
		}

		claims, errStr := tokenHelper.ValidateToken(clientToken)
		if errStr != "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": errStr,
			})
		}

		// store user info in request context
		c.Locals("email", claims.Email)
		c.Locals("first_name", claims.First_name)
		c.Locals("last_name", claims.Last_name)
		c.Locals("uid", claims.Uid)
		c.Locals("user_type", claims.User_type)

		return c.Next()
	}
}
