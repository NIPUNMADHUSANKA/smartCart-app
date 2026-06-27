package controllers

import (
	"context"
	"smartCart-app/database"
	"smartCart-app/models"
	"smartCart-app/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func ResetPassword() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userIdStr, ok := userIdInterface.(string)
		if !ok || userIdStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON((fiber.Map{
				"error": "UserId not found",
			}))
		}

		var passwordupdate models.PasswordUpdate

		if err := c.BodyParser(&passwordupdate); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input data",
			})
		}

		if err := validate.Struct(passwordupdate); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation Failed",
				"details": err.Error(),
			})
		}

		var storedPassword string
		err := database.DBPool.QueryRow(
			ctx,
			`SELECT "password" FROM "User" WHERE "userId" = $1`,
			userIdStr,
		).Scan(&storedPassword)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "User does not found",
				})
			}
			return err
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(passwordupdate.Password)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Current Password is Invalid",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(passwordupdate.NewPassword)); err == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "New Password must be different from Current Password",
			})
		}

		hashedPassword, err := HashPassword(passwordupdate.NewPassword)
		if err != nil {
			return err
		}

		_, err = database.DBPool.Exec(
			ctx,
			`UPDATE "User" SET "password" = $1, "updatedAt" = $2 WHERE "userId" = $3`,
			hashedPassword, time.Now(), userIdStr,
		)
		if err != nil {
			return err
		}

		return nil

	}
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

		user.Role = "USER"

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
			`SELECT EXISTS(SELECT 1 FROM "User" WHERE "userName" = $1 OR "email" = $2)`,
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
			`INSERT INTO "User" ("userId","fullName", "userName", "email", "password", "createdAt", "updatedAt", "role")
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING "userId"`,
			userID, user.FullName, user.UserName, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.Role,
		).Scan(&userID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}

		var userRes models.UserRegisterRes = models.UserRegisterRes{
			UserId:    userID,
			UserName:  user.UserName,
			FullName:  user.FullName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		return c.Status(fiber.StatusCreated).JSON(userRes)
	}
}

func LoginUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var UserLogin models.UserLogin

		if err := c.BodyParser(&UserLogin); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input data",
			})
		}

		if err := validate.Struct(UserLogin); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation Failed",
				"details": err.Error(),
			})
		}

		var foundUser models.User

		err := database.DBPool.QueryRow(
			ctx,
			`SELECT "userId","fullName","userName","email","password","role"
     FROM "User" WHERE "userName" = $1`,
			UserLogin.UserName,
		).Scan(
			&foundUser.UserId,
			&foundUser.FullName,
			&foundUser.UserName,
			&foundUser.Email,
			&foundUser.Password,
			&foundUser.Role,
		)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(UserLogin.Password))

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		token, refreshToken, err := utils.GernerateAllTokens(foundUser.UserId, foundUser.UserName, foundUser.Email, foundUser.FullName, foundUser.Role)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = utils.UpdateAllTokend(foundUser.UserName, token, refreshToken)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update tokens",
			})
		}

		var response models.UserResponse = models.UserResponse{
			UserId:       foundUser.UserId,
			UserName:     foundUser.UserName,
			Token:        token,
			RefreshToken: refreshToken,
		}

		return c.Status(fiber.StatusOK).JSON(response)

	}
}

func GetUserInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {

		userIdInterface := c.Locals("userId")
		userId, ok := userIdInterface.(string)

		if !ok || userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON((fiber.Map{
				"error": "UserId not found",
			}))
		}

		userName := c.Locals("userName")

		if userName == "" || userName == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "UserName not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"userId":   userId,
			"userName": userName,
		})

	}
}

func GetUserDetails() fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userIdStr, ok := userIdInterface.(string)
		if !ok || userIdStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON((fiber.Map{
				"error": "UserId not found",
			}))
		}

		userNameInterface := c.Locals("userName")
		userName, ok := userNameInterface.(string)
		if !ok || userName == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "UserName not found",
			})
		}

		parsedUserID, err := uuid.Parse(userIdStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		var foundUser models.User

		err = database.DBPool.QueryRow(
			ctx,
			`SELECT "fullName","email","role", "token", "refreshToken", "createdAt", "updatedAt"
      		FROM "User" WHERE "userName" = $1 and "userId" = $2`,
			userName, parsedUserID,
		).Scan(
			&foundUser.FullName,
			&foundUser.Email,
			&foundUser.Role,
			&foundUser.Token,
			&foundUser.RefreshToken,
			&foundUser.CreatedAt,
			&foundUser.UpdatedAt,
		)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var response models.UserInfo = models.UserInfo{
			UserId:    parsedUserID,
			UserName:  userName,
			FullName:  foundUser.FullName,
			Email:     foundUser.Email,
			CreatedAt: foundUser.CreatedAt,
			UpdatedAt: foundUser.UpdatedAt,
		}
		return c.Status(fiber.StatusOK).JSON(response)

	}

}

func DeleteUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userIdInterface := c.Locals("userId")
		userIdStr, ok := userIdInterface.(string)
		if !ok || userIdStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON((fiber.Map{
				"error": "UserId not found",
			}))
		}

		// Channel to collect errors from concurrent operations
		errChan := make(chan error, 5)

		// Delete AI items concurrently
		go func() {
			_, err := database.DBPool.Exec(ctx,
				`DELETE FROM "AIItem"
				WHERE "categoryId" IN (
					SELECT "id"
					FROM "AICategory"
					WHERE "suggestionId" IN (
						SELECT "id"
						FROM "AiSuggestion"
						WHERE "userId" = $1
					)
				)`,
				userIdStr,
			)
			errChan <- err
		}()

		// Delete shopping items concurrently
		go func() {
			_, err := database.DBPool.Exec(ctx,
				`DELETE FROM "ShoppingItem" si
				USING "Category" c
				WHERE si."categoryId" = c."categoryId"
				AND c."userId" = $1`,
				userIdStr,
			)
			errChan <- err
		}()

		// Delete AI categories concurrently
		go func() {
			_, err := database.DBPool.Exec(ctx,
				`DELETE FROM "AICategory"
				WHERE "suggestionId" IN (
					SELECT "id"
					FROM "AiSuggestion"
					WHERE "userId" = $1
				)`,
				userIdStr,
			)
			errChan <- err
		}()

		// Delete AI suggestions concurrently
		go func() {
			_, err := database.DBPool.Exec(ctx,
				`DELETE FROM "AiSuggestion"
				WHERE "userId" = $1`,
				userIdStr,
			)
			errChan <- err
		}()

		// Delete categories concurrently
		go func() {
			_, err := database.DBPool.Exec(ctx,
				`DELETE FROM "Category"
				WHERE "userId" = $1`,
				userIdStr,
			)
			errChan <- err
		}()

		// Collect errors from concurrent operations
		for i := 0; i < 5; i++ {
			if err := <-errChan; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to delete user related data",
				})
			}
		}

		// Finally, delete the user record
		res, err := database.DBPool.Exec(ctx,
			`DELETE FROM "User"
			WHERE "userId" = $1`,
			userIdStr,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete user",
			})
		}

		if res.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User and all related data deleted successfully",
		})
	}
}
