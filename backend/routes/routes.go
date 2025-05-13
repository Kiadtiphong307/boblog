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
	// ✅ Group: /users ต้องใช้ token ผ่าน middleware
	users := r.Group("/users")
	users.Put("/:id", middleware.Protected(), controller.UpdateUser)
	users.Delete("/:id", middleware.Protected(), controller.DeleteUser)

	// ✅ ใช้ดึงข้อมูล user จาก JWT (เช่นหน้าโปรไฟล์)
	r.Get("/user", middleware.Protected(), controller.GetCurrentUser)
}

func RegisterArticleRoutes(app *fiber.App) {
	articles := app.Group("/articles")

	articles.Get("/my-articles", middleware.Protected(), controller.GetMyArticles) // ดูบทความที่ตัวเองเขียน
	articles.Get("/", controller.SearchArticles)                                   // ดูทั้งหมด
	articles.Get("/:slug", controller.GetArticleBySlug)                            // ดูบทความตาม slug

	articles.Post("/", middleware.Protected(), controller.CreateArticle) // สร้างบทความ
	articles.Put("/:slug", middleware.Protected(), controller.UpdateArticle) // แก้ไขบทความ
	articles.Delete("/:slug", middleware.Protected(), controller.DeleteArticle) // ลบบทความ

}

func SearchedCategoryRoutes(app *fiber.App) {
	category := app.Group("/categories")
	category.Get("/", controller.GetCategories)
}

func GetTagsAll(app *fiber.App) {
	api := app.Group("/tags")
	api.Get("/", controller.GetTags)
}

