package main

import (
	"blog-db/database"
	"blog-db/middleware"
	"blog-db/models"
	"blog-db/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	database.Init()
	database.DB.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Comment{},
		&models.Category{},
		&models.Tag{},
	)

	routes.RegisterAuthRoutes(app)
	routes.SearchedCategoryRoutes(app)    
	routes.RegisterArticleRoutes(app)     
	
	// ✅ ใช้ middleware เฉพาะ route ที่ต้องการเท่านั้น
	protected := app.Group("/", middleware.Protected())
	routes.RegisterUserRoutes(protected)  
	


	app.Listen(":8080")

}
