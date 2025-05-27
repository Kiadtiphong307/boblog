package utils

import "github.com/gofiber/fiber/v2"

func SuccessResponse(data interface{}, message string) fiber.Map {
	return fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	}
}

func ErrorResponse(message string) fiber.Map {
	return fiber.Map{
		"success": false,
		"message": message,
	}
}
