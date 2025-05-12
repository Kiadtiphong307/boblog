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
func GetArticles(c *fiber.Ctx) error {
	var articles []models.Article

	err := database.DB.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments").
		Find(&articles).Error

	if err != nil {
		log.Println("❌ Error getting articles:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get articles"))
	}

	log.Println("✅ Retrieved all articles with relations")
	return c.JSON(utils.SuccessResponse(articles, "Articles retrieved successfully"))
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
