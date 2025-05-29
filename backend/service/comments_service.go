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
	// ✅ รับ slug ที่ encode จาก URL
	slugEncoded := c.Params("slug")

	// ✅ decode slug เพื่อให้ match กับในฐานข้อมูล
	slug, err := url.QueryUnescape(slugEncoded)
	if err != nil {
		log.Println("❌ Invalid slug encoding:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid slug encoding",
		})
	}
	log.Println("📥 POST comment on slug:", slug)

	// ✅ ค้นหาบทความจาก slug ที่ decode แล้ว
	var article models.Article
	if err := database.DB.Where("slug = ?", slug).First(&article).Error; err != nil {
		log.Println("❌ Article not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ไม่พบบทความ",
		})
	}

	// ✅ รับข้อมูลคอมเมนต์จาก body
	var input struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		log.Println("❌ Failed to parse request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบข้อมูลไม่ถูกต้อง",
		})
	}

	// ✅ ตรวจสอบว่า content ไม่ว่าง
	if strings.TrimSpace(input.Content) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ต้องกรอกข้อความความคิดเห็น",
		})
	}

	// ✅ ตรวจสอบ user จาก middleware (ควรแน่ใจว่า middleware เซ็ต userID ให้ถูกต้อง)
	userIDRaw := c.Locals("userID")
	userID, ok := userIDRaw.(uint)
	if !ok {
		log.Println("❌ Invalid or missing user ID in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: ไม่พบข้อมูลผู้ใช้งาน",
		})
	}

	// ✅ สร้าง comment ใหม่
	comment := models.Comment{
		Content:   input.Content,
		ArticleID: article.ID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		log.Println("❌ Failed to save comment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถบันทึกความคิดเห็นได้",
		})
	}

	log.Printf("✅ Comment created for articleID %d by userID %d\n", article.ID, userID)
	return c.Status(fiber.StatusCreated).JSON(comment)
}

// Update Comment
func HandleUpdateComment(c *fiber.Ctx) error {
	// ✅ รับ comment ID จาก URL parameter
	commentID := c.Params("commentId")
	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ต้องระบุ ID ของความคิดเห็น",
		})
	}

	// ✅ ตรวจสอบ user จาก middleware
	userIDRaw := c.Locals("userID")
	userID, ok := userIDRaw.(uint)
	if !ok {
		log.Println("❌ Invalid or missing user ID in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: ไม่พบข้อมูลผู้ใช้งาน",
		})
	}

	// ✅ ค้นหาคอมเมนต์ที่ต้องการแก้ไข
	var comment models.Comment
	if err := database.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		log.Println("❌ Comment not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ไม่พบความคิดเห็น",
		})
	}

	// ✅ ตรวจสอบว่าเป็นเจ้าของคอมเมนต์หรือไม่
	if comment.UserID != userID {
		log.Printf("❌ Unauthorized edit attempt: userID %d tried to edit comment by userID %d\n", userID, comment.UserID)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "คุณไม่มีสิทธิ์แก้ไขความคิดเห็นนี้",
		})
	}

	// ✅ รับข้อมูลใหม่จาก body
	var input struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		log.Println("❌ Failed to parse request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบข้อมูลไม่ถูกต้อง",
		})
	}

	// ✅ ตรวจสอบว่า content ไม่ว่าง
	if strings.TrimSpace(input.Content) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ต้องกรอกข้อความความคิดเห็น",
		})
	}

	// ✅ อัพเดตคอมเมนต์
	comment.Content = input.Content
	comment.UpdatedAt = time.Now()

	if err := database.DB.Save(&comment).Error; err != nil {
		log.Println("❌ Failed to update comment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถแก้ไขความคิดเห็นได้",
		})
	}

	log.Printf("✅ Comment ID %s updated by userID %d\n", commentID, userID)
	
	// ✅ โหลดข้อมูล User เพื่อ return กลับไป
	if err := database.DB.Preload("User").First(&comment, commentID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถโหลดข้อมูลความคิดเห็นได้",
		})
	}

	return c.JSON(comment)
}

// Delete Comment
func HandleDeleteComment(c *fiber.Ctx) error {
	// ✅ รับ comment ID จาก URL parameter
	commentID := c.Params("commentId")
	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ต้องระบุ ID ของความคิดเห็น",
		})
	}

	// ✅ ตรวจสอบ user จาก middleware
	userIDRaw := c.Locals("userID")
	userID, ok := userIDRaw.(uint)
	if !ok {
		log.Println("❌ Invalid or missing user ID in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: ไม่พบข้อมูลผู้ใช้งาน",
		})
	}

	// ✅ ค้นหาคอมเมนต์ที่ต้องการลบ
	var comment models.Comment
	if err := database.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		log.Println("❌ Comment not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ไม่พบความคิดเห็น",
		})
	}

	// ✅ ตรวจสอบว่าเป็นเจ้าของคอมเมนต์หรือไม่
	if comment.UserID != userID {
		log.Printf("❌ Unauthorized delete attempt: userID %d tried to delete comment by userID %d\n", userID, comment.UserID)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "คุณไม่มีสิทธิ์ลบความคิดเห็นนี้",
		})
	}

	// ✅ ลบคอมเมนต์
	if err := database.DB.Delete(&comment).Error; err != nil {
		log.Println("❌ Failed to delete comment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถลบความคิดเห็นได้",
		})
	}

	log.Printf("✅ Comment ID %s deleted by userID %d\n", commentID, userID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ลบความคิดเห็นเรียบร้อยแล้ว",
	})
}


