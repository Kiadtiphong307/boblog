package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret") // เปลี่ยนเป็น env ในโปรดักชัน

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึง Authorization header
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: missing token",
			})
		}

		// ตรวจรูปแบบ Bearer <token>
		parts := strings.SplitN(tokenStr, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}
		tokenStr = parts[1]

		// แปลงและตรวจสอบ JWT
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// แปลง claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["id"] == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// แปลง userID จาก float64 → uint อย่างปลอดภัย
		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token",
			})
		}
		userID := uint(userIDFloat)

		// เก็บ userID ลง context
		c.Locals("userID", userID)
		c.Locals("user", token)

		return c.Next()
	}
}
