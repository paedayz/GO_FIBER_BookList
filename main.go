package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/paedayz/GO_FIBER_BookList.git/driver"

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
