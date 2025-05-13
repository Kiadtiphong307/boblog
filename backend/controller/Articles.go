package controller

import (
	"blog-db/database"
	"blog-db/models"
	"blog-db/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

// Search articles
func SearchArticles(c *fiber.Ctx) error {
	var articles []models.Article

	search := c.Query("search")
	categoryID := c.Query("category_id")

	tx := database.DB.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments")

	if search != "" {
		tx = tx.Where("title LIKE ?", "%"+search+"%")
	}

	if categoryID != "" {
		tx = tx.Where("category_id = ?", categoryID)
	}

	if err := tx.Find(&articles).Error; err != nil {
		log.Println("❌ Error filtering articles:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to filter articles"))
	}

	log.Println("✅ Retrieved filtered articles")
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
		Title      string   `json:"title" validate:"required"`
		Slug       string   `json:"slug" validate:"required"`
		Content    string   `json:"content" validate:"required"`
		CategoryID uint     `json:"category_id" validate:"required"`
		Tags       []string `json:"tags"`
	}

	// 1. รับข้อมูลจาก request body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid input"))
	}

	// 2. ตรวจสอบความถูกต้องของข้อมูล
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Validation failed: " + err.Error()))
	}

	// 3. ตรวจสอบ slug ไม่ซ้ำ
	var existArticle models.Article
	if err := database.DB.Where("slug = ?", input.Slug).First(&existArticle).Error; err == nil {
		return c.Status(400).JSON(utils.ErrorResponse("Slug already exists"))
	}

	// 4. ดึง AuthorID จาก JWT token
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	authorID := uint(claims["id"].(float64))

	// 5. เตรียมข้อมูล tags (สร้างใหม่ถ้ายังไม่มี)
	var tags []models.Tag
	for _, tagName := range input.Tags {
		var tag models.Tag
		if err := database.DB.Where("name = ?", tagName).First(&tag).Error; err != nil {
			tag = models.Tag{Name: tagName}
			if err := database.DB.Create(&tag).Error; err != nil {
				return c.Status(500).JSON(utils.ErrorResponse("Failed to create tag: " + tagName))
			}
		}
		tags = append(tags, tag)
	}

	// 6. สร้าง Article ใหม่
	article := models.Article{
		Title:      input.Title,
		Slug:       input.Slug,
		Content:    input.Content,
		AuthorID:   authorID,
		CategoryID: input.CategoryID,
		Tags:       tags,
	}

	if err := database.DB.Create(&article).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to create article"))
	}

	return c.Status(201).JSON(utils.SuccessResponse(article, "Article created successfully"))
}


// UpdateArticle updates an existing article
func UpdateArticle(c *fiber.Ctx) error {
    slug := c.Params("slug")
    user := c.Locals("user").(*jwt.Token)
    claims := user.Claims.(jwt.MapClaims)
    userID := uint(claims["id"].(float64))

    // ค้นหาบทความจาก slug
    var article models.Article
    if err := database.DB.Preload("Tags").First(&article, "slug = ?", slug).Error; err != nil {
        return c.Status(404).JSON(utils.ErrorResponse("ไม่พบบทความ"))
    }

    // ตรวจสอบสิทธิ์ผู้ใช้
    if article.AuthorID != userID {
        return c.Status(403).JSON(utils.ErrorResponse("คุณไม่มีสิทธิ์แก้ไขบทความนี้"))
    }

    // รับ input จาก client
    var input struct {
        Title      string `json:"title"`
        Content    string `json:"content"`
        CategoryID uint   `json:"category_id"`
        TagIDs     []uint `json:"tag_ids"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(utils.ErrorResponse("ข้อมูลไม่ถูกต้อง"))
    }

    // ตรวจสอบ category ที่เลือกมีอยู่จริง
    var category models.Category
    if err := database.DB.First(&category, input.CategoryID).Error; err != nil {
        return c.Status(400).JSON(utils.ErrorResponse("หมวดหมู่ไม่ถูกต้อง"))
    }

    // ดึง tag ที่เลือกทั้งหมด
    var tags []models.Tag
    if len(input.TagIDs) > 0 {
        if err := database.DB.Where("id IN ?", input.TagIDs).Find(&tags).Error; err != nil {
            return c.Status(400).JSON(utils.ErrorResponse("ไม่สามารถโหลดแท็กได้"))
        }
    }

    // อัปเดตบทความ
    article.Title = input.Title
    article.Content = input.Content
    article.CategoryID = input.CategoryID

    // อัปเดต Many2Many relation (แท็ก)
    if err := database.DB.Model(&article).Association("Tags").Replace(&tags); err != nil {
        return c.Status(500).JSON(utils.ErrorResponse("อัปเดตแท็กล้มเหลว"))
    }

    if err := database.DB.Save(&article).Error; err != nil {
        return c.Status(500).JSON(utils.ErrorResponse("บันทึกบทความล้มเหลว"))
    }

    return c.JSON(utils.SuccessResponse(article, "แก้ไขบทความเรียบร้อยแล้ว"))
}


// DeleteArticle deletes an article by slug (only by the author)
func DeleteArticle(c *fiber.Ctx) error {
	slug := c.Params("slug")

	// รับข้อมูล user จาก JWT
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	// หาบทความตาม slug
	var article models.Article
	if err := database.DB.First(&article, "slug = ?", slug).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("ไม่พบบทความ"))
	}

	// ตรวจสอบสิทธิ์
	if article.AuthorID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("คุณไม่มีสิทธิ์ลบบทความนี้"))
	}

	// ลบบทความโดยไม่ยุ่งกับ tags
	if err := database.DB.Delete(&article).Error; err != nil {
		log.Println("🔥 ลบไม่สำเร็จ:", err)
		return c.Status(500).JSON(utils.ErrorResponse("ลบบทความไม่สำเร็จ"))
	}

	return c.JSON(utils.SuccessResponse(nil, "ลบบทความเรียบร้อยแล้ว"))
}

