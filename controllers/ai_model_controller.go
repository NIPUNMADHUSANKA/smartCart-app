package controllers

import (
	"context"
	"smartCart-app/database"
	"smartCart-app/models"
	"smartCart-app/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetAllAICategory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "userId is not found",
			})
		}

		// Fetch AI suggestions for the user, then their categories and items
		sugRows, err := database.DBPool.Query(
			ctx,
			`SELECT s."id", s."userId", s."prompt"
			FROM "AiSuggestion" s
			WHERE s."userId" = $1`,
			userId,
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch AI suggestions"})
		}
		defer sugRows.Close()

		var suggestions []models.AiSuggestion

		for sugRows.Next() {
			var s models.AiSuggestion
			if err := sugRows.Scan(&s.Id, &s.UserId, &s.Prompt); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan AI suggestion"})
			}

			// fetch categories for this suggestion
			catRows, err := database.DBPool.Query(
				ctx,
				`SELECT c."id", c."suggestionId", c."categoryName", c."priority"
				FROM "AICategory" c
				WHERE c."suggestionId" = $1`,
				s.Id,
			)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch AI categories"})
			}

			var categories []models.AICategory

			for catRows.Next() {
				var cat models.AICategory
				var catId uuid.UUID
				var suggestionId uuid.UUID

				if err := catRows.Scan(&catId, &suggestionId, &cat.CategoryName, &cat.Priority); err != nil {
					catRows.Close()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan AI category"})
				}

				cat.Id = catId
				cat.SuggestionId = suggestionId.String()

				// fetch items for this category
				itemRows, err := database.DBPool.Query(
					ctx,
					`SELECT i."id", i."categoryId", i."itemName", i."quantity", i."unit", i."priority"
					FROM "AIItem" i
					WHERE i."categoryId" = $1`,
					cat.Id,
				)

				if err != nil {
					catRows.Close()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch AI items"})
				}

				var items []models.AIItem

				for itemRows.Next() {
					var it models.AIItem
					var itemId uuid.UUID
					var categoryId uuid.UUID
					if err := itemRows.Scan(&itemId, &categoryId, &it.ItemName, &it.Quantity, &it.Unit, &it.Priority); err != nil {
						itemRows.Close()
						catRows.Close()
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan AI item"})
					}
					it.Id = itemId
					it.CategoryId = categoryId.String()
					items = append(items, it)
				}

				itemRows.Close()

				if err := itemRows.Err(); err != nil {
					catRows.Close()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error reading AI item rows"})
				}

				cat.Items = items
				categories = append(categories, cat)
			}

			catRows.Close()
			if err := catRows.Err(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error reading AI category rows"})
			}

			s.Categories = categories
			suggestions = append(suggestions, s)
		}

		if err := sugRows.Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error reading AI suggestion rows"})
		}

		return c.Status(fiber.StatusOK).JSON(suggestions)

	}
}

func DeleteAICategory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		categoryId := c.Params("categoryId")

		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "userId is not found"})
		}

		if categoryId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "categoryId is not found"})
		}

		tx, err := database.DBPool.Begin(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
		}
		defer tx.Rollback(ctx)

		_, err = tx.Exec(ctx,
			`DELETE FROM "AIItem"
			WHERE "categoryId" = $1`,
			categoryId,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete AI items"})
		}

		res, err := tx.Exec(ctx,
			`DELETE FROM "AICategory"
			USING "AiSuggestion"
			WHERE "AICategory"."id" = $1
			AND "AICategory"."suggestionId" = "AiSuggestion"."id"
			AND "AiSuggestion"."userId" = $2`,
			categoryId,
			userId,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete AI category"})
		}

		rowsAffected := res.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "AI category not found"})
		}

		if err := tx.Commit(ctx); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit deletion"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "AI category deleted successfully",
		})
	}
}

func DeleteAIShoppingItem() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		categoryId := c.Params("categoryId")
		itemId := c.Params("itemId")

		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "userId is not found"})
		}

		if categoryId == "" || itemId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "categoryId or itemId is not found"})
		}

		var deletedItem models.AIItem

		err := database.DBPool.QueryRow(
			ctx,
			`DELETE FROM "AIItem" i
			USING "AICategory" c, "AiSuggestion" s
			WHERE i."id" = $1
			AND i."categoryId" = $2
			AND i."categoryId" = c."id"
			AND c."suggestionId" = s."id"
			AND s."userId" = $3
			RETURNING i."id", i."categoryId", i."itemName", i."quantity", i."unit", i."priority"`,
			itemId,
			categoryId,
			userId,
		).Scan(
			&deletedItem.Id,
			&deletedItem.CategoryId,
			&deletedItem.ItemName,
			&deletedItem.Quantity,
			&deletedItem.Unit,
			&deletedItem.Priority,
		)

		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "AI item not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete AI item"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "AI item deleted successfully",
		})
	}
}

func AddAIShoppingItem() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var AIItem models.AIItem

		if err := c.BodyParser(&AIItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input data",
			})
		}

		if err := validate.Struct(AIItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation Failed",
				"details": err.Error(),
			})
		}

		var itemId uuid.UUID

		err := database.DBPool.QueryRow(
			ctx,
			`INSERT INTO "AIItem" ("Id", "CategoryId", "ItemName", "quantity", "unit", "priority")
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING "id"`,
			AIItem.Id, AIItem.CategoryId, AIItem.ItemName, AIItem.Quantity, AIItem.Unit, AIItem.Priority,
		).Scan(&itemId)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create category",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(AIItem)
	}
}

func UpdateAIShoppingItem() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "userId is not found"})
		}

		var payload models.AIItem

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input data"})
		}

		if err := validate.Struct(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation Failed", "details": err.Error()})
		}

		var updatedItem models.AIItem
		err := database.DBPool.QueryRow(
			ctx,
			`UPDATE "AIItem" i
			SET "itemName" = $1,
				"quantity" = $2,
				"unit" = $3,
				"priority" = $4
			FROM "AICategory" c
			JOIN "AiSuggestion" s ON c."suggestionId" = s."id"
			WHERE i."id" = $5
			AND i."categoryId" = $6
			AND i."categoryId" = c."id"
			AND s."userId" = $7
			RETURNING i."id", i."categoryId", i."itemName", i."quantity", i."unit", i."priority"`,
			payload.ItemName,
			payload.Quantity,
			payload.Unit,
			payload.Priority,
			payload.Id,
			payload.CategoryId,
			userId,
		).Scan(
			&updatedItem.Id,
			&updatedItem.CategoryId,
			&updatedItem.ItemName,
			&updatedItem.Quantity,
			&updatedItem.Unit,
			&updatedItem.Priority,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "AI item not found or not authorized"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update AI item"})
		}

		return c.Status(fiber.StatusOK).JSON(updatedItem)
	}
}

func ConfirmAICategory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "userId is not found"})
		}

		var payload struct {
			CategoryId string `json:"categoryId" validate:"required,uuid4"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input data"})
		}

		tx, err := database.DBPool.Begin(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
		}
		defer tx.Rollback(ctx)

		// Fetch AICategory and AIItems
		var aiCat models.AICategory
		var suggestionId uuid.UUID

		err = tx.QueryRow(ctx,
			`SELECT c."id", c."suggestionId", c."categoryName", c."priority"
			FROM "AICategory" c
			JOIN "AiSuggestion" s
			ON c."suggestionId" = s."id"
			WHERE c."id" = $1
			AND s."userId" = $2`,
			payload.CategoryId,
			userId,
		).Scan(&aiCat.Id, &suggestionId, &aiCat.CategoryName, &aiCat.Priority)

		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "AI category not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch AI category"})
		}

		// Fetch AIItems for this category
		itemRows, err := tx.Query(ctx,
			`SELECT "id", "categoryId", "itemName", "quantity", "unit", "priority"
			FROM "AIItem"
			WHERE "categoryId" = $1`,
			payload.CategoryId,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch AI items"})
		}
		defer itemRows.Close()

		var aiItems []models.AIItem
		for itemRows.Next() {
			var it models.AIItem
			if err := itemRows.Scan(&it.Id, &it.CategoryId, &it.ItemName, &it.Quantity, &it.Unit, &it.Priority); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan AI item"})
			}
			aiItems = append(aiItems, it)
		}

		if err := itemRows.Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error reading AI item rows"})
		}

		// Create regular Category
		categoryId := uuid.New()
		now := time.Now()

		_, err = tx.Exec(ctx,
			`INSERT INTO "Category" ("categoryId", "categoryName", "description", "status", "userId", "icon", "priority", "createdAt", "updatedAt")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			categoryId,
			aiCat.CategoryName,
			"Confirmed from AI suggestion: "+aiCat.CategoryName,
			"active",
			userId,
			"shopping_cart",
			aiCat.Priority,
			now,
			now,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create category"})
		}

		// Create ShoppingItems
		for _, aiItem := range aiItems {
			itemId := uuid.New()
			_, err := tx.Exec(ctx,
				`INSERT INTO "ShoppingItem" ("itemId", "itemName", "description", "status", "categoryId", "priority", "quantity", "unit", "createdAt", "updatedAt")
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
				itemId,
				aiItem.ItemName,
				"Added from AI suggestion",
				"active",
				categoryId.String(),
				aiItem.Priority,
				aiItem.Quantity,
				aiItem.Unit,
				now,
				now,
			)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create shopping item"})
			}
		}

		// Delete all AI rows for this suggestion explicitly
		_, err = tx.Exec(ctx,
			`DELETE FROM "AIItem"
			WHERE "categoryId" IN (
				SELECT "id"
				FROM "AICategory"
				WHERE "suggestionId" = $1
			)`,
			suggestionId,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete AI items"})
		}

		_, err = tx.Exec(ctx,
			`DELETE FROM "AICategory"
			WHERE "suggestionId" = $1`,
			suggestionId,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete AI categories"})
		}

		_, err = tx.Exec(ctx,
			`DELETE FROM "AiSuggestion"
			WHERE "id" = $1
			AND "userId" = $2`,
			suggestionId,
			userId,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete AI suggestion"})
		}

		if err := tx.Commit(ctx); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit confirmation"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "AI category confirmed and saved successfully",
		})
	}
}

func GenetateAIPrompt() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "userId is not found"})
		}

		var payload struct {
			Prompt string `json:"prompt" validate:"required"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input data"})
		}

		if err := validate.Struct(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation Failed", "details": err.Error()})
		}

		err := utils.GenerateAI(ctx, payload.Prompt, userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "AI prompt generated"})
	}
}

func DeleteAISuggestion() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		suggestionId := c.Params("suggestionId")

		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "userId is not found"})
		}

		if suggestionId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "suggestionId is not found"})
		}

		var deletedSuggestion models.AiSuggestion

		tx, err := database.DBPool.Begin(ctx)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
		}
		defer tx.Rollback(ctx)

		_, err = tx.Exec(ctx,
			`DELETE FROM "AIItem"
			WHERE "categoryId" IN (
				SELECT "id"
				FROM "AICategory"
				WHERE "suggestionId" = $1
			)`,
			suggestionId,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete AI items"})
		}

		_, err = tx.Exec(ctx,
			`DELETE FROM "AICategory"
			WHERE "suggestionId" = $1`,
			suggestionId,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete AI categories"})
		}

		err = tx.QueryRow(ctx,
			`DELETE FROM "AiSuggestion"
			WHERE "id" = $1
			AND "userId" = $2
			RETURNING "id", "userId", "prompt"`,
			suggestionId,
			userId,
		).Scan(
			&deletedSuggestion.Id,
			&deletedSuggestion.UserId,
			&deletedSuggestion.Prompt,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "AI suggestion not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete AI suggestion"})
		}

		if err := tx.Commit(ctx); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit AI suggestion deletion"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "AI suggestion deleted successfully",
		})
	}
}

func ReGenetateAIPrompt() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First call the DeleteAISuggestion handler
		if err := DeleteAISuggestion()(c); err != nil {
			return err
		}

		// Then call the GenetateAIPrompt handler
		if err := GenetateAIPrompt()(c); err != nil {
			return err
		}

		return nil
	}
}
