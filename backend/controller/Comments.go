package controller

import (
    "blog-db/database"
	"blog-db/models"
	"time"
	"github.com/gofiber/fiber/v2"
    "log"
    "strings"
)


// Get Comments
func GetComments(c *fiber.Ctx) error {
    slug := c.Params("slug")
    var article models.Article
    if err := database.DB.Where("slug = ?", slug).First(&article).Error; err != nil {
      return c.Status(404).JSON(fiber.Map{"error": "Article not found"})
    }
  
    var comments []models.Comment
    if err := database.DB.
    Preload("User").
    Order("created_at desc").
    Where("article_id = ?", article.ID).
    Find(&comments).Error; err != nil {
      return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch comments"})
    }
  
    return c.JSON(comments)
  }

// Create Comment
func CreateComment(c *fiber.Ctx) error {
	slug := c.Params("slug")
	log.Println("📥 POST comment on slug:", slug)

	// ค้นหาบทความจาก slug
	var article models.Article
	if err := database.DB.Where("slug = ?", slug).First(&article).Error; err != nil {
		log.Println("❌ Article not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Article not found",
		})
	}

	// รับข้อมูลจาก body
	var input struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		log.Println("❌ Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// ตรวจสอบว่า content ว่างหรือไม่
	if strings.TrimSpace(input.Content) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Comment content is required",
		})
	}

	// ป้องกัน panic: ตรวจสอบ userID ที่ middleware เซ็ตไว้
	userIDRaw := c.Locals("userID")
	userID, ok := userIDRaw.(uint)
	if !ok {
		log.Println("❌ Invalid or missing user ID in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid user context",
		})
	}

	// สร้าง comment ใหม่
	comment := models.Comment{
		Content:   input.Content,
		ArticleID: article.ID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		log.Println("❌ Failed to create comment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save comment",
		})
	}

	log.Printf("✅ Comment created for articleID %d by userID %d\n", article.ID, userID)
	return c.Status(fiber.StatusCreated).JSON(comment)
}


