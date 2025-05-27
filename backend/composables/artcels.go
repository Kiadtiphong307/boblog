package composables

import (
	"backend/database"
	"backend/models"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// ดึง userID จาก token
func GetCurrentUserID(c *fiber.Ctx) (uint, error) {
	userToken := c.Locals("user").(*jwt.Token)
	if userToken == nil {
		return 0, errors.New("missing user token")
	}
	claims := userToken.Claims.(jwt.MapClaims)
	id, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token")
	}
	return uint(id), nil
}

// ค้นหาหรือสร้าง category
func FindOrCreateCategory(name string) (*models.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("empty category name")
	}
	var category models.Category
	if err := database.DB.Where("name = ?", name).First(&category).Error; err == nil {
		return &category, nil
	}
	category = models.Category{Name: name}
	if err := database.DB.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// ค้นหาหรือสร้าง tag หลายตัว
func FindOrCreateTags(names []string) ([]models.Tags, error) {
	var tags []models.Tags
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var tag models.Tags
		if err := database.DB.Where("name = ?", name).First(&tag).Error; err == nil {
			tags = append(tags, tag)
		} else {
			tag = models.Tags{Name: name}
			if err := database.DB.Create(&tag).Error; err != nil {
				return nil, err
			}
			tags = append(tags, tag)
		}
	}
	return tags, nil
}
