package routes

import (
	"log"
	"time"

	"hse-results/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	app.Static("/static/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	app.Get("/search", func(c *fiber.Ctx) error {
		return c.SendFile("./public/search.html")
	})

	app.Get("/search/indivitual", func(c *fiber.Ctx) error {
		name := c.Query("name")

		log.Println(name)
		res, _ := services.SearchStudentByName(name)
		return c.JSON(res)
	})

	app.Get("/hello", func(c *fiber.Ctx) error {
		time.Sleep(2 * time.Second)
		// return c.SendString("<li>Hello World!</li>")
		return c.SendFile("./public/index.html")
	})

	log.Fatal(app.Listen(":3000"))
}
