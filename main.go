package main

import (
	"log"
	"smartCart-app/database"

	"github.com/gofiber/fiber/v2"
)

func main() {

	router := fiber.New()

	database.Connection()

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	if err := router.Listen(":8080"); err != nil {
		log.Fatal("Failed to Start Server", err)
	}
}
