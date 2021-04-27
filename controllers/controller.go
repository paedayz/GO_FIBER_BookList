package controllers

import "github.com/gofiber/fiber/v2"

type Controller struct{}

func (controller Controller) GetBooks(c *fiber.Ctx) error {
	return c.SendString("get books")
}
