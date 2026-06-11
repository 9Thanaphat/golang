package main

import (
	"fmt"
	"os"

	_ "fiber-http-server/docs"

	swaggo "github.com/gofiber/contrib/v3/swaggo"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func isAdmin(c fiber.Ctx) error {
	// ดึง Token ออกมาจาก Context ด้วยฟังก์ชันเฉพาะของ v3
	user := jwtware.FromContext(c)

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Token not found in context")
	}

	claims := user.Claims.(jwt.MapClaims)
	fmt.Println("ข้อมูลใน Token คือ:", claims)

	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	// init engine
	engine := html.New("./views", ".html")

	// pass engine to fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Setup Swagger
	app.Get("/docs/*", swaggo.HandlerDefault)

	// Setup route
	app.Get("/", renderTemplate)

	app.Get("/hello", func(c fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	// init book data
	books = append(books, Book{ID: 1, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"})
	books = append(books, Book{ID: 2, Title: "To Kill a Mockingbird", Author: "Harper Lee"})
	books = append(books, Book{ID: 3, Title: "1984", Author: "George Orwell"})

	app.Post("/login", loginHandler)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(os.Getenv("JWT_SECRET")),
		},
		Extractor: extractors.FromAuthHeader("Bearer"),
	}))

	bookGroup := app.Group("/book")

	bookGroup.Use(isAdmin)

	bookGroup.Get("/books", getBooks)
	bookGroup.Get("/books/:id", getBook)
	bookGroup.Post("/books", createBook)
	bookGroup.Put("/books/:id", updateBook)
	bookGroup.Delete("/books/:id", deleteBook)

	app.Post("/upload", uploadFile)

	app.Get("config", getConfig)
	app.Get("config2", getConfig2)

	app.Listen(":8080")
	fmt.Println("Server is running on http://localhost:8080")

}

func renderTemplate(c fiber.Ctx) error {
	// Render the template with variable data
	return c.Render("template", fiber.Map{
		"Name": "9Thanaphat",
	})
}

func uploadFile(c fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	err = c.SaveFile(file, "./uploads/"+file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusCreated).SendString("File uploaded successfully")
}

func getConfig(c fiber.Ctx) error {
	if value, exists := os.LookupEnv("SECRET"); exists {
		return c.JSON(fiber.Map{
			"secret": value,
		})
	}
	return c.Status(fiber.StatusNotFound).SendString("Config not found")
}

// get SECRET2 from .env file and return as JSON response
func getConfig2(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"secret": os.Getenv("SECRET2"),
	})
}
