package routes

import (
	"backend/controller"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)




// Register and Login Routes
func AuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", controller.Register)
	auth.Post("/login", controller.Login)
	
}

// Profile Routes
func UserRoutes(r fiber.Router) {
	users := r.Group("/user", middleware.Protected())

	users.Get("/", controller.Profile) // ดูข้อมูลผู้ใช้งานปัจจุบัน
	users.Put("/", controller.UpdateProfile) // แก้ไขข้อมูลผู้ใช้งานปัจจุบัน
}