package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	app.Get("/hello", func(c *fiber.Ctx) error {
		time.Sleep(2 * time.Second)
		return c.SendString("<li>Hello World!</li>")
	})

	log.Fatal(app.Listen(":3000"))
}
