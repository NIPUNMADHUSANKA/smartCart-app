package controllers

import (
	"context"
	"smartCart-app/database"
	"smartCart-app/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(HashPassword), nil
}

func RegisterUser() fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var user models.User

		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input data",
			})
		}

		if err := validate.Struct(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation Failed",
				"details": err.Error(),
			})
		}

		hashPassword, err := HashPassword(user.Password)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to hash password",
			})
		}

		var exists bool

		qderr := database.DBPool.QueryRow(
			ctx,
			"SELECT EXISTS(SELECT 1 FROM \"User\" WHERE \"userName\" = $1 OR \"email\" = $2)",
			user.UserName, user.Email,
		).Scan(&exists)

		if qderr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check existing user",
			})
		}

		if exists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "This email or username is already registered. Please log in instead."})
		}

		user.Password = hashPassword
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		var userID uuid.UUID = uuid.New()

		err = database.DBPool.QueryRow(
			ctx,
			`INSERT INTO "User" ("userId","fullName", "userName", "email", "password", "createdAt", "updatedAt")
			VALUES($1, $2, $3, $4, $5, $6, $7)
			RETURNING "userId"`,
			userID, user.FullName, user.UserName, user.Email, user.Password, user.CreatedAt, user.UpdatedAt,
		).Scan(&userID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}

		var userRes models.UserRegisterRes = models.UserRegisterRes{
			UserId:    userID,
			FullName:  user.FullName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		return c.Status(fiber.StatusCreated).JSON(userRes)
	}
}
