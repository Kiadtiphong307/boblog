package routes

import (
	"blog-db/controller"
	"blog-db/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", controller.Register)
	auth.Post("/login", controller.Login)
}

func RegisterUserRoutes(app *fiber.App) {
	users := app.Group("/users").Use(middleware.Protected())
	users.Put("/:id", controller.UpdateUser)
	users.Delete("/:id", controller.DeleteUser)

	app.Get("/user", middleware.Protected(), controller.GetCurrentUser)
}

func RegisterArticleRoutes(app *fiber.App) {
	articles := app.Group("/articles")
	articles.Get("/", controller.GetArticles)
	articles.Post("/", middleware.Protected(), controller.CreateArticle)
	articles.Get("/:slug", controller.GetArticleBySlug)
}
