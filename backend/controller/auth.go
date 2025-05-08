package controller

import (
    "blog-db/database"
    "blog-db/models"
    "blog-db/utils"
    "github.com/gofiber/fiber/v2"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "time"
    "log"
)

var jwtSecret = []byte("secret") // 👉 เปลี่ยนเป็น os.Getenv("JWT_SECRET") ภายหลัง

// Register
func Register(c *fiber.Ctx) error {
    var input models.User
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(utils.ErrorResponse("Invalid input"))
    }

    // เช็ค email ซ้ำ
    var existing models.User
    database.DB.Where("email = ?", input.Email).First(&existing)
    if existing.ID != 0 {
        return c.Status(400).JSON(utils.ErrorResponse("Email already in use"))
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), 14)
    input.PasswordHash = string(hashedPassword)

    database.DB.Create(&input)
    input.PasswordHash = ""

    return c.Status(201).JSON(utils.SuccessResponse(input, "User registered"))
}

func Login(c *fiber.Ctx) error {
    var input models.User
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(utils.ErrorResponse("Invalid input"))
    }

    var user models.User
    database.DB.Where("email = ?", input.Email).First(&user)
    if user.ID == 0 {
        return c.Status(400).JSON(utils.ErrorResponse("Invalid email or password"))
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.PasswordHash)); err != nil {
        return c.Status(400).JSON(utils.ErrorResponse("Invalid email or password"))
    }

    // ✅ Create JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id":    user.ID,
        "email": user.Email,
        "exp":   time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        log.Println("JWT error:", err)
        return c.SendStatus(500)
    }

    return c.JSON(utils.SuccessResponse(fiber.Map{
        "token": tokenString,
    }, "Login successful"))
}
