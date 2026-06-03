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

	/*
		11

			export const GENERATE_AI_PROMPT = `${host}/api/smart-cart/ai-model`;
			export const GET_ALL_AI_CATEGORY = `${host}/api/smart-cart/ai-model`;
			export const DELETE_AI_CATEGORY = `${host}/api/smart-cart/ai-model/:categoryId`;
			export const DELETE_AI_ALL_CATEGORY = `${host}/api/smart-cart/ai-model/deleteAISuggestion/:suggestionId`;
			export const DELETE_AI_SHOPPING_ITEM = `${host}/api/smart-cart/ai-model/deleteAIShoppingItem/:categoryId/:itemId`;
			export const UPDATE_AI_SHOPPING_ITEM = `${host}/api/smart-cart/ai-model/updateAIShoppingItem`;
			export const REGENERATE_AI_PROMPT = `${host}/api/smart-cart/ai-model/regenerateAIShopping`;
			export const CONFIRM_AI_CATEGORY = `${host}/api/smart-cart/ai-model/confirmAIShopping`;
			export const ADD_AI_SHOPPING_ITEM = `${host}/api/smart-cart/ai-model/addAIShoppingItem`;

			export const REMOVE_USER = `${host}/api/smart-cart/auth/remove`;

	*/
}
