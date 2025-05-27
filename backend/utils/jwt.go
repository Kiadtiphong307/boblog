package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// คือฟังก์ชันที่จะสร้าง JWT token ให้กับผู้ใช้
func GenerateJWT(userID uint, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
