package main

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func main() {
	app := fiber.New()

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

	app.Listen(":8080")
	fmt.Println("Server is running on http://localhost:8080")

}

func getBooks(c fiber.Ctx) error {
	return c.JSON(books)
}

func getBook(c fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for _, book := range books {
		if book.ID == bookId {
			return c.JSON(book)
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}

func createBook(c fiber.Ctx) error {
	book := new(Book)
	if err := c.Bind().Body(book); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	books = append(books, *book)
	return c.JSON(book)
}

func updateBook(c fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	bookUpdate := new(Book)
	if err := c.Bind().Body(bookUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, book := range books {
		if book.ID == bookId {
			books[i].Title = bookUpdate.Title
			books[i].Author = bookUpdate.Author
			return c.JSON(books[i])
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}

func deleteBook(c fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, book := range books {
		if book.ID == bookId {
			// สมมติว่าเจอหนังสือที่ตรงกับ bookId ในรอบที่ i พอดี
			// เราจะทำการหั่น (Slice) อาเรย์เดิมออกเป็น 2 ท่อน:

			// ท่อนที่ 1: books[:i]
			// คือการดึงข้อมูลตั้งแต่ตัวแรกสุด มาจนถึง "ก่อน" ตำแหน่ง i

			// ท่อนที่ 2: books[i+1:]
			// คือการดึงข้อมูลตั้งแต่ "หลัง" ตำแหน่ง i (ก็คือ i+1) ไปจนถึงตัวสุดท้าย

			// เครื่องหมาย ... (Spread/Unpack operator)
			// ทำหน้าที่แตกข้อมูลในท่อนที่ 2 ออกมาทีละชิ้น เพื่อให้ฟังก์ชัน append() รับเข้าไปได้

			// สุดท้ายจับท่อน 1 ต่อด้วยท่อน 2 แล้วเอาไปทับตัวแปร books ตัวเดิม
			books = append(books[:i], books[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}
