package routes

import (
	"blog-db/controller"
	"blog-db/middleware"

	"github.com/gofiber/fiber/v2"
)

// สมัครสมาชิก และเข้าสู่ระบบ
func RegisterAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", controller.Register)
	auth.Post("/login", controller.Login)
}

// ข้อมูลผู้ใช้งานปัจจุบัน
func RegisterUserRoutes(r fiber.Router) {
	users := r.Group("/user", middleware.Protected())

	users.Get("/", controller.GetCurrentUser) // ดูข้อมูลผู้ใช้งานปัจจุบัน
	users.Put("/", controller.UpdateCurrentUser) // แก้ไขข้อมูลผู้ใช้งานปัจจุบัน

}

// บทความ
func RegisterArticleRoutes(app *fiber.App) {
	articles := app.Group("/articles")

	articles.Get("/my-articles", middleware.Protected(), controller.GetMyArticles) // ดูบทความที่ตัวเองเขียน
	articles.Get("/", controller.SearchArticles)                                   // ดูทั้งหมด
	articles.Get("/:slug", controller.GetArticleBySlug)                            // ดูบทความตาม slug

	articles.Post("/", middleware.Protected(), controller.CreateArticle) // สร้างบทความ
	articles.Put("/:slug", middleware.Protected(), controller.UpdateArticle) // แก้ไขบทความ
	articles.Delete("/:slug", middleware.Protected(), controller.DeleteArticle) // ลบบทความ

}

// ค้นหาหมวดหมู่
func SearchedCategoryRoutes(app *fiber.App) {
	category := app.Group("/categories")
	category.Get("/", controller.GetCategories)
}

// ดึงข้อมูลแท็ก
func GetTagsAll(app *fiber.App) {
	api := app.Group("/tags")
	api.Get("/", controller.GetTags)
}

