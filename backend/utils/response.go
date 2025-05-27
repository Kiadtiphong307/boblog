package utils

import "github.com/gofiber/fiber/v2"

// คือฟังก์ชันที่จะส่งข้อมูลสำเร็จกลับให้กับผู้ใช้
func SuccessResponse(data interface{}, message string) fiber.Map {
	return fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	}
}

// คือฟังก์ชันที่จะส่งข้อมูลผิดพลาดกลับให้กับผู้ใช้
func ErrorResponse(message string) fiber.Map {
	return fiber.Map{
		"success": false,
		"message": message,
	}
}
