package routes

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/main.html")
	})

	app.Get("/hello", func(c *fiber.Ctx) error {
		time.Sleep(2 * time.Second)
		// return c.SendString("<li>Hello World!</li>")
		return c.SendFile("./public/main.html")
	})

	log.Fatal(app.Listen(":3000"))
}
