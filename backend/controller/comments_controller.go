package controller

import (
	service "backend/service"
	"github.com/gofiber/fiber/v2"
)


// คือฟังก์ชันที่จะดึงข้อมูลคอมเมนต์ทั้งหมดจากฐานข้อมูล
func GetComments(c *fiber.Ctx) error {
	return service.HandleGetComments(c)
}

// คือฟังก์ชันที่จะสร้างคอมเมนต์ใหม่
func CreateComment(c *fiber.Ctx) error {
	return service.HandleCreateComment(c)
}

// คือฟังก์ชันที่จะแก้ไขคอมเมนต์
func UpdateComment(c *fiber.Ctx) error {
	return service.HandleUpdateComment(c)
}

// คือฟังก์ชันที่จะลบคอมเมนต์
func DeleteComment(c *fiber.Ctx) error {
	return service.HandleDeleteComment(c)
}



