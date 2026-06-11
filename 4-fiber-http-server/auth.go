package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// dummy user data
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var memberUser = User{
	Email:    "user@example.com",
	Password: "password123",
}

func loginHandler(c fiber.Ctx) error {
	// สร้างกล่องเปล่าๆ ที่มีช่อง Email กับ Password มารอรับข้อมูล ตามโครงสร้างของ struct User
	user := new(User)
	// อ่าน JSON ดิบจาก Client แล้วเอาค่าที่ชื่อตรงกันมาเติมใส่ในกล่องเปล่า 'user' อัตโนมัติ
	// โดยใช้ c.Bind().Body() ซึ่งจะทำการแปลง JSON เป็น struct User ให้เรา
	if err := c.Bind().Body(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user.Email == memberUser.Email && user.Password == memberUser.Password {
		// Create Token
		token := jwt.New(jwt.SigningMethodHS256)

		// set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["email"] = user.Email
		claims["role"] = "admin"
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// generate encoded token and send it as response
		t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.JSON(fiber.Map{
			"message": "Login Success",
			"token":   t,
		})

	}
	return fiber.ErrUnauthorized

}
