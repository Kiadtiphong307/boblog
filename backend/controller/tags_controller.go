package controller

import (
	service "backend/service"
	"github.com/gofiber/fiber/v2"
)

// Get All Tags
func GetTags(c *fiber.Ctx) error {
	return service.HandleGetTags(c)
}
