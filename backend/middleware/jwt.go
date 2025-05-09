package middleware

import (
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret")

func Protected() fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenStr := c.Get("Authorization")
        if tokenStr == "" {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
        }

        // รองรับ prefix "Bearer "
        parts := strings.Split(tokenStr, " ")
        if len(parts) == 2 && parts[0] == "Bearer" {
            tokenStr = parts[1]
        }

        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
        }

        // ✅ บันทึก token ลง context เพื่อให้ controller ใช้ได้
        c.Locals("user", token)

        return c.Next()
    }
}
