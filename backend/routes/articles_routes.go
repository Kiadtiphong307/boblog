package routes

import (
	"backend/controller"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)



// Article Routes
func ArticleRoutes(app *fiber.App) {
	articles := app.Group("/articles")

	articles.Get("/my-articles", middleware.Protected(), controller.GetMyArticles) // ดูบทความที่ตัวเองเขียน
	articles.Get("/", controller.SearchArticlesTags)                                // ค้นหาบทความและแท็ก
	articles.Get("/:slug", controller.GetArticleBySlug)                            	// ดูบทความตาม slug

	articles.Post("/", middleware.Protected(), controller.CreateArticle) 			// สร้างบทความ
	articles.Put("/:slug", middleware.Protected(), controller.UpdateArticle) // แก้ไขบทความ
	articles.Delete("/:slug", middleware.Protected(), controller.DeleteArticle) // ลบบทความ

}

// Search Category Routes
func CategoryRoutes(app *fiber.App) {
	category := app.Group("/categories")
	category.Get("/", controller.GetCategories) // ดูหมวดหมู่บทความ
}





