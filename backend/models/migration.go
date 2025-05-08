package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique" json:"username"`
	Email        string    `gorm:"unique" json:"email"`
	PasswordHash string    `gorm:"not null" json:"password"`
	Bio          *string   `gorm:"default:null" json:"bio,omitempty"`
	Image        *string   `gorm:"default:null" json:"image,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Article struct {
	ID         uint      `gorm:"primaryKey"`
	Title      string    `gorm:"not null"`
	Slug       string    `gorm:"unique;not null"`
	Content    string    `gorm:"type:text;not null"`
	AuthorID   uint      `gorm:"not null"`
	Author     User      `gorm:"foreignKey:AuthorID"`
	CategoryID uint      `gorm:"not null"`
	Category   Category  `gorm:"foreignKey:CategoryID"`
	Comments   []Comment `gorm:"foreignKey:ArticleID"`
	Tags       []Tag     `gorm:"many2many:article_tags"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	ArticleID uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Article   Article   `gorm:"foreignKey:ArticleID"`
	CreatedAt time.Time `json:"created_at"`

}

type Category struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"unique;not null"`
	Articles []Article
}

type Tag struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"unique;not null"`
	Articles []Article `gorm:"many2many:article_tags"`
}
