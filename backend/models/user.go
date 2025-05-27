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
	Bio          *string   `json:"bio,omitempty"`
	Image        *string   `json:"image,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Articles     []Article `gorm:"foreignKey:AuthorID" json:"-"`
	Comments     []Comment `gorm:"foreignKey:UserID" json:"-"`
}
