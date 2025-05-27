package controller

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"backend/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Register
func HandleRegister(c *fiber.Ctx) error {
	var input validation.RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "❌ ไม่สามารถอ่านข้อมูลจากแบบฟอร์มได้",
		})
	}

	if errs := validation.ValidateStructRegister(input); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}

	if input.Password != input.ConfirmPassword {
		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]string{
				"confirm_password": "รหัสผ่านไม่ตรงกัน",
			},
		})
	}

	// Check email exists
	var existing models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]string{
				"email": "อีเมลนี้มีผู้ใช้งานแล้ว",
			},
		})
	}

	// Check username exists
	var existingUsername models.User
	if err := database.DB.Where("username = ?", input.Username).First(&existingUsername).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]string{
				"username": "ชื่อผู้ใช้งานนี้มีผู้ใช้งานแล้ว",
			},
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการเข้ารหัสรหัสผ่าน",
		})
	}

	// Create user
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

	user.PasswordHash = ""
	return c.Status(201).JSON(utils.SuccessResponse(user, "✅ สมัครสมาชิกสำเร็จ"))
}

// Login
func HandleLogin(c *fiber.Ctx) error {
	var input validation.LoginInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid input"))
	}

	if errs := validation.ValidateStructLogin(input); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println("JWT error:", err)
		return c.SendStatus(500)
	}

	return c.JSON(utils.SuccessResponse(fiber.Map{
		"token": tokenString,
	}, "Login successful"))
}

// get current user
func HandleProfile(c *fiber.Ctx) error {
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

// UpdateCurrentUser อัปเดตข้อมูลผู้ใช้ที่เข้าสู่ระบบ
func HandleUpdateProfile(c *fiber.Ctx) error {
	userIDAny := c.Locals("userID")
	userID, ok := userIDAny.(uint)
	if !ok {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("User not found"))
	}

	formHeader := c.Get("Content-Type")
	isMultipart := strings.Contains(strings.ToLower(formHeader), "multipart/form-data")

	// ✅ เคส: multipart/form-data (รองรับอัปโหลดไฟล์)
	if isMultipart {
		firstName := c.FormValue("first_name")
		lastName := c.FormValue("last_name")
		nickname := c.FormValue("nickname")
		bio := c.FormValue("bio")

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

		// ✅ ตรวจสอบและจัดการไฟล์รูปภาพ
		file, err := c.FormFile("avatar")
		if err == nil && file != nil {
			if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
				return c.Status(400).JSON(utils.ErrorResponse("Invalid image format"))
			}

			ext := filepath.Ext(file.Filename)
			filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().Unix(), ext)
			savePath := fmt.Sprintf("./uploads/avatars/%s", filename)

			if err := os.MkdirAll("./uploads/avatars", os.ModePerm); err != nil {
				log.Println("❌ MkdirAll error:", err)
				return c.Status(500).JSON(utils.ErrorResponse("Failed to prepare upload directory"))
			}

			if err := c.SaveFile(file, savePath); err != nil {
				log.Println("❌ SaveFile error:", err)
				return c.Status(500).JSON(utils.ErrorResponse("Failed to save image"))
			}

			imagePath := "/uploads/avatars/" + filename
			user.Image = &imagePath
		}
	} else {
		// ✅ เคส: application/json
		var input struct {
			FirstName string  `json:"first_name"`
			LastName  string  `json:"last_name"`
			Nickname  string  `json:"nickname"`
			Bio       *string `json:"bio"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(utils.ErrorResponse("Invalid JSON format"))
		}

		if input.FirstName != "" {
			user.FirstName = input.FirstName
		}
		if input.LastName != "" {
			user.LastName = input.LastName
		}
		if input.Nickname != "" {
			user.Nickname = input.Nickname
		}
		if input.Bio != nil {
			user.Bio = input.Bio
		}
	}

	if err := database.DB.Save(&user).Error; err != nil {
		log.Println("❌ DB save error:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to update user"))
	}

	// ❌ ห้ามคืน password hash
	user.PasswordHash = ""
	return c.JSON(utils.SuccessResponse(user, "Profile updated successfully"))
}
