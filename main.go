package main

import (
	"log"
	"smartCart-app/database"
	"smartCart-app/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	router := fiber.New()

	/*Catches unexpected runtime panics (crashes) anywhere in your HTTP request lifecycle and sends a clean 500 Internal Server Error response back to the user instead of letting the entire server crash for everyone else.*/
	router.Use(recover.New())

	router.Use(helmet.New())

	router.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4200, https://smart-cart-web-eight.vercel.app",
		AllowMethods: "GET,HEAD,PUT,PATCH,POST,DELETE",
	}))
	database.Connection()

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	routes.SetupUnProtectedRoutes(router)
	routes.SetupProtectedRoutes(router)

	router.All("*", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "404 or fallback route",
		})
	})

	if err := router.Listen(":8080"); err != nil {
		log.Fatal("Failed to Start Server", err)
	}
}
