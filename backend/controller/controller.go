package controller

import (
	"blog-db/database"
	"blog-db/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Println("❌ Failed to fetch users:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	log.Printf("✅ Found %d user(s)", len(users))

	var safeUsers []fiber.Map
	for _, u := range users {
		safeUsers = append(safeUsers, fiber.Map{
			"id":        u.ID,
			"username":  u.Username,
			"email":     u.Email,
			"bio":       u.Bio,
			"password":  u.PasswordHash,
			"image":     u.Image,
			"createdAt": u.CreatedAt,
			"updatedAt": u.UpdatedAt,
		})
	}

	if safeUsers == nil {
		safeUsers = []fiber.Map{}
	}
	return c.JSON(safeUsers)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("❌ User ID %s not found", id)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	log.Printf("✅ Found user ID %s", id)

	return c.JSON(fiber.Map{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"bio":       user.Bio,
		"image":     user.Image,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Println("❌ Failed to parse request body for user")
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := database.DB.Create(&user).Error; err != nil {
		log.Println("❌ Failed to create user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	log.Printf("✅ User created: ID %d, Username %s", user.ID, user.Username)

	return c.Status(201).JSON(fiber.Map{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"bio":       user.Bio,
		"image":     user.Image,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("❌ User ID %s not found for update", id)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		log.Println("❌ Failed to parse update data")
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Update เฉพาะ field ที่ส่งมา
	if username, ok := updateData["username"].(string); ok && username != "" {
		user.Username = username
	}
	if email, ok := updateData["email"].(string); ok && email != "" {
		user.Email = email
	}
	if bio, ok := updateData["bio"].(string); ok {
		user.Bio = &bio
	}
	if image, ok := updateData["image"].(string); ok {
		user.Image = &image
	}

	if err := database.DB.Save(&user).Error; err != nil {
		log.Printf("❌ Failed to update user ID %s: %v", id, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	log.Printf("✅ Updated user ID %s", id)

	return c.JSON(fiber.Map{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"bio":       user.Bio,
		"image":     user.Image,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("❌ User ID %s not found for delete", id)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		log.Printf("❌ Failed to delete user ID %s: %v", id, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	log.Printf("✅ Deleted user ID %s", id)
	return c.SendStatus(204)
}
