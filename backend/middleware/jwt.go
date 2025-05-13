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

        // แยก Bearer token ออก
        parts := strings.Split(tokenStr, " ")
        if len(parts) == 2 && parts[0] == "Bearer" {
            tokenStr = parts[1]
        }

        // แปลง token
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        if err != nil || !token.Valid {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
        }

        // แปลง claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || claims["id"] == nil {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
        }

        // ดึง userID จาก claims แล้วเก็บใน context
        userIDFloat, ok := claims["id"].(float64)
        if !ok {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid user ID in token"})
        }

        userID := uint(userIDFloat)
        c.Locals("userID", userID)
        c.Locals("user", token)

        return c.Next()
    }
}
