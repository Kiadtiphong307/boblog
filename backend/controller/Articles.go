package controller

import (
	"blog-db/database"
	"blog-db/models"
	"blog-db/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Get all articles
func GetArticles(c *fiber.Ctx) error {
	var articles []models.Article
	if result := database.DB.Find(&articles); result.Error != nil {
		log.Println("❌ Error getting articles:", result.Error)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get articles"))
	}
	log.Println("✅ Retrieved all articles")
	return c.JSON(utils.SuccessResponse(articles, "Articles retrieved successfully"))
}

// Get article by ID
func GetArticle(c *fiber.Ctx) error {
	id := c.Params("id")
	var article models.Article
	result := database.DB.First(&article, id)
	if result.Error != nil {
		log.Println("❌ Article not found with ID:", id)
		return c.Status(404).JSON(utils.ErrorResponse("Article not found"))
	}
	log.Println("✅ Retrieved article with ID:", id)
	return c.JSON(utils.SuccessResponse(article, "Article retrieved successfully"))
}
