package controller

import (
	service "backend/service"
	"github.com/gofiber/fiber/v2"
)

// คือฟังก์ชันที่จะจัดการกับการลงทะเบียนของผู้ใช้ 
func Register(c *fiber.Ctx) error {
	return service.HandleRegister(c)
}
// คือฟังก์ชันที่จะจัดการกับการลงชื่อเข้าใช้งานของผู้ใช้
func Login(c *fiber.Ctx) error {
	return service.HandleLogin(c)
}

// คือฟังก์ชันที่จะจัดการกับการดึงข้อมูลผู้ใช้
func Profile(c *fiber.Ctx) error {
	return service.HandleProfile(c)
}
// คือฟังก์ชันที่จะจัดการกับการอัปเดตข้อมูลผู้ใช้
func UpdateProfile(c *fiber.Ctx) error {
	return service.HandleUpdateProfile(c)
}
