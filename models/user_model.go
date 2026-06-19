package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId       uuid.UUID `json:"userId"`
	FullName     string    `json:"fullName" validate:"required,min=2,max=100"`
	UserName     string    `json:"userName" validate:"required,min=2,max=100"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password" validate:"required,min=6"`
	Role         string    `json:"role" validate:"oneof=ADMIN USER"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Token        string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

type UserLogin struct {
	UserName string `json:"userName" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=6"`
}

type PasswordUpdate struct {
	Password    string `json:"currentPassword" validate:"required,min=6"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

type UserResponse struct {
	UserId       uuid.UUID `json:"userId"`
	UserName     string    `json:"userName"`
	Token        string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

type UserRegisterRes struct {
	UserId    uuid.UUID `json:"userId"`
	UserName  string    `json:"userName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
}

type UserInfo struct {
	UserId    uuid.UUID `json:"userId"`
	UserName  string    `json:"userName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
}
