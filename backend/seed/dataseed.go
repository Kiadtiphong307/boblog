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

func SeedUserAndArticles() {
	// Seed user
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

	// Get category to use
	var techCat models.Category
	if err := database.DB.First(&techCat, "name = ?", "Technology").Error; err != nil {
		log.Println("❌ Category not found, skipping article seeding")
		return
	}

	// Seed articles
	articles := []models.Article{
		{Title: "Go Fiber Framework", Slug: "go-fiber-framework", Content: strings.Repeat("Go Fiber is blazing fast. ", 10), AuthorID: user.ID, CategoryID: techCat.ID},
		{Title: "Understanding Goroutines", Slug: "understanding-goroutines", Content: strings.Repeat("Goroutines are lightweight threads. ", 10), AuthorID: user.ID, CategoryID: techCat.ID},
		{Title: "Docker for Developers", Slug: "docker-for-developers", Content: strings.Repeat("Docker simplifies deployment. ", 10), AuthorID: user.ID, CategoryID: techCat.ID},
		{Title: "RESTful API Design", Slug: "restful-api-design", Content: strings.Repeat("Designing APIs with REST principles. ", 10), AuthorID: user.ID, CategoryID: techCat.ID},
		{Title: "Intro to Kubernetes", Slug: "intro-to-kubernetes", Content: strings.Repeat("Kubernetes manages containers. ", 10), AuthorID: user.ID, CategoryID: techCat.ID},
	}

	for _, article := range articles {
		var existing models.Article
		database.DB.Where("slug = ?", article.Slug).First(&existing)
		if existing.ID == 0 {
			database.DB.Create(&article)
		}
	}

	log.Println("✅ Seeded user and articles")
}
