package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique" json:"username"`
	FirstName    string    `gorm:"not null" json:"first_name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	Nickname     string    `gorm:"not null" json:"nickname"`
	Email        string    `gorm:"unique" json:"email"`
	PasswordHash string    `gorm:"not null" json:"password"`
	Bio          *string   `gorm:"default:null" json:"bio,omitempty"`
	Image        *string   `gorm:"default:null" json:"image,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Article struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"not null" json:"title"`
	Slug       string    `gorm:"unique;not null" json:"slug"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	AuthorID   uint      `gorm:"not null" json:"-"`
	Author     User      `gorm:"foreignKey:AuthorID" json:"author"`
	CategoryID uint      `gorm:"not null" json:"-"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category"`
	Tags       []Tags     `gorm:"many2many:article_tags;" json:"tags"` 
	Comments   []Comment `gorm:"foreignKey:ArticleID" json:"comments"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}


type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	ArticleID uint      `gorm:"not null" json:"-"`
	UserID    uint      `gorm:"not null" json:"-"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Article   Article   `gorm:"foreignKey:ArticleID" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"unique;not null" json:"name"`
	Articles []Article `json:"-"`
}

type Tags struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}

