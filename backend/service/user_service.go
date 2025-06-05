package service  

import (
	"backend/composables"
	"backend/database"
	"backend/models"
	"backend/utils"
	"backend/validation"
	

	"github.com/gofiber/fiber/v2"
	"strings"
)

// register user
func HandleRegister(c *fiber.Ctx) error {
	var input validation.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Failed to read form data")) 
	}
	if errs := validation.ValidateStructRegister(input); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	if input.Password != input.ConfirmPassword {
		return c.Status(400).JSON(fiber.Map{"errors": map[string]string{"confirm_password": "Password does not match"}})
	}
	if exists, _ := composables.IsEmailExists(input.Email); exists {
		return c.Status(400).JSON(fiber.Map{"errors": map[string]string{"email": "Email already exists"}})
	}
	if exists, _ := composables.IsUsernameExists(input.Username); exists {
		return c.Status(400).JSON(fiber.Map{"errors": map[string]string{"username": "Username already exists"}})
	}

	hashedPassword := utils.HashPassword(input.Password)
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Nickname:     input.Nickname,
		PasswordHash: hashedPassword,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to create user"))
	}
	user.PasswordHash = ""
	return c.Status(201).JSON(utils.SuccessResponse(user, "Register successful"))
}

// login user
func HandleLogin(c *fiber.Ctx) error {
	var input validation.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Failed to read form data"))
	}
	if errs := validation.ValidateStructLogin(input); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}

	var user models.User
	database.DB.
		Where("email = ? OR username = ?", input.EmailOrUsername, input.EmailOrUsername).
		First(&user)
	if user.ID == 0 || !utils.CheckPassword(input.Password, user.PasswordHash) {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid email or password"))
	}

	tokenString, err := utils.GenerateJWT(user.ID, user.Email, "user") 
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to generate token"))
	}

	return c.JSON(utils.SuccessResponse(fiber.Map{"token": tokenString}, "Login successful"))
}

// get user profile
func HandleProfile(c *fiber.Ctx) error {
	userID, err := composables.ExtractUserIDFromJWT(c)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}
	user, err := composables.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("User not found"))
	}
	user.PasswordHash = ""
	return c.JSON(utils.SuccessResponse(user, "Current user"))
}

// update user profile
func HandleUpdateProfile(c *fiber.Ctx) error {
	userID, err := composables.ExtractUserIDFromJWT(c)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}

	user, err := composables.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("User not found"))
	}

	formHeader := c.Get("Content-Type")
	isMultipart := formHeader != "" && strings.Contains(strings.ToLower(formHeader), "multipart/form-data")

	if isMultipart {
		user.FirstName = c.FormValue("first_name")
		user.LastName = c.FormValue("last_name")
		user.Nickname = c.FormValue("nickname")
		bio := c.FormValue("bio")
		if bio != "" {
			user.Bio = &bio
		}
		if path, err := composables.HandleAvatarUpload(c, userID); err == nil && path != nil {
			user.Image = path
		}
	} else {
		var input struct {
			FirstName string  `json:"first_name"`
			LastName  string  `json:"last_name"`
			Nickname  string  `json:"nickname"`
			Bio       *string `json:"bio"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(utils.ErrorResponse("Invalid JSON"))
		}
		user.FirstName = input.FirstName
		user.LastName = input.LastName
		user.Nickname = input.Nickname
		user.Bio = input.Bio
	}

	if err := database.DB.Save(user).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to update user"))
	}
	user.PasswordHash = ""
	return c.JSON(utils.SuccessResponse(user, "Profile updated successfully"))
}
