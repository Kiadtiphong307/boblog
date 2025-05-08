package controller

import (
	"blog-db/database"
	"blog-db/models"
	"log"
	"github.com/go-playground/validator/v10"
	"blog-db/utils"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

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

// CreateUser creates a new user
func CreateUser(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        log.Println("❌ Failed to parse request body:", err)
        return c.Status(400).JSON(utils.ErrorResponse("Cannot parse JSON"))
    }

    // ✅ Validation
    if err := validate.Struct(&user); err != nil {
        log.Println("❌ Validation failed:", err)
        return c.Status(400).JSON(utils.ErrorResponse("Validation failed: " + err.Error()))
    }

    if result := database.DB.Create(&user); result.Error != nil {
        log.Println("❌ Failed to create user:", result.Error)
        return c.Status(500).JSON(utils.ErrorResponse("Failed to create user"))
    }

    log.Println("✅ Created new user:", user)
    return c.Status(201).JSON(utils.SuccessResponse(user, "User created successfully"))
}

// UpdateUser updates an existing user by ID
func UpdateUser(c *fiber.Ctx) error {
    id := c.Params("id")
    var user models.User
    if err := database.DB.First(&user, id).Error; err != nil {
        log.Println("❌ User not found with ID:", id)
        return c.Status(404).JSON(utils.ErrorResponse("User not found"))
    }

    var updateData models.User
    if err := c.BodyParser(&updateData); err != nil {
        log.Println("❌ Failed to parse update data:", err)
        return c.Status(400).JSON(utils.ErrorResponse("Cannot parse JSON"))
    }

    // ✅ Validation
    if err := validate.Struct(&updateData); err != nil {
        log.Println("❌ Validation failed:", err)
        return c.Status(400).JSON(utils.ErrorResponse("Validation failed: " + err.Error()))
    }

    user.Username = updateData.Username
    user.Email = updateData.Email
    database.DB.Save(&user)

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