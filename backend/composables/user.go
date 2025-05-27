package composables

import (
	"backend/database"
	"backend/models"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// คือฟังก์ชันที่จะตรวจสอบว่าอีเมลมีอยู่แล้วหรือไม่
func IsEmailExists(email string) (bool, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return true, nil
	}
	return false, nil
}
// คือฟังก์ชันที่จะตรวจสอบว่าชื่อผู้ใช้งานมีอยู่แล้วหรือไม่
func IsUsernameExists(username string) (bool, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err == nil {
		return true, nil
	}
	return false, nil
}

// คือฟังก์ชันที่จะดึงข้อมูลผู้ใช้งานจากฐานข้อมูล
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// คือฟังก์ชันที่จะดึงข้อมูลผู้ใช้งานจาก JWT
func ExtractUserIDFromJWT(c *fiber.Ctx) (uint, error) {
	userToken := c.Locals("user")
	token, ok := userToken.(*jwt.Token)
	if !ok {
		return 0, errors.New("Invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	idFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("Invalid claims")
	}
	return uint(idFloat), nil
}
// คือฟังก์ชันที่จะอัปเดตรูปภาพผู้ใช้
func HandleAvatarUpload(c *fiber.Ctx, userID uint) (*string, error) {
	file, err := c.FormFile("avatar")
	if err != nil || file == nil {
		return nil, nil
	}
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return nil, errors.New("Invalid image format")
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().Unix(), ext)
	savePath := fmt.Sprintf("./uploads/avatars/%s", filename)
	if err := os.MkdirAll("./uploads/avatars", os.ModePerm); err != nil {
		return nil, errors.New("Failed to prepare upload directory")
	}
	if err := c.SaveFile(file, savePath); err != nil {
		return nil, errors.New("Failed to save image")
	}
	imagePath := "/uploads/avatars/" + filename
	return &imagePath, nil
}
