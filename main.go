package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/paedayz/GO_FIBER_BookList.git/driver"

	"github.com/subosito/gotenv"
)

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

	app.Listen(":3000")
}
