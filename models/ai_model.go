package models

import (
	"time"

	"github.com/google/uuid"
)

type AiSuggestion struct {
	Id         uuid.UUID    `json:"id"`
	UserId     string       `json:"userId"`
	Prompt     string       `json:"prompt"`
	CreatedAt  time.Time    `json:"createdAt"`
	Categories []AICategory `json:"categories"`
}

type AICategory struct {
	Id           uuid.UUID      `json:"id"`
	SuggestionId string         `json:"suggestionId" validate:"required,uuid4"`
	CategoryName string         `json:"categoryName" validate:"required"`
	Priority     PriorityStatus `json:"priority" validate:"required,oneof=normal low medium high"`
	Items        []AIItem       `json:"items" validate:"required,dive"`
}

type AIItem struct {
	Id         uuid.UUID      `json:"id"`
	CategoryId string         `json:"categoryId" validate:"required,uuid4"`
	ItemName   string         `json:"itemName" validate:"required"`
	Quantity   float32        `json:"quantity" validate:"required,gt=0"`
	Unit       UnitStatus     `json:"unit" validate:"required"`
	Priority   PriorityStatus `json:"priority" validate:"required,oneof=normal low medium high"`
}

type AIItemRec struct {
	Id         uuid.UUID      `json:"itemId"`
	CategoryId string         `json:"categoryId" validate:"required,uuid4"`
	ItemName   string         `json:"itemName" validate:"required"`
	Quantity   float32        `json:"itemQty" validate:"required,gt=0"`
	Unit       UnitStatus     `json:"itemUnit" validate:"required"`
	Priority   PriorityStatus `json:"priority" validate:"required,oneof=normal low medium high"`
}
