package seed

import (
	"blog-db/database"
	"blog-db/models"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func SeedCategories() {
	categories := []models.Category{
		{Name: "Technology"},
		{Name: "Health"},
		{Name: "Business"},
		{Name: "Education"},
		{Name: "Entertainment"},
	}

	for _, category := range categories {
		var existing models.Category
		database.DB.Where("name = ?", category.Name).First(&existing)
		if existing.ID == 0 {
			database.DB.Create(&category)
		}
	}
	log.Println("✅ Seeded categories")
}

func SeedTags() {
	tags := []models.Tags{
		{Name: "Go"},
		{Name: "Docker"},
		{Name: "API"},
		{Name: "Kubernetes"},
	}

	for _, tag := range tags {
		var existing models.Tags
		database.DB.Where("name = ?", tag.Name).First(&existing)
		if existing.ID == 0 {
			database.DB.Create(&tag)
		}
	}
	log.Println("✅ Seeded tags")
}

func SeedUserAndArticles() {
	// Create user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 14)
	user := models.User{
		Username:     "admin",
		FirstName:    "Admin",
		LastName:     "User",
		Nickname:     "Boss",
		Email:        "admin@example.com",
		PasswordHash: string(hashedPassword),
	}
	database.DB.Where("email = ?", user.Email).FirstOrCreate(&user)

	// Get category
	var techCat models.Category
	if err := database.DB.First(&techCat, "name = ?", "Technology").Error; err != nil {
		log.Println("❌ Category not found, skipping article seeding")
		return
	}

	// Get tags
	var tags []models.Tags
	if err := database.DB.Find(&tags).Error; err != nil || len(tags) < 4 {
		log.Println("❌ Not enough tags found")
		return
	}

	// Article seed data
	articleData := []struct {
		Title   string
		Slug    string
		Content string
		TagSet  []models.Tags
	}{
		{"Go Fiber Framework", "go-fiber-framework", strings.Repeat("Go Fiber is blazing fast. ", 10), []models.Tags{tags[0]}},
		{"Understanding Goroutines", "understanding-goroutines", strings.Repeat("Goroutines are lightweight threads. ", 10), []models.Tags{tags[0]}},
		{"Docker for Developers", "docker-for-developers", strings.Repeat("Docker simplifies deployment. ", 10), []models.Tags{tags[1]}},
		{"RESTful API Design", "restful-api-design", strings.Repeat("Designing APIs with REST principles. ", 10), []models.Tags{tags[2]}},
		{"Intro to Kubernetes", "intro-to-kubernetes", strings.Repeat("Kubernetes manages containers. ", 10), []models.Tags{tags[3]}},
	}

	for _, item := range articleData {
		var existing models.Article
		database.DB.Where("slug = ?", item.Slug).First(&existing)
		if existing.ID == 0 {
			article := models.Article{
				Title:      item.Title,
				Slug:       item.Slug,
				Content:    item.Content,
				AuthorID:   user.ID,
				CategoryID: techCat.ID,
			}
			if err := database.DB.Create(&article).Error; err != nil {
				log.Println("❌ Failed to create article:", article.Slug)
				continue
			}

			// ✅ ผูกแท็กหลายตัว
			if err := database.DB.Model(&article).Association("Tags").Replace(&item.TagSet); err != nil {
				log.Println("❌ Failed to associate tags for", article.Slug)
			}
		}
	}

	log.Println("✅ Seeded user and articles with tags")
}

