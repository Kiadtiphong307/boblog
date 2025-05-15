package controller

import (
    "blog-db/database"
	"blog-db/models"
	"time"
	"github.com/gofiber/fiber/v2"
    "log"
    "strings"
	"net/url"
)


// Get Comments
func GetComments(c *fiber.Ctx) error {
    slugEncoded := c.Params("slug")
    slug, err := url.QueryUnescape(slugEncoded)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid slug encoding"})
    }

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


