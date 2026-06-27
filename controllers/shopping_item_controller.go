package controllers

import (
	"context"
	"smartCart-app/database"
	"smartCart-app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetShoppingItems() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "UserId is not found",
			})
		}

		rows, err := database.DBPool.Query(
			ctx,
			`SELECT si."itemId", si."itemName", si."description", si."status", si."categoryId", si."priority", si."createdAt", si."updatedAt", si."quantity", si."unit"
			FROM "ShoppingItem" si
			INNER JOIN "Category" c ON si."categoryId" = c."categoryId"
			WHERE c."userId" = $1`,
			userId,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		defer rows.Close()

		var shoppingItems []models.ShoppingItem

		for rows.Next() {
			var shoppingItem models.ShoppingItem
			if err := rows.Scan(
				&shoppingItem.ItemId,
				&shoppingItem.ItemName,
				&shoppingItem.Description,
				&shoppingItem.Status,
				&shoppingItem.CategoryId,
				&shoppingItem.Priority,
				&shoppingItem.CreatedAt,
				&shoppingItem.UpdatedAt,
				&shoppingItem.Quantity,
				&shoppingItem.Unit,
			); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to scan Shopping Item",
				})
			}
			shoppingItems = append(shoppingItems, shoppingItem)
		}
		if err := rows.Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error reading Shopping Item rows",
			})
		}

		return c.Status(fiber.StatusOK).JSON(shoppingItems)

	}
}

func GetShoppingItemByItemId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "UserId is not found",
			})
		}

		itemId := c.Params("itemId")

		if itemId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Shopping Item Id is not found",
			})
		}

		var foundShoppingItem models.ShoppingItem

		err := database.DBPool.QueryRow(
			ctx,
			`SELECT si."itemId", si."itemName", si."description", si."status", si."categoryId", si."priority", si."createdAt", si."updatedAt", si."quantity", si."unit"
			FROM "ShoppingItem" si
			INNER JOIN "Category" c ON si."categoryId" = c."categoryId"
			WHERE c."userId" = $1
			AND si."itemId" = $2`,
			userId, itemId,
		).Scan(
			&foundShoppingItem.ItemId,
			&foundShoppingItem.ItemName,
			&foundShoppingItem.Description,
			&foundShoppingItem.Status,
			&foundShoppingItem.CategoryId,
			&foundShoppingItem.Priority,
			&foundShoppingItem.CreatedAt,
			&foundShoppingItem.UpdatedAt,
			&foundShoppingItem.Quantity,
			&foundShoppingItem.Unit,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(foundShoppingItem)

	}
}

func GetShoppingItemByCategoryId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "UserId is not found",
			})
		}

		categoryId := c.Params("categoryId")

		if categoryId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Category Id is not found",
			})
		}

		var foundShoppingItem models.ShoppingItem

		err := database.DBPool.QueryRow(
			ctx,
			`SELECT si."itemId", si."itemName", si."description", si."status", si."categoryId", si."priority", si."createdAt", si."updatedAt", si."quantity", si."unit"
			FROM "ShoppingItem" si
			INNER JOIN "Category" c ON si."categoryId" = c."categoryId"
			WHERE c."userId" = $1
			AND si."categoryId" = $2`,
			userId, categoryId,
		).Scan(
			&foundShoppingItem.ItemId,
			&foundShoppingItem.ItemName,
			&foundShoppingItem.Description,
			&foundShoppingItem.Status,
			&foundShoppingItem.CategoryId,
			&foundShoppingItem.Priority,
			&foundShoppingItem.CreatedAt,
			&foundShoppingItem.UpdatedAt,
			&foundShoppingItem.Quantity,
			&foundShoppingItem.Unit,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(foundShoppingItem)

	}
}

func DeleteShoppingItemByItemId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second*100)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)

		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "UserId is not found",
			})
		}

		itemId := c.Params("itemId")

		if itemId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Shopping Item Id is not found",
			})
		}

		var deletedShoppingItem models.ShoppingItem

		err := database.DBPool.QueryRow(
			ctx,
			`DELETE FROM "ShoppingItem" si
			USING "Category" c
			WHERE si."categoryId" = c."categoryId"
			AND c."userId" = $1
			AND si."itemId" = $2
			RETURNING si."itemId", si."itemName", si."description", si."status", si."categoryId", si."priority", si."createdAt", si."updatedAt", si."quantity", si."unit"`,
			userId,
			itemId,
		).Scan(
			&deletedShoppingItem.ItemId,
			&deletedShoppingItem.ItemName,
			&deletedShoppingItem.Description,
			&deletedShoppingItem.Status,
			&deletedShoppingItem.CategoryId,
			&deletedShoppingItem.Priority,
			&deletedShoppingItem.CreatedAt,
			&deletedShoppingItem.UpdatedAt,
			&deletedShoppingItem.Quantity,
			&deletedShoppingItem.Unit,
		)

		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Shopping Item not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete shopping item"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":       "Shopping Item deleted successfully",
			"shopping_item": deletedShoppingItem,
		})
	}
}

func CreateShoppingItems() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var shoppingItem models.ShoppingItem

		if err := c.BodyParser(&shoppingItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input data",
			})
		}

		shoppingItem.ItemId = uuid.New()
		shoppingItem.CreatedAt = time.Now()
		shoppingItem.UpdatedAt = time.Now()

		err := database.DBPool.QueryRow(
			ctx,
			`INSERT INTO "ShoppingItem" ("itemId", "itemName", "description", "status", "categoryId", "priority", "createdAt", "updatedAt", "quantity", "unit")
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING "itemId"`,
			shoppingItem.ItemId,
			shoppingItem.ItemName,
			shoppingItem.Description,
			shoppingItem.Status,
			shoppingItem.CategoryId,
			shoppingItem.Priority,
			shoppingItem.CreatedAt,
			shoppingItem.UpdatedAt,
			shoppingItem.Quantity,
			shoppingItem.Unit,
		).Scan(&shoppingItem.ItemId)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create Shopping Item",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(shoppingItem)
	}
}

func UpdateShoppingItem() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		shoppingItemId := c.Params("itemId")
		if shoppingItemId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Shopping Item Id not found",
			})
		}

		itemUUID, err := uuid.Parse(shoppingItemId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid shopping item id",
			})
		}

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "UserId not found",
			})
		}

		// Step 1: Fetch the existing record
		var existingItem models.ShoppingItem
		err = database.DBPool.QueryRow(
			ctx,
			`SELECT si."itemId", si."itemName", si."description", si."status", si."categoryId", si."priority", si."createdAt", si."updatedAt", si."quantity", si."unit"
			FROM "ShoppingItem" si
			INNER JOIN "Category" c ON si."categoryId" = c."categoryId"
			WHERE c."userId" = $1
			AND si."itemId" = $2`,
			userId, itemUUID,
		).Scan(
			&existingItem.ItemId,
			&existingItem.ItemName,
			&existingItem.Description,
			&existingItem.Status,
			&existingItem.CategoryId,
			&existingItem.Priority,
			&existingItem.CreatedAt,
			&existingItem.UpdatedAt,
			&existingItem.Quantity,
			&existingItem.Unit,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Shopping Item not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Step 2: Parse request body and overlay new values onto the existing record
		var input models.ShoppingItem
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input data",
			})
		}

		if input.ItemName != "" {
			existingItem.ItemName = input.ItemName
		}
		if input.Description != "" {
			existingItem.Description = input.Description
		}
		if input.Status != "" {
			existingItem.Status = input.Status
		}
		if input.CategoryId != "" {
			existingItem.CategoryId = input.CategoryId
		}
		if input.Priority != "" {
			existingItem.Priority = input.Priority
		}
		if input.Quantity != 0 {
			existingItem.Quantity = input.Quantity
		}
		if input.Unit != "" {
			existingItem.Unit = input.Unit
		}
		existingItem.UpdatedAt = time.Now()

		if err := validate.Struct(existingItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Step 3: Update the record with merged values
		err = database.DBPool.QueryRow(
			ctx,
			`UPDATE "ShoppingItem" si
			SET "itemName" = $1,
			    "description" = $2,
			    "status" = $3,
			    "categoryId" = $4,
			    "priority" = $5,
			    "updatedAt" = $6,
			    "quantity" = $7,
			    "unit" = $8
			FROM "Category" c
			WHERE si."categoryId" = c."categoryId"
			  AND c."userId" = $9
			  AND si."itemId" = $10
			  AND EXISTS (
			      SELECT 1 FROM "Category" c2
			      WHERE c2."categoryId" = $4
			        AND c2."userId" = $9
			  )
			RETURNING si."itemId", si."itemName", si."description", si."status", si."categoryId", si."priority", si."createdAt", si."updatedAt", si."quantity", si."unit"`,
			existingItem.ItemName,
			existingItem.Description,
			existingItem.Status,
			existingItem.CategoryId,
			existingItem.Priority,
			existingItem.UpdatedAt,
			existingItem.Quantity,
			existingItem.Unit,
			userId,
			itemUUID,
		).Scan(
			&existingItem.ItemId,
			&existingItem.ItemName,
			&existingItem.Description,
			&existingItem.Status,
			&existingItem.CategoryId,
			&existingItem.Priority,
			&existingItem.CreatedAt,
			&existingItem.UpdatedAt,
			&existingItem.Quantity,
			&existingItem.Unit,
		)

		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Shopping Item not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update shopping item",
			})
		}

		return c.Status(fiber.StatusOK).JSON(existingItem)
	}
}
