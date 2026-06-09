package models

import (
	"github.com/google/uuid"
)

type AiSuggestion struct {
	Id         uuid.UUID    `json:"id"`
	UserId     string       `json:"user_id" validate:"required,uuid4"`
	Prompt     string       `json:"prompt" validate:"required"`
	Categories []AICategory `json:"categories" validate:"required,dive"`
}

type AICategory struct {
	Id           uuid.UUID      `json:"id"`
	SuggestionId string         `json:"suggestion_id" validate:"required,uuid4"`
	CategoryName string         `json:"category_name" validate:"required"`
	Priority     PriorityStatus `json:"priority" validate:"required,oneof=normal low medium high"`
	Items        []AIItem       `json:"items" validate:"required,dive"`
}

type AIItem struct {
	Id         uuid.UUID      `json:"id"`
	CategoryId string         `json:"category_id" validate:"required,uuid4"`
	ItemName   string         `json:"item_name" validate:"required"`
	Quantity   float32        `json:"quantity" validate:"required,gt=0"`
	Unit       UnitStatus     `json:"unit" validate:"required"`
	Priority   PriorityStatus `json:"priority" validate:"required,oneof=normal low medium high"`
}
