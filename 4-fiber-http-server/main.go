package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func main() {
	// init engine
	engine := html.New("./views", ".html")

	// pass engine to fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Setup route
	app.Get("/", renderTemplate)

	app.Get("/hello", func(c fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	// init book data
	books = append(books, Book{ID: 1, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"})
	books = append(books, Book{ID: 2, Title: "To Kill a Mockingbird", Author: "Harper Lee"})
	books = append(books, Book{ID: 3, Title: "1984", Author: "George Orwell"})

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	app.Post("/upload", uploadFile)

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
