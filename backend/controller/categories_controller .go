package controller 

import (
	service "backend/service"
	"github.com/gofiber/fiber/v2"
)

// Get All Categories
func GetCategories(c *fiber.Ctx) error {
	return service.HandleGetCategories(c)
}
