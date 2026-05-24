package middleware

import (
	"smartCart-app/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddlware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := utils.GetAccessToken(ctx)

		if err != nil || token == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		claims, err := utils.ValidateToken(token)

		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		ctx.Locals("userName", claims.UserName)
		ctx.Locals("userId", claims.UserId.String())
		ctx.Locals("role", claims.Role)

		ctx.Next()

		return nil
	}
}
