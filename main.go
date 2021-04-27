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

	app.Get("/book/:id", func(c *fiber.Ctx) error {
		var book models.Book
		params_id := c.Params("id")

		rows := db.QueryRow("select * from books where id=$1", params_id) // Query data from database
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year) // Store data from database to var book
		logFatal(err)
		return c.JSON(book)
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

	app.Put("/book", func(c *fiber.Ctx) error {
		b := new(models.Book)
		if err := c.BodyParser(b); err != nil {
			return err
		}
		_, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id",
			b.Title, b.Author, b.Year, b.ID)
		logFatal(err)

		return c.JSON(b)
	})

	app.Delete("/book/:id", func(c *fiber.Ctx) error {
		params_id := c.Params("id")
		_, err := db.Exec("delete from books where id = $1", params_id)
		logFatal(err)

		return c.SendString(params_id)
	})

	app.Listen(":3000")
}
