package routes

import (
	"smartCart-app/controllers"
	"smartCart-app/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUnProtectedRoutes(router *fiber.App) {
	router.Use(middleware.RateLimiting())

	api := router.Group("/api/smart-cart/")

	api.Post("auth/register", controllers.RegisterUser())
	api.Post("auth/login", controllers.LoginUser())

}
