package controller 

import (
	service "backend/service"
	"github.com/gofiber/fiber/v2"
)

// คือฟังก์ชันที่จะดึงข้อมูลหมวดหมู่ทั้งหมดจากฐานข้อมูล
func GetCategories(c *fiber.Ctx) error {
	return service.HandleGetCategories(c)
}
