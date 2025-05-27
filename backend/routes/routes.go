package routes

import (
	"backend/controller"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

// Register Auth Routes
func RegisterAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", controller.Register)
	auth.Post("/login", controller.Login)
}

// User Routes
func RegisterUserRoutes(r fiber.Router) {
	users := r.Group("/user", middleware.Protected())

	users.Get("/", controller.Profile) // ดูข้อมูลผู้ใช้งานปัจจุบัน
	users.Put("/", controller.UpdateProfile) // แก้ไขข้อมูลผู้ใช้งานปัจจุบัน

}

// Article Routes
func RegisterArticleRoutes(app *fiber.App) {
	articles := app.Group("/articles")

	articles.Get("/my-articles", middleware.Protected(), controller.GetMyArticles) // ดูบทความที่ตัวเองเขียน
	articles.Get("/", controller.SearchArticlesTags)                                   // ค้นหาบทความและแท็ก
	articles.Get("/:slug", controller.GetArticleBySlug)                            // ดูบทความตาม slug

	articles.Post("/", middleware.Protected(), controller.CreateArticle) // สร้างบทความ
	articles.Put("/:slug", middleware.Protected(), controller.UpdateArticle) // แก้ไขบทความ
	articles.Delete("/:slug", middleware.Protected(), controller.DeleteArticle) // ลบบทความ

}

// Search Category Routes
func SearchedCategoryRoutes(app *fiber.App) {
	category := app.Group("/categories")
	category.Get("/", controller.GetCategories)
}

// Get Tags Routes
func GetTagsAll(app *fiber.App) {
	api := app.Group("/tags")
	api.Get("/", controller.GetTags)
}

// Comment Routes
func RegisterCommentRoutes(app *fiber.App) {
    articles := app.Group("/articles")
    articles.Post("/:slug/comments", middleware.Protected(), controller.CreateComment)
    articles.Get("/:slug/comments", controller.GetComments)
}


// 