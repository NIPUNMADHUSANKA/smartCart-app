package utils

import (
	"context"
	"errors"
	"os"
	"smartCart-app/database"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SignedDetails struct {
	UserId   uuid.UUID
	Email    string
	FullName string
	Role     string
	UserName string
	jwt.RegisteredClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var SECRET_REFRESH_KEY string = os.Getenv("SECRET_REFRESH_KEY")

func GernerateAllTokens(userId uuid.UUID, UserName, email, FullName, Role string) (string, string, error) {
	claims := &SignedDetails{
		UserId:   userId,
		Email:    email,
		FullName: FullName,
		UserName: UserName,
		Role:     Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "SmartCart",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refreshclaims := &SignedDetails{
		UserId:   userId,
		Email:    email,
		FullName: FullName,
		UserName: UserName,
		Role:     Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "SmartCart",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims)
	signedrefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))

	if err != nil {
		return "", "", err
	}

	return signedToken, signedrefreshToken, nil

}

func UpdateAllTokend(userName, token, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var userId uuid.UUID

	err := database.DBPool.QueryRow(
		ctx,
		`UPDATE "User" 
		 SET "token" = $1,
		 "refreshToken" = $2,
		 "updatedAt" = $3
		 WHERE "userName" = $4
		RETURNING "userId"`,
		token, refreshToken, time.Now(), userName,
	).Scan(&userId)

	if err != nil {
		return err
	}
	return nil
}

func GetAccessToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}

	tokenstring := strings.Split(authHeader, "Bearer ")
	if len(tokenstring) == 0 {
		return "", errors.New("Bearer token is required")
	}

	return tokenstring[1], nil
}

func ValidateToken(tokenstring string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenstring, claims, func(t *jwt.Token) (any, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
