package utils

import "golang.org/x/crypto/bcrypt"

// hash password
func HashPassword(pw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), 12)
	return string(hash)
}

// check password
func CheckPassword(pw string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err == nil
}
