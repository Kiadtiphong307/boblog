package utils

import "golang.org/x/crypto/bcrypt"

// คือฟังก์ชันที่จะจัดการกับการลบรหัสผ่านของผู้ใช้
func HashPassword(pw string) string {
	// สร้างรหัสผ่านของผู้ใช้ใหม่
	// ใช้ฟังก์ชัน GenerateFromPassword จาก bcrypt เพื่อสร้างรหัสผ่านของผู้ใช้ใหม่
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), 12)
	return string(hash)
}

// คือฟังก์ชันที่จะจัดการกับการตรวจสอบรหัสผ่านของผู้ใช้
func CheckPassword(pw string, hash string) bool {
	// ตรวจสอบรหัสผ่านว่ามีความถูกต้องหรือไม่
	// ใช้ฟังก์ชัน CompareHashAndPassword จาก bcrypt เพื่อตรวจสอบรหัสผ่านของผู้ใช้
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err == nil
}
