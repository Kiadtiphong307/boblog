package controller

import (
	"blog-db/database"
	"blog-db/models"
	"blog-db/utils"
	"log"
	"time"
	"strings"
	"fmt"
	"path/filepath"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

var jwtSecret = []byte("secret") 

// Register

func Register(c *fiber.Ctx) error {
	var input struct {
		Username        string `json:"username" validate:"required,min=3"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
		FirstName       string `json:"first_name" validate:"required"`
		LastName        string `json:"last_name" validate:"required"`
		Nickname        string `json:"nickname" validate:"required"`
	}

	// check body is valid
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "❌ ไม่สามารถอ่านข้อมูลจากแบบฟอร์มได้",
		})
	}

	// check body is valid
	if err := validate.Struct(input); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = "จำเป็นต้องกรอก"
			case "email":
				errors[field] = "อีเมลไม่ถูกต้อง"
			case "min":
				errors[field] = fmt.Sprintf("ต้องมีความยาวอย่างน้อย %s ตัวอักษร", e.Param())
			default:
				errors[field] = "ข้อมูลไม่ถูกต้อง"
			}
		}
		return c.Status(400).JSON(fiber.Map{
			"errors": errors,
		})
	}

	// check password is valid
	if input.Password != input.ConfirmPassword {
		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]string{
				"confirm_password": "รหัสผ่านไม่ตรงกัน",
			},
		})
	}

	// check email is exist
	var existing models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]string{
				"email": "อีเมลนี้มีผู้ใช้งานแล้ว",
			},
		})
	}

	// check username is exist
	var existingUsername models.User
	if err := database.DB.Where("username = ?", input.Username).First(&existingUsername).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]string{
				"username": "ชื่อผู้ใช้งานนี้มีผู้ใช้งานแล้ว",
			},
		})
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการเข้ารหัสรหัสผ่าน",
		})
	}

	// create user
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Nickname:     input.Nickname,
		PasswordHash: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "ไม่สามารถสร้างบัญชีผู้ใช้ได้",
		})
	}

	// don't return password
	user.PasswordHash = ""

	return c.Status(201).JSON(utils.SuccessResponse(user, "✅ สมัครสมาชิกสำเร็จ"))
}

// Login
func Login(c *fiber.Ctx) error {
	var input struct {
		EmailOrUsername string `json:"email"`
		Password        string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid input"))
	}

	var user models.User
	database.DB.
		Where("email = ? OR username = ?", input.EmailOrUsername, input.EmailOrUsername).
		First(&user)

	if user.ID == 0 {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid credentials"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid credentials"))
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("JWT error:", err)
		return c.SendStatus(500)
	}

	return c.JSON(utils.SuccessResponse(fiber.Map{
		"token": tokenString,
	}, "Login successful"))
}

// get current user
func GetCurrentUser(c *fiber.Ctx) error {
	// check if token is sent
	userToken := c.Locals("user")
	if userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Unauthorized"))
	}

	// Convert to jwt.Token safely
	token, ok := userToken.(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid token"))
	}

	// Convert claims and get ID
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	// Get user from database
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Println("❌ User not found with ID from token:", userID)
		return c.Status(404).JSON(utils.ErrorResponse("User not found"))
	}

	// Remove password before sending response
	user.PasswordHash = ""
	return c.JSON(utils.SuccessResponse(user, "Current user"))
}

// Update Current User
func UpdateCurrentUser(c *fiber.Ctx) error {
	userIDAny := c.Locals("userID")
	userID, ok := userIDAny.(uint)
	if !ok {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}

	formHeader := c.Get("Content-Type")
	isMultipart := strings.Contains(strings.ToLower(formHeader), "multipart/form-data")

	var (
		firstName = c.FormValue("first_name")
		lastName  = c.FormValue("last_name")
		nickname  = c.FormValue("nickname")
		bio       = c.FormValue("bio")
		imageURL  *string
	)

	if isMultipart {
		file, err := c.FormFile("avatar")
		if err == nil && file != nil {
			if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
				return c.Status(400).JSON(utils.ErrorResponse("Invalid image format"))
			}

			ext := filepath.Ext(file.Filename)
			filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().Unix(), ext)
			savePath := fmt.Sprintf("./uploads/avatars/%s", filename)

			// ✅ สร้างโฟลเดอร์หากยังไม่มี
			if err := os.MkdirAll("./uploads/avatars", os.ModePerm); err != nil {
				log.Println("❌ MkdirAll error:", err)
				return c.Status(500).JSON(utils.ErrorResponse("Failed to prepare upload directory"))
			}

			// ✅ บันทึกไฟล์
			if err := c.SaveFile(file, savePath); err != nil {
				log.Println("❌ SaveFile error:", err)
				return c.Status(500).JSON(utils.ErrorResponse("Failed to save image"))
			}

			imagePath := "/uploads/avatars/" + filename
			imageURL = &imagePath
		}
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("User not found"))
	}

	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}
	if nickname != "" {
		user.Nickname = nickname
	}
	if bio != "" {
		user.Bio = &bio
	}
	if imageURL != nil {
		user.Image = imageURL
	}

	if err := database.DB.Save(&user).Error; err != nil {
		log.Println("❌ DB save error:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to update user"))
	}

	user.PasswordHash = ""
	return c.JSON(utils.SuccessResponse(user, "Profile updated successfully"))
}


	




