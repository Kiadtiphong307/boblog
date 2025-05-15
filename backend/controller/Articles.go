package controller

import (
	"blog-db/database"
	"blog-db/models"
	"blog-db/utils"
	"log"
	"net/url"
	"strings"


	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)




// Get all articles
func GetAllArticles(c *fiber.Ctx) error {
	var articles []models.Article

	err := database.DB.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments").
		Find(&articles).Error

	if err != nil {
		log.Println("❌ Error getting all articles:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get articles"))
	}

	log.Println("✅ Retrieved all articles")
	return c.JSON(utils.SuccessResponse(articles, "All articles retrieved"))
}

// filter articles and tags
func SearchArticlesTags(c *fiber.Ctx) error {
	var articles []models.Article

	search := c.Query("search")
	categoryID := c.Query("category_id")

	tx := database.DB.
		Model(&models.Article{}).
		Joins("LEFT JOIN article_tags ON article_tags.article_id = articles.id").
		Joins("LEFT JOIN tags ON tags.id = article_tags.tags_id").
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Distinct()

	// ✅ ถ้ามี keyword → ค้นหาทั้ง title, content, tags.name
	if search != "" {
		keywords := strings.Fields(strings.ToLower(search)) // แยกเป็นคำ เช่น “go fiber”
		for _, kw := range keywords {
			tx = tx.Where(`
				LOWER(articles.title) LIKE ?
				OR LOWER(articles.content) LIKE ?
				OR LOWER(tags.name) LIKE ?
			`, "%"+kw+"%", "%"+kw+"%", "%"+kw+"%")
		}
	}

	// ✅ กรองหมวดหมู่ (ถ้ามี)
	if categoryID != "" {
		tx = tx.Where("articles.category_id = ?", categoryID)
	}

	// ✅ ดึงผลลัพธ์
	if err := tx.Order("articles.created_at DESC").Find(&articles).Error; err != nil {
		log.Println("❌ Error filtering articles:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to filter articles"))
	}

	return c.JSON(utils.SuccessResponse(articles, "Filtered articles retrieved"))
}

// Get article by slug
func GetArticleBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	var article models.Article

	err := database.DB.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments").
		First(&article, "slug = ?", slug).Error

	if err != nil {
		log.Println("❌ Article not found:", err)
		return c.Status(404).JSON(utils.ErrorResponse("Article not found"))
	}

	return c.JSON(utils.SuccessResponse(article, "Article retrieved successfully"))
}

// Get my articles
func GetMyArticles(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	var articles []models.Article
	if err := database.DB.Where("author_id = ?", uint(userID)).Find(&articles).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "ไม่สามารถดึงบทความได้",
		})
	}

	return c.JSON(fiber.Map{
		"data": articles,
	})
}

// CreateArticle creates a new article
func CreateArticle(c *fiber.Ctx) error {
	var input struct {
		Title        string   `json:"title" validate:"required"`
		Slug         string   `json:"slug" validate:"required"`
		Content      string   `json:"content" validate:"required"`
		CategoryName string   `json:"category_name" validate:"required"` // ✅ รับชื่อหมวดหมู่
		TagNames     []string `json:"tag_names"`
	}

	// ✅ แปลง JSON เป็น struct
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid input"))
	}

	// ✅ ตรวจสอบข้อมูล
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Validation failed: " + err.Error()))
	}

	// ✅ ตรวจ slug ซ้ำ
	var exist models.Article
	if err := database.DB.Where("slug = ?", input.Slug).First(&exist).Error; err == nil {
		return c.Status(400).JSON(utils.ErrorResponse("Slug already exists"))
	}

	// ✅ ดึง user ID จาก JWT
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	authorID := uint(claims["id"].(float64))

	// ✅ ตรวจสอบหรือสร้างหมวดหมู่
	var category models.Category
	if err := database.DB.Where("name = ?", input.CategoryName).First(&category).Error; err != nil {
		category = models.Category{Name: input.CategoryName}
		if err := database.DB.Create(&category).Error; err != nil {
			return c.Status(500).JSON(utils.ErrorResponse("Failed to create category: " + input.CategoryName))
		}
	}

	// ✅ ตรวจสอบหรือสร้าง tag
	var tags []models.Tags
	for _, name := range input.TagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var tag models.Tags
		if err := database.DB.Where("name = ?", name).First(&tag).Error; err != nil {
			tag = models.Tags{Name: name}
			if err := database.DB.Create(&tag).Error; err != nil {
				return c.Status(500).JSON(utils.ErrorResponse("Failed to create tag: " + name))
			}
		}
		tags = append(tags, tag)
	}

	// ✅ บันทึกบทความ
	article := models.Article{
		Title:      input.Title,
		Slug:       input.Slug,
		Content:    input.Content,
		AuthorID:   authorID,
		CategoryID: category.ID,
		Tags:       tags,
	}

	if err := database.DB.Create(&article).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to create article"))
	}

	return c.Status(201).JSON(utils.SuccessResponse(article, "✅ Article created successfully"))
}

// UpdateArticle updates an existing article
func UpdateArticle(c *fiber.Ctx) error {
	slug := c.Params("slug")

	// ✅ อ่านข้อมูลผู้ใช้จาก JWT
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	// ✅ โหลดบทความพร้อมแท็กเดิม
	var article models.Article
	if err := database.DB.Preload("Tags").First(&article, "slug = ?", slug).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("ไม่พบบทความ"))
	}

	// ✅ ตรวจสอบสิทธิ์เจ้าของบทความ
	if article.AuthorID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("คุณไม่มีสิทธิ์แก้ไขบทความนี้"))
	}

	// ✅ รับ input จาก client
	var input struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		CategoryID uint     `json:"category_id"`
		TagIDs     []uint   `json:"tag_ids"`
		NewTags    []string `json:"new_tags"` // 👈 รับแท็กใหม่
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("ข้อมูลไม่ถูกต้อง"))
	}

	// ✅ ตรวจสอบหมวดหมู่
	var category models.Category
	if err := database.DB.First(&category, input.CategoryID).Error; err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("หมวดหมู่ไม่ถูกต้อง"))
	}

	// ✅ อัปเดตข้อมูลหลัก
	article.Title = input.Title
	article.Content = input.Content
	article.CategoryID = input.CategoryID

	if err := database.DB.Save(&article).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("บันทึกบทความล้มเหลว"))
	}

	// ✅ ตรวจสอบและสร้างแท็กใหม่
	var createdTags []models.Tags
	for _, name := range input.NewTags {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		var existing models.Tags
		if err := database.DB.Where("LOWER(name) = ?", strings.ToLower(name)).First(&existing).Error; err == nil {
			// พบอยู่แล้ว → ข้าม
			continue
		} else if err != gorm.ErrRecordNotFound {
			// Error อื่น
			return c.Status(500).JSON(utils.ErrorResponse("ตรวจสอบแท็กใหม่ล้มเหลว"))
		}

		newTag := models.Tags{Name: name}
		if err := database.DB.Create(&newTag).Error; err != nil {
			return c.Status(500).JSON(utils.ErrorResponse("สร้างแท็กใหม่ล้มเหลว"))
		}
		createdTags = append(createdTags, newTag)
	}

	// ✅ รวม TagIDs เดิม + ใหม่
	allTagIDs := input.TagIDs
	for _, tag := range createdTags {
		allTagIDs = append(allTagIDs, tag.ID)
	}

	// ✅ แทนที่แท็กทั้งหมด
	if len(allTagIDs) > 0 {
		var allTags []models.Tags
		if err := database.DB.Where("id IN ?", allTagIDs).Find(&allTags).Error; err != nil {
			return c.Status(500).JSON(utils.ErrorResponse("โหลดแท็กใหม่ล้มเหลว"))
		}
		if err := database.DB.Model(&article).Association("Tags").Replace(&allTags); err != nil {
			return c.Status(500).JSON(utils.ErrorResponse("อัปเดตแท็กล้มเหลว"))
		}
	} else {
		// ✅ ถ้าไม่มีแท็ก → ลบความสัมพันธ์
		if err := database.DB.Model(&article).Association("Tags").Clear(); err != nil {
			return c.Status(500).JSON(utils.ErrorResponse("ลบแท็กเก่าไม่สำเร็จ"))
		}
	}

	return c.JSON(utils.SuccessResponse(article, "แก้ไขบทความเรียบร้อยแล้ว"))
}

// DeleteArticle deletes an article by slug (only by the author)
func DeleteArticle(c *fiber.Ctx) error {
	encodedSlug := c.Params("slug")
	slug, err := url.QueryUnescape(encodedSlug)
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Slug ไม่ถูกต้อง"))
	}

	// ดึงข้อมูล user จาก JWT
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	// ค้นหาบทความจาก slug พร้อม preload แท็ก
	var article models.Article
	if err := database.DB.Preload("Tags").First(&article, "slug = ?", slug).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("ไม่พบบทความ"))
	}

	// ตรวจสอบสิทธิ์เจ้าของบทความ
	if article.AuthorID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("คุณไม่มีสิทธิ์ลบบทความนี้"))
	}

	// ลบความสัมพันธ์กับแท็ก
	if err := database.DB.Model(&article).Association("Tags").Clear(); err != nil {
		log.Println("❌ ลบความสัมพันธ์กับ Tags ไม่สำเร็จ:", err)
		return c.Status(500).JSON(utils.ErrorResponse("ลบแท็กของบทความไม่สำเร็จ"))
	}

	// ลบบทความ
	if err := database.DB.Delete(&article).Error; err != nil {
		log.Println("🔥 ลบบทความไม่สำเร็จ:", err)
		return c.Status(500).JSON(utils.ErrorResponse("ลบบทความไม่สำเร็จ"))
	}

	return c.JSON(utils.SuccessResponse(nil, "ลบบทความเรียบร้อยแล้ว"))
}
