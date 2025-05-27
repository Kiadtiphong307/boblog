package middleware

import (
	"os"
	"fmt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ContextKey: "user", // เก็บ token ไว้ใน c.Locals("user")
		ErrorHandler: jwtError,
		SuccessHandler: func(c *fiber.Ctx) error {
			userToken := c.Locals("user")
			token, ok := userToken.(*jwt.Token)
			if !ok {
				return jwtError(c, fmt.Errorf("invalid token"))
			}

			claims := token.Claims.(jwt.MapClaims)
			if idFloat, ok := claims["id"].(float64); ok {
				c.Locals("userID", uint(idFloat)) // ✅ Set userID ที่เราจะใช้งานใน Controller
			}

			return c.Next()
		},
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
}