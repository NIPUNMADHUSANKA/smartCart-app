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

func GetCategory() fiber.Handler {
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

		rows, err := database.DBPool.Query(
			ctx,
			`SELECT c."categoryId", c."categoryName", c."description", c."status", c."userId", c."icon", c."priority", c."createdAt", c."updatedAt"
			FROM "Category" c
			WHERE c."userId" = $1`,
			userId,
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch categories",
			})
		}

		defer rows.Close()

		var categories []models.Category

		for rows.Next() {
			var category models.Category
			if err := rows.Scan(
				&category.CategoryId,
				&category.CategoryName,
				&category.Description,
				&category.Status,
				&category.UserId,
				&category.Icon,
				&category.Priority,
				&category.CreatedAt,
				&category.UpdatedAt,
			); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to scan category",
				})
			}
			categories = append(categories, category)
		}

		if err := rows.Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error reading category rows",
			})
		}

		return c.Status(fiber.StatusOK).JSON(categories)

	}
}

func GetCategoryByCategoryId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		categoryId := c.Params("categoryId")

		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "userId is not found",
			})
		}

		if categoryId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "categoryId is not found",
			})
		}

		var foundCategory models.Category

		err := database.DBPool.QueryRow(
			ctx,
			`SELECT c."categoryId", c."categoryName", c."description", c."status", c."userId", c."icon", c."priority", c."createdAt", c."updatedAt"
			FROM "Category" c
			WHERE c."userId" = $1
			AND c."categoryId" = $2`,
			userId, categoryId,
		).Scan(
			&foundCategory.CategoryId,
			&foundCategory.CategoryName,
			&foundCategory.Description,
			&foundCategory.Status,
			&foundCategory.UserId,
			&foundCategory.Icon,
			&foundCategory.Priority,
			&foundCategory.CreatedAt,
			&foundCategory.UpdatedAt,
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch Category",
			})
		}

		return c.Status(fiber.StatusOK).JSON(foundCategory)

	}
}

func DeleteCategoryByCategoryId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		categoryId := c.Params("categoryId")

		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "userId is not found",
			})
		}

		if categoryId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "categoryId is not found",
			})
		}

		var deletedCategory models.Category

		err := database.DBPool.QueryRow(
			ctx,
			`DELETE FROM "Category"
			WHERE "userId" = $1
			AND "categoryId" = $2
			RETURNING "categoryId", "categoryName", "description", "status", "userId", "icon", "priority", "createdAt", "updatedAt"`,
			userId, categoryId,
		).Scan(
			&deletedCategory.CategoryId,
			&deletedCategory.CategoryName,
			&deletedCategory.Description,
			&deletedCategory.Status,
			&deletedCategory.UserId,
			&deletedCategory.Icon,
			&deletedCategory.Priority,
			&deletedCategory.CreatedAt,
			&deletedCategory.UpdatedAt,
		)

		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Category not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete category",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "Category deleted successfully",
			"category": deletedCategory,
		})
	}
}

func CreateCategory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var category models.Category

		if err := c.BodyParser(&category); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input data",
			})
		}

		if err := validate.Struct(category); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation Failed",
				"details": err.Error(),
			})
		}

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)
		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "userId not found",
			})
		}

		category.CategoryId = uuid.New()
		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()
		category.UserId = userId

		var CategoryId uuid.UUID
		err := database.DBPool.QueryRow(
			ctx,
			`INSERT INTO "Category" ("categoryId", "categoryName","description", "status", "userId", "icon", "priority", "createdAt", "updatedAt")
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING "categoryId"`,
			category.CategoryId, category.CategoryName, category.Description, category.Status, category.UserId, category.Icon, category.Priority, category.CreatedAt, category.UpdatedAt,
		).Scan(&CategoryId)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create category",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(category)
	}
}
