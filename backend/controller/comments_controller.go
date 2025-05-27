package controller

import (
	service "backend/service"
	"github.com/gofiber/fiber/v2"
)


// Get Comments
func GetComments(c *fiber.Ctx) error {
	return service.HandleGetComments(c)
}

// Create Comment
func CreateComment(c *fiber.Ctx) error {
	return service.HandleCreateComment(c)
}

