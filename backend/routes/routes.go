package routes

import (
	"blog-db/controller"
    "github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
    // Public routes
    app.Post("/register", controller.Register)
    app.Post("/login", controller.Login)

    // Protected routes
    users := app.Group("/users", controller.JWTMiddleware())
    users.Get("/", controller.GetUsers)
    users.Get("/:id", controller.GetUser)
    users.Post("/", controller.CreateUser)
    users.Put("/:id", controller.UpdateUser)
    users.Delete("/:id", controller.DeleteUser)
}
