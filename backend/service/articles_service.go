package controller

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"backend/validation"
	"log"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// คือฟังก์ชันที่จะดึงข้อมูลบทความทั้งหมดจากฐานข้อมูล
func HandleGetAllArticles(c *fiber.Ctx) error {
	var articles []models.Article
	err := database.DB.Preload("Author").Preload("Category").Preload("Tags").Preload("Comments").Find(&articles).Error
	if err != nil {
		log.Println("❌ Error getting all articles:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get articles"))
	}
	return c.JSON(utils.SuccessResponse(articles, "All articles retrieved"))
}

// คือฟังก์ชันที่จะค้นหาบทความตามหมวดหมู่และแท็ก
func HandleSearchArticlesTags(c *fiber.Ctx) error {
	var articles []models.Article
	search := c.Query("search")
	categoryID := c.Query("category_id")
	tx := database.DB.Model(&models.Article{}).
		Joins("LEFT JOIN article_tags ON article_tags.article_id = articles.id").
		Joins("LEFT JOIN tags ON tags.id = article_tags.tags_id").
		Preload("Author").Preload("Category").Preload("Tags").Distinct()

	if search != "" {
		keywords := strings.Fields(strings.ToLower(search))
		for _, kw := range keywords {
			tx = tx.Where(`LOWER(articles.title) LIKE ? OR LOWER(articles.content) LIKE ? OR LOWER(tags.name) LIKE ?`, "%"+kw+"%", "%"+kw+"%", "%"+kw+"%")
		}
	}
	if categoryID != "" {
		tx = tx.Where("articles.category_id = ?", categoryID)
	}
	if err := tx.Order("articles.created_at DESC").Find(&articles).Error; err != nil {
		log.Println("❌ Error filtering articles:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to filter articles"))
	}
	return c.JSON(utils.SuccessResponse(articles, "Filtered articles retrieved"))
}

// คือฟังก์ชันที่จะดึงข้อมูลบทความตาม slug
func HandleGetArticleBySlug(c *fiber.Ctx) error {
	slugParam := c.Params("slug")
	slug, err := url.PathUnescape(slugParam)
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid slug format"))
	}
	var article models.Article
	err = database.DB.Preload("Author").Preload("Category").Preload("Tags").Preload("Comments").First(&article, "slug = ?", slug).Error
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Article not found"))
	}
	return c.JSON(utils.SuccessResponse(article, "Article retrieved successfully"))
}

// คือฟังก์ชันที่จะดึงข้อมูลบทความของผู้ใช้งานปัจจุบัน
func HandleGetMyArticles(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)
	var articles []models.Article
	if err := database.DB.Where("author_id = ?", uint(userID)).Find(&articles).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("ไม่สามารถดึงบทความได้"))
	}
	return c.JSON(utils.SuccessResponse(articles, "My articles retrieved"))
}

// คือฟังก์ชันที่จะสร้างบทความใหม่
func HandleCreateArticle(c *fiber.Ctx) error {
	var input validation.CreateArticleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid input format"))
	}
	if errs := validation.ValidateStructArticle(input); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}

	var titleExists models.Article
	if err := database.DB.Where("title = ?", input.Title).First(&titleExists).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"errors": map[string]string{"title": "ชื่อบทความนี้มีอยู่แล้ว"}})
	}
	var slugExists models.Article
	if err := database.DB.Where("slug = ?", input.Slug).First(&slugExists).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"errors": map[string]string{"slug": "Slug นี้ถูกใช้แล้ว"}})
	}

	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	authorID := uint(claims["id"].(float64))

	var category models.Category
	if err := database.DB.Where("name = ?", input.CategoryName).First(&category).Error; err != nil {
		category = models.Category{Name: input.CategoryName}
		database.DB.Create(&category)
	}

	tags := []models.Tags{}
	for _, name := range input.TagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var tag models.Tags
		if err := database.DB.Where("name = ?", name).First(&tag).Error; err != nil {
			tag = models.Tags{Name: name}
			database.DB.Create(&tag)
		}
		tags = append(tags, tag)
	}

	article := models.Article{
		Title:      input.Title,
		Slug:       input.Slug,
		Content:    input.Content,
		AuthorID:   authorID,
		CategoryID: category.ID,
		Tags:       tags,
	}
	if err := database.DB.Create(&article).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("ไม่สามารถสร้างบทความ"))
	}
	return c.Status(201).JSON(utils.SuccessResponse(article, "✅ สร้างบทความสำเร็จ"))
}

// คือฟังก์ชันที่จะแก้ไขบทความ
func HandleUpdateArticle(c *fiber.Ctx) error {
	slug := c.Params("slug")
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))
	var article models.Article
	if err := database.DB.Preload("Tags").First(&article, "slug = ?", slug).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("ไม่พบบทความ"))
	}
	if article.AuthorID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("คุณไม่มีสิทธิ์แก้ไขบทความนี้"))
	}

	var input validation.UpdateArticleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("ข้อมูลไม่ถูกต้อง"))
	}

	if input.Title != nil {
		article.Title = *input.Title
	}
	if input.Content != nil {
		article.Content = *input.Content
	}
	if input.CategoryID != nil {
		article.CategoryID = *input.CategoryID
	}
	if err := database.DB.Save(&article).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("บันทึกบทความล้มเหลว"))
	}

	var createdTags []models.Tags
	for _, name := range input.NewTags {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var existing models.Tags
		if err := database.DB.Where("LOWER(name) = ?", strings.ToLower(name)).First(&existing).Error; err == nil {
			continue
		} else if err != gorm.ErrRecordNotFound {
			return c.Status(500).JSON(utils.ErrorResponse("ตรวจสอบแท็กใหม่ล้มเหลว"))
		}
		newTag := models.Tags{Name: name}
		database.DB.Create(&newTag)
		createdTags = append(createdTags, newTag)
	}

	allTagIDs := input.TagIDs
	for _, tag := range createdTags {
		allTagIDs = append(allTagIDs, tag.ID)
	}

	if input.TagIDs != nil || len(createdTags) > 0 {
		if len(allTagIDs) > 0 {
			var allTags []models.Tags
			database.DB.Where("id IN ?", allTagIDs).Find(&allTags)
			database.DB.Model(&article).Association("Tags").Replace(&allTags)
		} else {
			database.DB.Model(&article).Association("Tags").Clear()
		}
	}
	return c.JSON(utils.SuccessResponse(article, "แก้ไขบทความเรียบร้อยแล้ว"))
}

// คือฟังก์ชันที่จะลบบทความ
func HandleDeleteArticle(c *fiber.Ctx) error {
	slug, err := url.QueryUnescape(c.Params("slug"))
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Slug ไม่ถูกต้อง"))
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))
	var article models.Article
	if err := database.DB.Preload("Tags").First(&article, "slug = ?", slug).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("ไม่พบบทความ"))
	}
	if article.AuthorID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("คุณไม่มีสิทธิ์ลบบทความนี้"))
	}
	database.DB.Model(&article).Association("Tags").Clear()
	if err := database.DB.Delete(&article).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("ลบบทความไม่สำเร็จ"))
	}
	return c.JSON(utils.SuccessResponse(nil, "ลบบทความเรียบร้อยแล้ว"))
}
