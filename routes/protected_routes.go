package routes

import (
	"smartCart-app/controllers"
	"smartCart-app/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupProtectedRoutes(router *fiber.App) {
	router.Use(middleware.RateLimiting())
	router.Use(middleware.AuthMiddlware())

	api := router.Group("/api/smart-cart/")

	api.Get("category", controllers.GetCategory())
	api.Post("category", controllers.CreateCategory())
	api.Get("category/:categoryId", controllers.GetCategoryByCategoryId())
	api.Delete("category/:categoryId", controllers.DeleteCategoryByCategoryId())
	api.Patch("category/:categoryId", controllers.UpdateCategory())

	api.Get("auth/me", controllers.GetUserInfo())
	api.Get("auth/info", controllers.GetUserDetails())
	api.Patch("auth/resetPassword", controllers.ResetPassword())
	api.Delete("auth/remove", controllers.DeleteUser())

	api.Get("shopping-item", controllers.GetShoppingItems())
	api.Get("shopping-item/:itemId", controllers.GetShoppingItemByItemId())
	api.Get("shopping-item/findByCategory/:categoryId", controllers.GetShoppingItemByCategoryId())
	api.Delete("shopping-item/:itemId", controllers.DeleteShoppingItemByItemId())
	api.Post("shopping-item", controllers.CreateShoppingItems())
	api.Patch("shopping-item/:itemId", controllers.UpdateShoppingItem())

	api.Get("ai-model", controllers.GetAllAICategory())
	api.Delete("ai-model/:categoryId", controllers.DeleteAICategory())
	api.Delete("ai-model/deleteAISuggestion/:suggestionId", controllers.DeleteAISuggestion())
	api.Delete("ai-model/deleteAIShoppingItem/:categoryId/:itemId", controllers.DeleteAIShoppingItem())
	api.Patch("ai-model/updateAIShoppingItem", controllers.UpdateAIShoppingItem())
	api.Post("ai-model/addAIShoppingItem", controllers.AddAIShoppingItem())
	api.Post("ai-model/confirmAIShopping", controllers.ConfirmAICategory())
	api.Post("ai-model", controllers.GenetateAIPrompt())
	api.Post("ai-model/regenerateAIShopping", controllers.ReGenetateAIPrompt())
}
