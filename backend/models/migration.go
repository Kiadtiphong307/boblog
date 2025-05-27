package models

import "time"

// User Model
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique" json:"username"`
	FirstName    string    `gorm:"not null" json:"first_name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	Nickname     string    `gorm:"not null" json:"nickname"`
	Email        string    `gorm:"unique" json:"email"`
	PasswordHash string    `gorm:"not null" json:"password"`
	Bio          *string   `json:"bio,omitempty"`
	Image        *string   `json:"image,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Articles     []Article `gorm:"foreignKey:AuthorID" json:"-"`
	Comments     []Comment `gorm:"foreignKey:UserID" json:"-"`
}

// Article Model
type Article struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"title"`
	Slug       string    `gorm:"unique;not null" json:"slug"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	AuthorID   uint      `gorm:"not null" json:"-"`
	Author     User      `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"author"`
	CategoryID uint      `gorm:"not null" json:"-"`
	Category   Category  `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"category"`
	Tags       []Tags    `gorm:"many2many:article_tags;" json:"tags"`
	Comments   []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Comment Model
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	ArticleID uint      `gorm:"not null;index" json:"-"`
	Article   Article   `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	UserID    uint      `gorm:"not null" json:"-"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

// Category Model
type Category struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"unique;not null" json:"name"`
	Articles []Article `gorm:"foreignKey:CategoryID" json:"-"`
}

// Tags Model
type Tags struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}
