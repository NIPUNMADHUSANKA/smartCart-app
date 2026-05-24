package models

import (
	"time"

	"github.com/google/uuid"
)

type CategoryStatus string
type ItemStatus string
type PriorityStatus string
type UnitStatus string

const (
	CategoryStatusActive    CategoryStatus = "active"
	CategoryStatusCompleted CategoryStatus = "completed"
)
const (
	ItemStatusActive    ItemStatus = "active"
	ItemStatusCompleted ItemStatus = "archived"
)

const (
	PriorityStatusNormal PriorityStatus = "normal"
	PriorityStatusLow    PriorityStatus = "low"
	PriorityStatusMedium PriorityStatus = "medium"
	PriorityStatusHigh   PriorityStatus = "high"
)

const (
	Kg     UnitStatus = "kg"
	Piece  UnitStatus = "piece"
	Pack   UnitStatus = "pack"
	Dozen  UnitStatus = "dozen"
	Box    UnitStatus = "box"
	G      UnitStatus = "g"
	L      UnitStatus = "l"
	Ml     UnitStatus = "ml"
	Bottle UnitStatus = "bottle"
	Can    UnitStatus = "can"
	Cup    UnitStatus = "cup"
	Other  UnitStatus = "other"
)

type ShoppingItem struct {
	ItemId      uuid.UUID      `json:"item_id" validate:"required"`
	ItemName    string         `json:"item_name" validate:"required"`
	Description string         `json:"description" validate:"required,min=2,max=500"`
	Status      ItemStatus     `json:"status" validate:"oneof=active archived"`
	CategoryId  string         `json:"category_id" validate:"required"`
	Priority    PriorityStatus `json:"priority" validate:"oneof=normal low medium high"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updateded_at"`
	Quantity    int            `json:"quantity" validate:"required"`
	Unit        UnitStatus     `json:"unit" validate:"required"`
}

type Category struct {
	CategoryId    uuid.UUID      `json:"category_id"`
	CategoryName  string         `json:"category_name" validate:"required"`
	Description   string         `json:"description"`
	Status        CategoryStatus `json:"status" validate:"oneof=active completed"`
	UserId        string         `json:"user_id"`
	Icon          string         `json:"icon" validate:"required"`
	Priority      PriorityStatus `json:"priority" validate:"oneof=normal low medium high"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updateded_at"`
	ShoppingItems []ShoppingItem `json:"shopping_items"`
}
