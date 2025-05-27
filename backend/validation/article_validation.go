package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Struct สำหรับสร้างบทความ
type CreateArticleInput struct {
	Title        string   `json:"title" validate:"required"`
	Slug         string   `json:"slug" validate:"required"`
	Content      string   `json:"content" validate:"required"`
	CategoryName string   `json:"category_name" validate:"required"`
	TagNames     []string `json:"tag_names"`
}

// Struct สำหรับแก้ไขบทความ
type UpdateArticleInput struct {
	Title      *string  `json:"title"`
	Content    *string  `json:"content"`
	CategoryID *uint    `json:"category_id"`
	TagIDs     []uint   `json:"tag_ids"`
	NewTags    []string `json:"new_tags"`
}

// ฟังก์ชันตรวจสอบ struct ทั่วไป
func ValidateStructArticle(data interface{}) map[string]string {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			errors[field] = "จำเป็นต้องกรอก"
		default:
			errors[field] = fmt.Sprintf("ข้อมูลไม่ถูกต้อง")
		}
	}
	return errors
}
