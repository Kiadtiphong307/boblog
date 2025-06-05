package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var userValidator = validator.New()

// Register Input Struct
type RegisterInput struct {
	Username        string `json:"username" validate:"required,min=3"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Nickname        string `json:"nickname" validate:"required"`
}

// Login Input Struct
type LoginInput struct {
	EmailOrUsername string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

// validate register input
func ValidateStructRegister(data interface{}) map[string]string {
	err := userValidator.Struct(data)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			errors[field] = "This field is required"
		case "email":
			errors[field] = "Invalid email"
		case "min":
			errors[field] = fmt.Sprintf("Minimum %s characters", e.Param())
		default:
			errors[field] = "Invalid data"
		}
	}
	return errors
}

// validate login input
func ValidateStructLogin(data interface{}) map[string]string {
	err := userValidator.Struct(data)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errors[field] = "Invalid data"
	}
	return errors
}
