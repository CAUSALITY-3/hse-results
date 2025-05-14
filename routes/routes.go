package routes

import (
	"log"
	"os"
	"time"

	"hse-results/services"
	"hse-results/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	if os.Getenv("ENABLE_RESOURCE_LIMIT_FOR_RATE_LIMIT") == "true" {
		go utils.MnitorSystem()
	}
	go utils.CeanupVisitors()

	app.Use(utils.RateLimiterMiddleware)
	// app.Use(compress.New())

	// app.Static("/static/", "./public")

	app.Get("/static/*", utils.ServeCompressedFile("./public"))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	app.Get("/chart", func(c *fiber.Ctx) error {
		return c.SendFile("./public/chart.html")
	})

	app.Get("/:resultType/search", func(c *fiber.Ctx) error {
		resultType := c.Params("resultType") // Get dynamic class name
		searchType := c.Query("searchType")  // Get dynamic class name
		log.Println(resultType, searchType)
		if searchType == "school" {
			return c.SendFile("./public/pages/schoolSearch.html")
		}
		return c.SendFile("./public/pages/search.html")
	})

	app.Get("/search/indivitual", func(c *fiber.Ctx) error {
		name := c.Query("name")

		log.Println(name)
		res, _ := services.SearchStudentByName(name)
		return c.JSON(res)
	})

	app.Get("/:resultType/search/student/:rollno", compress.New(), func(c *fiber.Ctx) error {
		resultType := c.Params("resultType")
		rollNo := c.Params("rollno")
		log.Println(resultType, rollNo)

		err := services.GetStudentResults(c, resultType, "23274662")
		if err != nil {
			log.Println("Error executing template:", err)
			return c.Status(500).SendString("Error rendering template")
		}

		return nil

	})

	app.Get("/:resultType/search/school/:schoolCode", func(c *fiber.Ctx) error {
		resultType := c.Params("resultType")
		schoolCode := c.Params("schoolCode")
		log.Println(resultType, schoolCode)

		err := services.GetSchoolResultsServerSide(c, resultType, schoolCode)
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

	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
				c.Status(500).SendString("Internal Server Error")
			}
		}()
		return c.Next()
	})

	log.Fatal(app.Listen(":3000"))
}
