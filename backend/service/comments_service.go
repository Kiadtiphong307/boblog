package service

import (
    "backend/database"
	"backend/models"
	"time"
	"github.com/gofiber/fiber/v2"
    "log"
    "strings"
	"net/url"
)
 

// Get Comments
func HandleGetComments(c *fiber.Ctx) error {
	//  Params & Query
    slugEncoded := c.Params("slug")
	// ทำงานยังไง
    slug, err := url.QueryUnescape(slugEncoded)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid slug encoding"})
    }
	// วิธีการใช้ var , := 
    var article models.Article
	//  ORM 
    if err := database.DB.Where("slug = ?", slug).First(&article).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Article not found"})
    }
    var comments []models.Comment
    if err := database.DB.
        Preload("User").
        Order("created_at desc").
		// SQL 
        Where("article_id = ?", article.ID).
        Find(&comments).Error; err != nil {
		// Find & Frist 
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch comments"})
    }
    return c.JSON(comments)
}

// Create Comment
func HandleCreateComment(c *fiber.Ctx) error {

	slugEncoded := c.Params("slug")
			
	slug, err := url.QueryUnescape(slugEncoded)
	if err != nil {
		log.Println("❌ Invalid slug encoding:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid slug encoding",
		})
	}
	log.Println("📥 POST comment on slug:", slug)


	var article models.Article
	if err := database.DB.Where("slug = ?", slug).First(&article).Error; err != nil {
		log.Println("❌ Article not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		}) 
	}


	var input struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		log.Println("❌ Failed to parse request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		})
	}

	if strings.TrimSpace(input.Content) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		})
	}

	userIDRaw := c.Locals("userID")
	userID, ok := userIDRaw.(uint)
	if !ok {
		log.Println("❌ Invalid or missing user ID in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		})
	}


	comment := models.Comment{
		Content:   input.Content,
		ArticleID: article.ID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		log.Println("❌ Failed to save comment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		})
	}

	log.Printf("✅ Comment created for articleID %d by userID %d\n", article.ID, userID)
	return c.Status(fiber.StatusCreated).JSON(comment)
}

// Update Comment
func HandleUpdateComment(c *fiber.Ctx) error {
	commentID := c.Params("commentId")
	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		})
	}

	userIDRaw := c.Locals("userID")
	userID, ok := userIDRaw.(uint)
	if !ok {
		log.Println("❌ Invalid or missing user ID in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		})
	}

	var comment models.Comment
	if err := database.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		log.Println("❌ Comment not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		})
	}

	if comment.UserID != userID {
		log.Printf("❌ Unauthorized edit attempt: userID %d tried to edit comment by userID %d\n", userID, comment.UserID)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		})
	}

	var input struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		log.Println("❌ Failed to parse request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		})
	}

	if strings.TrimSpace(input.Content) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		})
	}
	
	comment.Content = input.Content
	comment.UpdatedAt = time.Now()

	if err := database.DB.Save(&comment).Error; err != nil {
		log.Println("❌ Failed to update comment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		})
	}

	log.Printf("✅ Comment ID %s updated by userID %d\n", commentID, userID)
	
	//  โหลดข้อมูล User เพื่อ return กลับไป
	if err := database.DB.Preload("User").First(&comment, commentID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		})
	}

	return c.JSON(comment)
}

// Delete Comment
func HandleDeleteComment(c *fiber.Ctx) error {
	//  รับ comment ID จาก URL parameter
	commentID := c.Params("commentId")
	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		})
	}

	//  ตรวจสอบ user จาก middleware
	userIDRaw := c.Locals("userID")
	userID, ok := userIDRaw.(uint)
	if !ok {
		log.Println("❌ Invalid or missing user ID in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		})
	}

	// ค้นหาคอมเมนต์ที่ต้องการลบ
	var comment models.Comment
	if err := database.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		log.Println("❌ Comment not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		})
	}

	//  ตรวจสอบว่าเป็นเจ้าของคอมเมนต์หรือไม่
	if comment.UserID != userID {
		log.Printf("❌ Unauthorized delete attempt: userID %d tried to delete comment by userID %d\n", userID, comment.UserID)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		})
	}

	// ลบคอมเมนต์
	if err := database.DB.Delete(&comment).Error; err != nil {
		log.Println("❌ Failed to delete comment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		})
	}

	log.Printf("✅ Comment ID %s deleted by userID %d\n", commentID, userID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	})
}


