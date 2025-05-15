package controller

import (
	"blog-db/database"
	"blog-db/models"
	"blog-db/utils"
	"github.com/gofiber/fiber/v2"
	"log"
)

// Get All Tags
func GetTags(c *fiber.Ctx) error {
	var tags []models.Tags

	if result := database.DB.Find(&tags); result.Error != nil {
		log.Println("❌ Error loading tags:", result.Error)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to load tags"))
	}

	return c.JSON(utils.SuccessResponse(tags, "Tags retrieved"))
}

