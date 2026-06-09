package routes

import (
	"smartCart-app/controllers"
	"smartCart-app/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupProtectedRoutes(router *fiber.App) {
	router.Use(middleware.AuthMiddlware())

	router.Get("/api/smart-cart/category", controllers.GetCategory())
	router.Post("/api/smart-cart/category", controllers.CreateCategory())
	router.Get("/api/smart-cart/category/:categoryId", controllers.GetCategoryByCategoryId())
	router.Delete("/api/smart-cart/category/:categoryId", controllers.DeleteCategoryByCategoryId())
	router.Patch("/api/smart-cart/category/:categoryId", controllers.UpdateCategory())

	router.Get("/api/smart-cart/auth/me", controllers.GetUserInfo())
	router.Get("/api/smart-cart/auth/info", controllers.GetUserDetails())
	router.Post("/api/smart-cart/auth/resetPassword", controllers.ResetPassword())

	router.Get("/api/smart-cart/shopping-item", controllers.GetShoppingItems())
	router.Get("/api/smart-cart/shopping-item/:itemId", controllers.GetShoppingItemByItemId())
	router.Get("/api/smart-cart/shopping-item/:categoryId", controllers.GetShoppingItemByCategoryId())
	router.Delete("/api/smart-cart/shopping-item/:itemId", controllers.DeleteShoppingItemByItemId())
	router.Post("/api/smart-cart/shopping-item", controllers.CreateShoppingItems())
	router.Patch("/api/smart-cart/shopping-item/:itemId", controllers.UpdateShoppingItem())

	router.Get("/api/smart-cart/ai-model", controllers.GetAllAICategory())
	router.Delete("/api/smart-cart/ai-model/:categoryId", controllers.DeleteAICategory())
	router.Delete("/api/smart-cart/ai-model/deleteAISuggestion/:suggestionId", controllers.DeleteAISuggestion())
	router.Delete("/api/smart-cart/ai-model/deleteAIShoppingItem/:categoryId/:itemId", controllers.DeleteAIShoppingItem())
	router.Post("/api/smart-cart/ai-model/updateAIShoppingItem", controllers.UpdateAIShoppingItem())
	router.Post("/api/smart-cart/ai-model/confirmAIShopping", controllers.ConfirmAICategory())
	router.Post("/api/smart-cart/ai-model/addAIShoppingItem", controllers.AddAIShoppingItem())
	router.Post("/api/smart-cart/ai-model", controllers.GenetateAIPrompt())
	router.Post("/api/smart-cart/ai-model/regenerateAIShopping/:suggestionId", controllers.ReGenetateAIPrompt())

	/*
		1
			export const REMOVE_USER = `${host}/api/smart-cart/auth/remove`;

	*/
}
