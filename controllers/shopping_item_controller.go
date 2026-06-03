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
			`SELECT si."ItemId", si."ItemName", si."Description", si."Status", si."CategoryId", si."Priority", si."CreatedAt", si."UpdatedAt", si."Quantity", si."Unit"
			FROM "shoppingItem" si
			INNER JOIN "Category" c ON si."CategoryId" = c."categoryId"
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
			`SELECT si."ItemId", si."ItemName", si."Description", si."Status", si."CategoryId", si."Priority", si."CreatedAt", si."UpdatedAt", si."Quantity", si."Unit"
			FROM "shoppingItem" si
			INNER JOIN "Category" c ON si."CategoryId" = c."categoryId"
			WHERE c."userId" = $1
			AND si."ItemId" = $2`,
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
			`SELECT si."ItemId", si."ItemName", si."Description", si."Status", si."CategoryId", si."Priority", si."CreatedAt", si."UpdatedAt", si."Quantity", si."Unit"
			FROM "shoppingItem" si
			INNER JOIN "Category" c ON si."CategoryId" = c."categoryId"
			WHERE c."userId" = $1
			AND si."CategoryId" = $2`,
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
			`DELETE FROM "shoppingItem" si
			USING "Category" c
			WHERE si."CategoryId" = c."categoryId"
			AND c."userId" = $1
			AND si."ItemId" = $2
			RETURNING si."ItemId", si."ItemName", si."Description", si."Status", si."CategoryId", si."Priority", si."CreatedAt", si."UpdatedAt", si."Quantity", si."Unit"`,
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
			`INSERT INTO "shoppingItem" ("ItemId", "ItemName", "Description", "Status", "CategoryId", "Priority", "CreatedAt", "UpdatedAt", "Quantity", "Unit")
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING "ItemId"`,
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

		var shoppingItem models.ShoppingItem

		if err := c.BodyParser(&shoppingItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input data",
			})
		}

		itemUUID, err := uuid.Parse(shoppingItemId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid shopping item id",
			})
		}

		shoppingItem.ItemId = itemUUID

		if err := validate.Struct(shoppingItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation Failed",
				"details": err.Error(),
			})
		}

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)

		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "UserId not found",
			})
		}

		var updatedCategory models.ShoppingItem

		err = database.DBPool.QueryRow(
			ctx,
			`UPDATE "shoppingItem" si
			SET "ItemName" = $1,
			    "Description" = $2,
			    "Status" = $3,
			    "CategoryId" = $4,
			    "Priority" = $5,
			    "UpdatedAt" = $6,
			    "Quantity" = $7,
			    "Unit" = $8
			FROM "Category" c
			WHERE si."CategoryId" = c."categoryId"
			  AND c."userId" = $9
			  AND si."ItemId" = $10
			  AND EXISTS (
			      SELECT 1 FROM "Category" c2
			      WHERE c2."categoryId" = $4
			        AND c2."userId" = $9
			  )
			RETURNING si."ItemId", si."ItemName", si."Description", si."Status", si."CategoryId", si."Priority", si."CreatedAt", si."UpdatedAt", si."Quantity", si."Unit`,
			shoppingItem.ItemName,
			shoppingItem.Description,
			shoppingItem.Status,
			shoppingItem.CategoryId,
			shoppingItem.Priority,
			time.Now(),
			shoppingItem.Quantity,
			shoppingItem.Unit,
			userId,
			itemUUID,
		).Scan(
			&updatedCategory.ItemId,
			&updatedCategory.ItemName,
			&updatedCategory.Description,
			&updatedCategory.Status,
			&updatedCategory.CategoryId,
			&updatedCategory.Priority,
			&updatedCategory.CreatedAt,
			&updatedCategory.UpdatedAt,
			&updatedCategory.Quantity,
			&updatedCategory.Unit,
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

		return c.Status(fiber.StatusOK).JSON(updatedCategory)
	}
}
