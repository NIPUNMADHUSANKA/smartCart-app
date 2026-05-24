package routes

import (
	"smartCart-app/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupUnProtectedRoutes(router *fiber.App) {

	router.Post("/api/smart-cart/auth/register", controllers.RegisterUser())
	router.Post("/api/smart-cart/auth/login", controllers.LoginUser())

}
