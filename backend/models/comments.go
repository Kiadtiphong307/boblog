package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	ArticleID uint      `gorm:"not null;index" json:"-"`
	Article   Article   `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	UserID    uint      `gorm:"not null" json:"-"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	CreatedAt time.Time `json:"created_at"`
}