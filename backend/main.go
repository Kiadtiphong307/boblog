package main

import (
	"blog-db/database"
	"blog-db/middleware"
	"blog-db/models"
	"blog-db/routes"
	"blog-db/seed"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// For sending image data
	app.Static("/uploads", "/app/uploads")

	database.Init()
	database.DB.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Comment{},
		&models.Category{},
		&models.Tags{},
	)

	routes.RegisterAuthRoutes(app) // For registering and logging in
	routes.SearchedCategoryRoutes(app) // For searching categories
	routes.RegisterArticleRoutes(app) // For creating, updating, and deleting articles
	routes.GetTagsAll(app) // For getting tag data
	routes.RegisterCommentRoutes(app) // For creating and getting comments

	
	protected := app.Group("/", middleware.Protected())
	routes.RegisterUserRoutes(protected) 
	
	// Simulate data for testing
	seed.SeedCategories()
	seed.SeedTags()  
	seed.SeedUserAndArticles()


	app.Listen(":8080")

}
