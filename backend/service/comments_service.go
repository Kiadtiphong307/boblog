package service

import (
    "backend/database"
	"backend/models"
	"time"
	"github.com/gofiber/fiber/v2"
    "strings"
	"net/url"
	"backend/utils"
)
 

// Get Comments
func HandleGetComments(c *fiber.Ctx) error {
	//  Params & Query
    slugEncoded := c.Params("slug")
	// ทำงานยังไง
    slug, err := url.QueryUnescape(slugEncoded)
    if err != nil {
        return c.Status(400).JSON(utils.ErrorResponse("Invalid slug encoding"))
    }
	// วิธีการใช้ var , := 
    var article models.Article
	//  ORM 
    if err := database.DB.Where("slug = ?", slug).First(&article).Error; err != nil {
        return c.Status(404).JSON(utils.ErrorResponse("Article not found"))
    }
    var comments []models.Comment
    if err := database.DB.
        Preload("User").
        Order("created_at desc").
		// SQL 
        Where("article_id = ?", article.ID).
        Find(&comments).Error; err != nil {
		// Find & Frist 
        return c.Status(500).JSON(utils.ErrorResponse("Failed to fetch comments"))
    }
    return c.JSON(utils.SuccessResponse(comments, "get comments success"))
}

// Create Comment
func HandleCreateComment(c *fiber.Ctx) error {

	slugEncoded := c.Params("slug")
			
	slug, err := url.QueryUnescape(slugEncoded)
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid slug encoding"))
	}

	var article models.Article
	if err := database.DB.Where("slug = ?", slug).First(&article).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Article not found"))
	}

	var input struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid request body"))
	}

	if strings.TrimSpace(input.Content) == "" {
		return c.Status(400).JSON(utils.ErrorResponse("Content is required"))
	}

	userIDRaw := c.Locals("userID")
	userID, ok := userIDRaw.(uint)
	if !ok {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}


	comment := models.Comment{
		Content:   input.Content,
		ArticleID: article.ID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to create comment"))
	}

	return c.JSON(comment)
}

// Update Comment
func HandleUpdateComment(c *fiber.Ctx) error {
	commentID := c.Params("commentId")
	if commentID == "" {
		return c.Status(400).JSON(utils.ErrorResponse("Comment ID is required"))
	}

	userIDRaw := c.Locals("userID")
	userID, ok := userIDRaw.(uint)
	if !ok {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}

	var comment models.Comment
	if err := database.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Comment not found"))
	}

	if comment.UserID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("Forbidden"))
	}

	var input struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid request body"))
	}

	if strings.TrimSpace(input.Content) == "" {
		return c.Status(400).JSON(utils.ErrorResponse("Content is required"))
	}
	
	comment.Content = input.Content
	comment.UpdatedAt = time.Now()

	if err := database.DB.Save(&comment).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to update comment"))
	}

	
	//  โหลดข้อมูล User เพื่อ return กลับไป
	if err := database.DB.Preload("User").First(&comment, commentID).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to update comment"))
	}

	return c.JSON(comment)
}

// Delete Comment
func HandleDeleteComment(c *fiber.Ctx) error {
	commentID := c.Params("commentId")
	if commentID == "" {
		return c.Status(400).JSON(utils.ErrorResponse("Comment ID is required"))
	}

	userIDRaw := c.Locals("userID")
	userID, ok := userIDRaw.(uint)
	if !ok {
		return c.Status(401).JSON(utils.ErrorResponse("Unauthorized"))
	}

	var comment models.Comment
	if err := database.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Comment not found"))
	}

	if comment.UserID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("Forbidden"))
	}


	if err := database.DB.Delete(&comment).Error; err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to delete comment"))
	}

	return c.JSON(utils.SuccessResponse(nil, "delete comment success"))
}


