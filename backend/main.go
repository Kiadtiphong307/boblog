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
	app.Use(middleware.Protected())
	routes.RegisterUserRoutes(app)
	routes.RegisterArticleRoutes(app)

	app.Listen(":8080")

}
