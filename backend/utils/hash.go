package utils

import "golang.org/x/crypto/bcrypt"

// คือฟังก์ชันที่จะจัดการกับการลบรหัสผ่านของผู้ใช้
func HashPassword(pw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), 12)
	return string(hash)
}

// คือฟังก์ชันที่จะจัดการกับการตรวจสอบรหัสผ่านของผู้ใช้
func CheckPassword(pw string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err == nil
}
