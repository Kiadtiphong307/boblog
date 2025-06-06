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

// check if email exists
func IsEmailExists(email string) (bool, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return true, nil
	}
	return false, nil
}

// check if username exists
func IsUsernameExists(username string) (bool, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err == nil {
		return true, nil
	}
	return false, nil
}

// get user by id
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// get user id from jwt
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

// upload avatar
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
