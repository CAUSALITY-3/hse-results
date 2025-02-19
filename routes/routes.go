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

	app.Get("/:resultType/search", func(c *fiber.Ctx) error {
		resultType := c.Params("resultType") // Get dynamic class name
		return c.SendFile("./public/" + resultType + "/search.html")
	})

	app.Get("/search/indivitual", func(c *fiber.Ctx) error {
		name := c.Query("name")

		log.Println(name)
		res, _ := services.SearchStudentByName(name)
		return c.JSON(res)
	})

	app.Get("/:resultType/search/:rollno", func(c *fiber.Ctx) error {
		resultType := c.Params("resultType")
		rollNo := c.Params("rollno")
		log.Println(resultType, rollNo)

		err := services.GetStudentResults(c, resultType, "23007265")
		if err != nil {
			log.Println("Error executing template:", err)
			return c.Status(500).SendString("Error rendering template")
		}

		return nil

	})

	app.Get("/hello", func(c *fiber.Ctx) error {
		time.Sleep(0 * time.Second)
		searchText := c.Query("nameOrRollNo")

		log.Println("query params : ", searchText)
		// return c.SendString("<li>Hello World!</li>")
		return c.SendFile("./public/index.html")
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Redirect("/", 302) // 302 Found (Temporary Redirect)
	})

	log.Fatal(app.Listen(":3000"))
}
