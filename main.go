package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/paedayz/GO_FIBER_BookList.git/driver"
	models "github.com/paedayz/GO_FIBER_BookList.git/model"

	"github.com/subosito/gotenv"
)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := fiber.New()
	db = driver.ConnectDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	app.Get("/books", func(c *fiber.Ctx) error {
		var book models.Book
		books := []models.Book{}

		rows, err := db.Query("select * from books")
		logFatal(err)

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year) // Check if every column of rows had data
			logFatal(err)

			books = append(books, book)
		}

		return c.JSON(books)
	})

	app.Post("/book", func(c *fiber.Ctx) error {
		var bookID int
		b := new(Book)
		if err := c.BodyParser(b); err != nil {
			return err
		}

		fmt.Println(b)
		err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;",
			b.Title, b.Author, b.Year).Scan(&bookID)

		logFatal(err)
		return c.SendString(strconv.Itoa(bookID))
	})

	app.Listen(":3000")
}
