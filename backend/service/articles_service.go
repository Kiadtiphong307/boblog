package service

import (
	"backend/composables"
	"backend/database"
	"backend/models"
	"backend/utils"
	"backend/validation"
	"log"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)



// คือฟังก์ชันที่จะดึงข้อมูลบทความทั้งหมดจากฐานข้อมูล
func HandleGetAllArticles(c *fiber.Ctx) error {
		var articles []models.Article
	if err := database.DB.Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments").
		Find(&articles).Error; err != nil {
		log.Println("❌ Error getting all articles:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get articles"))
	}
	return c.JSON(utils.SuccessResponse(articles, "All articles retrieved"))
}

// คือฟังก์ชันที่จะค้นหาบทความทั้งหมดจากฐานข้อมูล
func HandleSearchArticlesTags(c *fiber.Ctx) error {
	var articles []models.Article
	search, categoryID := c.Query("search"), c.Query("category_id")
	
	tx := database.DB.Model(&models.Article{}).
		Joins("LEFT JOIN article_tags ON article_tags.article_id = articles.id").
		Joins("LEFT JOIN tags ON tags.id = article_tags.tags_id").
		Preload("Author").
		Preload("Category").
		Preload("Tags").Distinct()

	if search != "" {
		for _, kw := range strings.Fields(strings.ToLower(search)) {
			tx = tx.Where(`LOWER(articles.title) LIKE ? OR LOWER(articles.content) LIKE ? OR LOWER(tags.name) LIKE ?`, 
				"%"+kw+"%", "%"+kw+"%", "%"+kw+"%")
		}
	}
	if categoryID != "" {
		tx = tx.Where("articles.category_id = ?", categoryID)
	}
	
	if err := tx.Order("articles.created_at DESC").Find(&articles).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to filter articles"))
	}
	return c.JSON(utils.SuccessResponse(articles, "Filtered articles retrieved"))
}

// คือฟังก์ชันที่จะดึงข้อมูลบตามslug
func HandleGetArticleBySlug(c *fiber.Ctx) error {
	slug, err := url.PathUnescape(c.Params("slug"))
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid slug format"))
	}
	
	var article models.Article
	if err := database.DB.
	Preload("Author").
	Preload("Category").
	Preload("Tags").
	Preload("Comments").
	First(&article, "slug = ?", slug).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Article not found"))
	}
	
	return c.JSON(utils.SuccessResponse(article, "Article retrieved successfully"))
}

// คือฟังก์ชันที่จะดึงข้อมูลบทความ เจ้าของบทความ
func HandleGetMyArticles(c *fiber.Ctx) error {
	userID, err := composables.GetCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}
	var articles []models.Article
	if err := database.DB.Where("author_id = ?", userID).Find(&articles).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("failed to get articles"))
	}
	return c.JSON(utils.SuccessResponse(articles, "My articles retrieved"))
}

// คือฟังก์ชันที่จะสร้างบทความ
func HandleCreateArticle(c *fiber.Ctx) error {
	var input validation.CreateArticleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid input format"))
	}
	if errs := validation.ValidateStructArticle(input); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	var exists models.Article
	if err := database.DB.Where("title = ? OR slug = ?", input.Title, input.Slug).First(&exists).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"errors": map[string]string{"title": "title or slug already exists"}})
	}
	userID, err := composables.GetCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}
	category, err := composables.FindOrCreateCategory(input.CategoryName)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to create category"))
	}
	tags, err := composables.FindOrCreateTags(input.TagNames)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to create tags"))
	}
	article := models.Article{
		Title:      input.Title,
		Slug:       input.Slug,
		Content:    input.Content,
		AuthorID:   userID,
		CategoryID: category.ID,
		Tags:       tags,
	}
	if err := database.DB.Create(&article).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("create article failed"))
	}
	return c.Status(201).JSON(utils.SuccessResponse(article, "create article success"))
}

// คือฟังก์ชันที่จะอัปเดตบทความ
func HandleUpdateArticle(c *fiber.Ctx) error {
	slug := c.Params("slug")
	userID, err := composables.GetCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}
	var article models.Article
	if err := database.DB.Preload("Tags").First(&article, "slug = ?", slug).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("article not found"))
	}
	if article.AuthorID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("you don't have permission to update this article"))
	}
	var input validation.UpdateArticleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("invalid data"))
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
		return c.Status(500).JSON(utils.ErrorResponse("save article failed")) 
	}
	// สร้าง tag ใหม่
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
			return c.Status(500).JSON(utils.ErrorResponse("check new tag failed"))
		}
		newTag := models.Tags{Name: name}
		database.DB.Create(&newTag)
		createdTags = append(createdTags, newTag)
	}
	// รวมค่า tag ใหม่กับ tag เดิม
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
	return c.JSON(utils.SuccessResponse(article, "update article success"))
}

// คือฟังก์ชันที่จะลบบทความ
func HandleDeleteArticle(c *fiber.Ctx) error {
	slug, err := url.QueryUnescape(c.Params("slug"))
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("invalid slug"))
	}
	userID, err := composables.GetCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}
	var article models.Article
	if err := database.DB.Preload("Tags").First(&article, "slug = ?", slug).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("article not found"))
	}
	if article.AuthorID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("you don't have permission to delete this article"))
	}
	database.DB.Model(&article).Association("Tags").Clear()
	if err := database.DB.Delete(&article).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("delete article failed"))
	}
	return c.JSON(utils.SuccessResponse(nil, "delete article success"))
}
