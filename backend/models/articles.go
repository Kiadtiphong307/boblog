package models

import "time"

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
