package entity

import (
	"database/sql"
	"time"
)

type PermissionModel struct {
	ID          sql.NullInt64  `db:"id"`
	Name        sql.NullString `db:"name"`
	IsActive    sql.NullBool   `db:"is_active"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

type PermissionEntity struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name" db:"name"`
	IsActive    bool   `json:"is_active" db:"is_active"`
	Description string `json:"description" db:"description"`
}

type PermissionResponse struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Services    []ServiceResponse `json:"service"`
	IsActive    bool              `json:"is_active"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type PermissionUpdateRequest struct {
	ID       int64  `json:"-"`
	Name     string `json:"name" validate:"required" db:"name"`
	IsActive bool   `json:"is_active"  db:"is_active"`
	// ServicesID  []int64 `json:"services_id"`
	ServiceID   int64  `json:"service_id" db:"service_id"`
	Description string `json:"description" db:"description"`
}

func (p PermissionModel) ToEntity() PermissionEntity {
	return PermissionEntity{
		ID:          p.ID.Int64,
		Name:        p.Name.String,
		IsActive:    p.IsActive.Bool,
		Description: p.Description.String,
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
		DeletedAt:   p.DeletedAt.Time,
	}
}

func (entity PermissionEntity) ToResponse() PermissionResponse {
	return PermissionResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		IsActive:    entity.IsActive,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
