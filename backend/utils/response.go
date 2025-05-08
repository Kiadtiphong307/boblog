package utils

import "github.com/gofiber/fiber/v2"

func SuccessResponse(data interface{}, message string) fiber.Map {
    return fiber.Map{
        "success": true,
        "data":    data,
        "message": message,
    }
}

func ErrorResponse(message string) fiber.Map {
    return fiber.Map{
        "success": false,
        "message": message,
    }
}
