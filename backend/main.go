package main

import (
	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/routes"
	"backend/seed"
	

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// คือการสร้าง folder สำหรับส่งข้อมูลรูปภาพ
	app.Static("/uploads", "/app/uploads")

	database.Init()
	database.DB.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Comment{},
		&models.Category{},
		&models.Tags{},
	)

	routes.AuthRoutes(app)  // ลงทะเบียนและเข้าสู่ระบบ
	routes.CategoryRoutes(app) // ดูหมวดหมู่บทความ
	routes.ArticleRoutes(app)  // ดูบทความ
	routes.GetTagsAll(app)  // ดูแท็ก
	routes.CommentRoutes(app)  // ดูคอมเมนต์

	
	protected := app.Group("/", middleware.Protected())
	routes.UserRoutes(protected)  // ดูข้อมูลผู้ใช้
	
	// Database Seed
	seed.SeedCategories()
	seed.SeedTags()  
	seed.SeedUserAndArticles()


	app.Listen(":8080")

}
