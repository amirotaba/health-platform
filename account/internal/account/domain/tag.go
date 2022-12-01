package domain

import (
	"time"
)

type TagEntity struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name" validate:"required"`
	CategoryId   int64     `json:"category_id"  validate:"required"`
	CategoryName string    `json:"category_name"`
	IsActive     bool      `json:"is_active"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"update_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type TagResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	CategoryId   int64     `json:"category_id"`
	CategoryName string    `json:"category_name"`
	IsActive     bool      `json:"is_active"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"update_at"`
}
