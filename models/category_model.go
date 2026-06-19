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
	ItemId      uuid.UUID      `json:"itemId" validate:"required"`
	ItemName    string         `json:"itemName" validate:"required"`
	Description string         `json:"description" validate:"required,min=2,max=500"`
	Status      ItemStatus     `json:"status" validate:"oneof=active archived"`
	CategoryId  string         `json:"categoryId" validate:"required"`
	Priority    PriorityStatus `json:"priority" validate:"oneof=normal low medium high"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Quantity    int            `json:"quantity" validate:"required"`
	Unit        UnitStatus     `json:"unit" validate:"required"`
}

type Category struct {
	CategoryId   uuid.UUID      `json:"categoryId"`
	CategoryName string         `json:"categoryName" validate:"required"`
	Description  string         `json:"description"`
	Status       CategoryStatus `json:"status" validate:"oneof=active completed"`
	UserId       string         `json:"userId"`
	Icon         string         `json:"icon" validate:"required"`
	Priority     PriorityStatus `json:"priority" validate:"oneof=normal low medium high"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}
