package main

import (
	"github.com/gofiber/fiber/v2"
	"blog-db/database"
	"blog-db/model"
)

func main() {
	app := fiber.New()

	database.Init()
	database.DB.AutoMigrate(&models.User{})

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Welcome to Go + Fiber + MySQL")
    })

	app.Listen(":8080")

}