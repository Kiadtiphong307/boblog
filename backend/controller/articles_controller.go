package controller

import (
	service "backend/service"
	"github.com/gofiber/fiber/v2"
)



// Get all articles
func GetAllArticles(c *fiber.Ctx) error {
	return service.HandleGetAllArticles(c)
}
// filter articles and tags
func SearchArticlesTags(c *fiber.Ctx) error {
	return service.HandleSearchArticlesTags(c)
}

// Get article by slug
func GetArticleBySlug(c *fiber.Ctx) error {
	return service.HandleGetArticleBySlug(c)
}

// Get my articles
func GetMyArticles(c *fiber.Ctx) error {
	return service.HandleGetMyArticles(c)
}

func CreateArticle(c *fiber.Ctx) error {
	return service.HandleCreateArticle(c)
}

// UpdateArticle updates an existing article
func UpdateArticle(c *fiber.Ctx) error {
	return service.HandleUpdateArticle(c)
}

// DeleteArticle deletes an article by slug (only by the author)
func DeleteArticle(c *fiber.Ctx) error {
	return service.HandleDeleteArticle(c)
}
