package main

import (
	"github.com/gofiber/fiber/v2"
	"blog-db/database"
	"blog-db/models"
	"blog-db/routes"
)

func main() {
	app := fiber.New()

	database.Init()
	database.DB.AutoMigrate(&models.User{})

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Welcome to Go + Fiber + MySQL")
    })
	routes.RegisterUserRoutes(app)

	app.Listen(":8080")

}