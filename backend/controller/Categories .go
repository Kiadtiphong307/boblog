package controller

import (
	"blog-db/database"
	"blog-db/models"
	"blog-db/utils"
	"github.com/gofiber/fiber/v2"
	"log"
)

// Get All Categories
func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category

	if result := database.DB.Find(&categories); result.Error != nil {
		log.Println("❌ Error loading categories:", result.Error)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to load categories"))
	}

	return c.JSON(utils.SuccessResponse(categories, "Categories retrieved"))
}

