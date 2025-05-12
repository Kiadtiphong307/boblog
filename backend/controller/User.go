package controller

import (
	"blog-db/database"
	"blog-db/models"
	"blog-db/utils"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

var jwtSecret = []byte("secret") // 👉 เปลี่ยนเป็น os.Getenv("JWT_SECRET") ภายหลัง

// GetCurrentUser
func GetCurrentUser(c *fiber.Ctx) error {
	// ตรวจสอบว่า token ถูกส่งมาหรือไม่
	userToken := c.Locals("user")
	if userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Unauthorized"))
	}

	// แปลงเป็น jwt.Token อย่างปลอดภัย
	token, ok := userToken.(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse("Invalid token"))
	}

	// แปลง claims และดึง ID
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	// ดึง user จาก database
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Println("❌ User not found with ID from token:", userID)
		return c.Status(404).JSON(utils.ErrorResponse("User not found"))
	}

	// ลบ password ออกก่อนตอบกลับ
	user.PasswordHash = ""
	return c.JSON(utils.SuccessResponse(user, "Current user"))
}


func Register(c *fiber.Ctx) error {
	// 1. รับค่าจาก JSON พร้อมฟิลด์ใหม่
	var input struct {
		Username        string `json:"username" validate:"required,min=3"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
		FirstName       string `json:"first_name" validate:"required"`
		LastName        string `json:"last_name" validate:"required"`
		Nickname        string `json:"nickname" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid input"))
	}

	// 2. ตรวจสอบความถูกต้อง
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Validation failed: " + err.Error()))
	}

	// 2.5 ตรวจว่ารหัสผ่านตรงกันหรือไม่
	if input.Password != input.ConfirmPassword {
		return c.Status(400).JSON(utils.ErrorResponse("Passwords do not match"))
	}

	// 3. ตรวจ email ซ้ำ
	var existing models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return c.Status(400).JSON(utils.ErrorResponse("Email already in use"))
	}

	// 4. เข้ารหัส password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to hash password"))
	}

	// 5. สร้าง user และบันทึก
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Nickname:     input.Nickname,
		PasswordHash: string(hashedPassword),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to create user"))
	}

	// 6. ลบ password ก่อนส่งกลับ
	user.PasswordHash = ""

	return c.Status(201).JSON(utils.SuccessResponse(user, "User registered successfully"))
}


// Login
func Login(c *fiber.Ctx) error {
	var input struct {
		EmailOrUsername string `json:"email"` // ใช้ key เดิมเพื่อไม่ต้องแก้ frontend
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

	// ✅ Create JWT
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



// แก้ไขข้อมูลผู้ใช้
func UpdateUser(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    var input struct {
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
        Nickname  string `json:"nickname"`
        Bio       *string `json:"bio"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(utils.ErrorResponse("Invalid input"))
    }

    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        return c.Status(404).JSON(utils.ErrorResponse("User not found"))
    }

    user.FirstName = input.FirstName
    user.LastName = input.LastName
    user.Nickname = input.Nickname
    user.Bio = input.Bio

    database.DB.Save(&user)
    user.PasswordHash = ""
    return c.JSON(utils.SuccessResponse(user, "Profile updated"))
}

// ลบบัญชีผู้ใช้
func DeleteUser(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)

    if err := database.DB.Delete(&models.User{}, userID).Error; err != nil {
        return c.Status(500).JSON(utils.ErrorResponse("Failed to delete account"))
    }

    return c.JSON(utils.SuccessResponse(nil, "Account deleted successfully"))
}
