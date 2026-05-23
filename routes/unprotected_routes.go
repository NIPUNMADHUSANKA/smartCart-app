package routes

import (
	"smartCart-app/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupUnProtectedRoutes(router *fiber.App) {
	router.Post("/signup", controllers.RegisterUser())
}
