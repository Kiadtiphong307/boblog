package controller

import (
	service "backend/service"
	"github.com/gofiber/fiber/v2"
)

// คือฟังก์ชันที่จะดึงข้อมูลแท็กทั้งหมดจากฐานข้อมูล
func GetTags(c *fiber.Ctx) error {
	return service.HandleGetTags(c)
}
