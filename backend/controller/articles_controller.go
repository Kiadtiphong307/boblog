package controller

import (
	service "backend/service"
	"github.com/gofiber/fiber/v2"
)



// คือฟังก์ชันที่จะดึงข้อมูลบทความทั้งหมดจากฐานข้อมูล
func GetAllArticles(c *fiber.Ctx) error {
	return service.HandleGetAllArticles(c)
}
// คือฟังก์ชันที่จะค้นหาบทความตามหมวดหมู่และแท็ก
func SearchArticlesTags(c *fiber.Ctx) error {
	return service.HandleSearchArticlesTags(c)
}

// คือฟังก์ชันที่จะดึงข้อมูลบทความตาม slug
func GetArticleBySlug(c *fiber.Ctx) error {
	return service.HandleGetArticleBySlug(c)
}

// CRUD Write Article
// คือฟังก์ชันที่จะดึงข้อมูลบทความของผู้ใช้งานปัจจุบัน 
func GetMyArticles(c *fiber.Ctx) error {
	return service.HandleGetMyArticles(c)
}

// คือฟังก์ชันที่จะสร้างบทความใหม่ จาก service/articles_service.go
func CreateArticle(c *fiber.Ctx) error {
	return service.HandleCreateArticle(c)
}

// คือฟังก์ชันที่จะแก้ไขบทความ จาก service/articles_service.go
func UpdateArticle(c *fiber.Ctx) error {
	return service.HandleUpdateArticle(c)
}

// คือฟังก์ชันที่จะลบบทความ
func DeleteArticle(c *fiber.Ctx) error {
	return service.HandleDeleteArticle(c)
}
