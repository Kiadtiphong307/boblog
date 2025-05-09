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

// Register
func Register(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid input"))
	}

	// เช็ค email ซ้ำ
	var existing models.User
	database.DB.Where("email = ?", input.Email).First(&existing)
	if existing.ID != 0 {
		return c.Status(400).JSON(utils.ErrorResponse("Email already in use"))
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), 14)
	input.PasswordHash = string(hashedPassword)

	database.DB.Create(&input)
	input.PasswordHash = ""

	return c.Status(201).JSON(utils.SuccessResponse(input, "User registered"))
}

// Login
func Login(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid input"))
	}

	var user models.User
	database.DB.Where("email = ?", input.Email).First(&user)
	if user.ID == 0 {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid email or password"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.PasswordHash)); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid email or password"))
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

// GetUsers retrieves all users from the database
func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	if result := database.DB.Find(&users); result.Error != nil {
		log.Println("❌ Error getting users:", result.Error)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get users"))
	}
	log.Println("✅ Retrieved all users")
	return c.JSON(utils.SuccessResponse(users, "Users retrieved successfully"))
}

// GetUser retrieves a user by ID
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		log.Println("❌ User not found with ID:", id)
		return c.Status(404).JSON(utils.ErrorResponse("User not found"))
	}
	log.Println("✅ Retrieved user with ID:", id)
	return c.JSON(utils.SuccessResponse(user, "User retrieved successfully"))
}

// UpdateUser updates an existing user by ID
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Println("❌ User not found with ID:", id)
		return c.Status(404).JSON(utils.ErrorResponse("User not found"))
	}

	// Parse request body into a map to check which fields are provided
	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		log.Println("❌ Failed to parse update data:", err)
		return c.Status(400).JSON(utils.ErrorResponse("Cannot parse JSON"))
	}

	// Create a temporary user for validation
	tempUser := models.User{
		Username: user.Username,
		Email:    user.Email,
	}

	// Update only provided fields
	if username, ok := updateData["username"].(string); ok && username != "" {
		tempUser.Username = username
	}
	if email, ok := updateData["email"].(string); ok && email != "" {
		tempUser.Email = email
	}

	// Validate the updated data
	if err := validate.Struct(&tempUser); err != nil {
		log.Println("❌ Validation failed:", err)
		return c.Status(400).JSON(utils.ErrorResponse("Validation failed: " + err.Error()))
	}

	// Update the user with validated data
	if username, ok := updateData["username"].(string); ok && username != "" {
		user.Username = username
	}
	if email, ok := updateData["email"].(string); ok && email != "" {
		user.Email = email
	}

	// Save the updated user
	if err := database.DB.Save(&user).Error; err != nil {
		log.Println("❌ Failed to update user:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to update user"))
	}

	log.Println("✅ Updated user with ID:", id)
	return c.JSON(utils.SuccessResponse(user, "User updated successfully"))
}

// DeleteUser deletes a user by ID
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Println("❌ User not found for delete with ID:", id)
		return c.Status(404).JSON(utils.ErrorResponse("User not found"))
	}

	database.DB.Delete(&user)
	log.Println("✅ Deleted user with ID:", id)
	return c.JSON(utils.SuccessResponse(nil, "User deleted successfully"))
}
