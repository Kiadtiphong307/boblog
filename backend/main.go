package main

import (
	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/routes"
	"backend/seed"
	

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

	routes.RegisterAuthRoutes(app)
	routes.SearchedCategoryRoutes(app)
	routes.RegisterArticleRoutes(app) 
	routes.GetTagsAll(app) 
	routes.RegisterCommentRoutes(app) 

	
	protected := app.Group("/", middleware.Protected())
	routes.RegisterUserRoutes(protected) 
	
	seed.SeedCategories()
	seed.SeedTags()  
	seed.SeedUserAndArticles()


	app.Listen(":8080")

}
