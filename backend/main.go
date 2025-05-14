package main

import (
	"blog-db/database"
	"blog-db/middleware"
	"blog-db/models"
	"blog-db/routes"
	"blog-db/seed"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	database.Init()
	database.DB.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Comment{},
		&models.Category{},
		&models.Tags{},
	)

	routes.RegisterAuthRoutes(app) // สำหรับการสมัครสมาชิก และเข้าสู่ระบบ
	routes.SearchedCategoryRoutes(app) // สำหรับการค้นหาหมวดหมู่
	routes.RegisterArticleRoutes(app) // สำหรับการสร้าง แก้ไข และลบบทความ
	routes.GetTagsAll(app) // สำหรับการดึงข้อมูลแท็ก
	routes.RegisterCommentRoutes(app) // สำหรับการสร้าง และดึงข้อมูลคอมเมนต์

	
	protected := app.Group("/", middleware.Protected())
	routes.RegisterUserRoutes(protected) 
	
	// จำลองข้อมูล สำหรับทดสอบ
	seed.SeedCategories()
	seed.SeedTags()  
	seed.SeedUserAndArticles()

	app.Listen(":8080")

}
