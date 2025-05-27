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
	"github.com/go-playground/validator/v10"
)

// Get all articles
func GetAllArticles(c *fiber.Ctx) error {
	var articles []models.Article

	//q
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

	// search articles and tags by keyword
	if search != "" {
		keywords := strings.Fields(strings.ToLower(search)) // split into words, e.g. "go fiber"
		for _, kw := range keywords {
			tx = tx.Where(`
				LOWER(articles.title) LIKE ?
				OR LOWER(articles.content) LIKE ?
				OR LOWER(tags.name) LIKE ?
			`, "%"+kw+"%", "%"+kw+"%", "%"+kw+"%")
		}
	}

	// filter category if exists
	if categoryID != "" {
		tx = tx.Where("articles.category_id = ?", categoryID)
	}

	// get results
	if err := tx.Order("articles.created_at DESC").Find(&articles).Error; err != nil {
		log.Println("❌ Error filtering articles:", err)
		return c.Status(500).JSON(utils.ErrorResponse("Failed to filter articles"))
	}

	return c.JSON(utils.SuccessResponse(articles, "Filtered articles retrieved"))
}

// Get article by slug
func GetArticleBySlug(c *fiber.Ctx) error {
	// decode slug that is encoded as %E0%B8...
	slugParam := c.Params("slug")
	slug, err := url.PathUnescape(slugParam)
	if err != nil {
		log.Println("❌ Invalid slug format:", err)
		return c.Status(400).JSON(utils.ErrorResponse("Invalid slug format"))
	}

	var article models.Article

	err = database.DB.
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
		CategoryName string   `json:"category_name" validate:"required"`
		TagNames     []string `json:"tag_names"`
	}

	// แปลง JSON จาก body เป็น struct
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input format",
		})
	}

	// ตรวจสอบ validation ด้วย validator
	if err := validate.Struct(input); err != nil {
		validationErrors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				validationErrors[field] = "จำเป็นต้องกรอก"
			default:
				validationErrors[field] = "ข้อมูลไม่ถูกต้อง"
			}
		}
		return c.Status(400).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	// ตรวจสอบชื่อบทความซ้ำ
	var titleExists models.Article
	if err := database.DB.Where("title = ?", input.Title).First(&titleExists).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]string{
				"title": "ชื่อบทความนี้มีอยู่แล้ว",
			},
		})
	}

	// ตรวจสอบ slug ซ้ำ
	var slugExists models.Article
	if err := database.DB.Where("slug = ?", input.Slug).First(&slugExists).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"errors": map[string]string{
				"slug": "Slug นี้ถูกใช้แล้ว",
			},
		})
	}

	// ดึง user ID จาก JWT
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized: Token missing or invalid",
		})
	}
	claims := userToken.Claims.(jwt.MapClaims)
	authorID := uint(claims["id"].(float64))

	// ตรวจสอบหรือสร้าง category
	var category models.Category
	if err := database.DB.Where("name = ?", input.CategoryName).First(&category).Error; err != nil {
		category = models.Category{Name: input.CategoryName}
		if err := database.DB.Create(&category).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "ไม่สามารถสร้างหมวดหมู่: " + input.CategoryName,
			})
		}
	}

	// ตรวจสอบหรือสร้างแท็ก
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
				return c.Status(500).JSON(fiber.Map{
					"error": "ไม่สามารถสร้างแท็ก: " + name,
				})
			}
		}
		tags = append(tags, tag)
	}

	// สร้างบทความใหม่
	article := models.Article{
		Title:      input.Title,
		Slug:       input.Slug,
		Content:    input.Content,
		AuthorID:   authorID,
		CategoryID: category.ID,
		Tags:       tags,
	}

	if err := database.DB.Create(&article).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "ไม่สามารถสร้างบทความ",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data":    article,
		"message": "✅ สร้างบทความสำเร็จ",
	})
}


// UpdateArticle updates an existing article
func UpdateArticle(c *fiber.Ctx) error {
	slug := c.Params("slug")

	// รับ user ID จาก JWT
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	// โหลดบทความพร้อม Tags
	var article models.Article
	if err := database.DB.Preload("Tags").First(&article, "slug = ?", slug).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("ไม่พบบทความ"))
	}

	// ตรวจสอบว่า user เป็นเจ้าของบทความหรือไม่
	if article.AuthorID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("คุณไม่มีสิทธิ์แก้ไขบทความนี้"))
	}

	// รับ input แบบ pointer เพื่อเช็กว่าผู้ใช้ส่งอะไรมาบ้าง
	var input struct {
		Title      *string  `json:"title"`
		Content    *string  `json:"content"`
		CategoryID *uint    `json:"category_id"`
		TagIDs     []uint   `json:"tag_ids"`
		NewTags    []string `json:"new_tags"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("ข้อมูลไม่ถูกต้อง"))
	}

	// อัปเดตเฉพาะฟิลด์ที่ถูกส่งมา
	if input.Title != nil {
		article.Title = *input.Title
	}
	if input.Content != nil {
		article.Content = *input.Content
	}
	if input.CategoryID != nil {
		var category models.Category
		if err := database.DB.First(&category, *input.CategoryID).Error; err != nil {
			return c.Status(400).JSON(utils.ErrorResponse("หมวดหมู่ไม่ถูกต้อง"))
		}
		article.CategoryID = *input.CategoryID
	}

	if err := database.DB.Save(&article).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("บันทึกบทความล้มเหลว"))
	}

	// ตรวจสอบและสร้าง new_tags หากมี
	var createdTags []models.Tags
	for _, name := range input.NewTags {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		var existing models.Tags
		if err := database.DB.Where("LOWER(name) = ?", strings.ToLower(name)).First(&existing).Error; err == nil {
			continue // แท็กซ้ำ
		} else if err != gorm.ErrRecordNotFound {
			return c.Status(500).JSON(utils.ErrorResponse("ตรวจสอบแท็กใหม่ล้มเหลว"))
		}

		newTag := models.Tags{Name: name}
		if err := database.DB.Create(&newTag).Error; err != nil {
			return c.Status(500).JSON(utils.ErrorResponse("สร้างแท็กใหม่ล้มเหลว"))
		}
		createdTags = append(createdTags, newTag)
	}

	// รวม tag_ids เก่ากับที่สร้างใหม่
	allTagIDs := input.TagIDs
	for _, tag := range createdTags {
		allTagIDs = append(allTagIDs, tag.ID)
	}

	// ถ้าส่ง tag_ids หรือ new_tags มาด้วย → จัดการแท็ก
	if input.TagIDs != nil || len(createdTags) > 0 {
		if len(allTagIDs) > 0 {
			var allTags []models.Tags
			if err := database.DB.Where("id IN ?", allTagIDs).Find(&allTags).Error; err != nil {
				return c.Status(500).JSON(utils.ErrorResponse("โหลดแท็กล้มเหลว"))
			}
			if err := database.DB.Model(&article).Association("Tags").Replace(&allTags); err != nil {
				return c.Status(500).JSON(utils.ErrorResponse("อัปเดตแท็กล้มเหลว"))
			}
		} else {
			// ไม่มีแท็กใด ๆ → ลบความสัมพันธ์ทั้งหมด
			if err := database.DB.Model(&article).Association("Tags").Clear(); err != nil {
				return c.Status(500).JSON(utils.ErrorResponse("ลบแท็กเก่าไม่สำเร็จ"))
			}
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

	// get user ID from JWT
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	// find article by slug with tags
	var article models.Article
	if err := database.DB.Preload("Tags").First(&article, "slug = ?", slug).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("ไม่พบบทความ"))
	}

	// check if user is the author
	if article.AuthorID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("คุณไม่มีสิทธิ์ลบบทความนี้"))
	}

	// delete tag relationship
	if err := database.DB.Model(&article).Association("Tags").Clear(); err != nil {
		log.Println("❌ ลบความสัมพันธ์กับ Tags ไม่สำเร็จ:", err)
		return c.Status(500).JSON(utils.ErrorResponse("ลบแท็กของบทความไม่สำเร็จ"))
	}

	// delete article
	if err := database.DB.Delete(&article).Error; err != nil {
		log.Println("🔥 Delete article failed:", err)
		return c.Status(500).JSON(utils.ErrorResponse("ลบบทความไม่สำเร็จ"))
	}

	return c.JSON(utils.SuccessResponse(nil, "ลบบทความเรียบร้อยแล้ว"))
}
