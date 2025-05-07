package routes

import (
	"blog-db/controller"
    "github.com/gofiber/fiber/v2"

)

func RegisterUserRoutes(app *fiber.App) {
    users := app.Group("/users")
    users.Get("/", controller.GetUsers)
    users.Get("/:id", controller.GetUser)
    users.Post("/", controller.CreateUser)
    users.Put("/:id", controller.UpdateUser)
    users.Delete("/:id", controller.DeleteUser)
}
