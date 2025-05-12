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

func RegisterUserRoutes(r fiber.Router) {
	users := r.Group("/users")
	users.Put("/:id", controller.UpdateUser)
	users.Delete("/:id", controller.DeleteUser)

	r.Get("/user", controller.GetCurrentUser)
}

func RegisterArticleRoutes(app *fiber.App) {
	articles := app.Group("/articles")
	articles.Post("/", middleware.Protected(), controller.CreateArticle)
	articles.Get("/", controller.SearchArticles)         // ใช้ใน frontend
	articles.Get("/all", middleware.Protected(), controller.GetAllArticles) // สำหรับ admin
	articles.Get("/:slug", controller.GetArticleBySlug)
}


func SearchedCategoryRoutes(app *fiber.App) {
	category := app.Group("/categories")
	category.Get("/", controller.GetCategories) 
}
