package service  

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"github.com/gofiber/fiber/v2"
)

// Get All Tags
func HandleGetTags(c *fiber.Ctx) error {
	var tags []models.Tags

	if result := database.DB.Find(&tags); result.Error != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to load tags"))
	}

	return c.JSON(utils.SuccessResponse(tags, "Tags retrieved"))
}

