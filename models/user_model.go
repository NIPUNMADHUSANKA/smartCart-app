package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID `json:"user_id"`
	FullName     string    `json:"full_name" validate:"required,min=2,max=100"`
	UserName     string    `json:"user_name" validate:"required,min=2,max=100"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password" validate:"required,min=6"`
	Role         string    `json:"role" validate:"oneof=ADMIN USER"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updateded_at"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	UserId       string `json:"user_id"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserRegisterRes struct {
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updateded_at"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
}
