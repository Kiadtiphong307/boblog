package routes

import (
	"blog-db/controller"
    "github.com/gofiber/fiber/v2"
    "blog-db/middleware"

)

func RegisterAuthRoutes(app *fiber.App) {
    auth := app.Group("/auth")
    auth.Post("/register",  controller.Register)
    auth.Post("/login", controller.Login)
}

func RegisterUserRoutes(app *fiber.App) {
    users := app.Group("/users").Use(middleware.Protected())
    users.Get("/", controller.GetUsers)
    users.Get("/:id", controller.GetUser)
    users.Put("/:id", controller.UpdateUser)
    users.Delete("/:id", controller.DeleteUser)

    app.Get("/user", middleware.Protected(), controller.GetCurrentUser)
}

