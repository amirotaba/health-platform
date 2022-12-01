package entity

import (
	"database/sql"
	"time"
)

type RoleModel struct {
	ID          sql.NullInt64  `db:"id"`
	Name        sql.NullString `db:"name"`
	IsActive    sql.NullBool   `db:"is_active"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

type RoleEntity struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name" validate:"required"`
	IsActive    bool               `json:"is_active"`
	Permissions []PermissionEntity `json:"permissions"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   time.Time          `json:"deleted_at"`
}

type RoleResponse struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	IsActive    bool               `json:"is_active"`
	Permissions []PermissionEntity `json:"permissions,omitempty"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type RoleUpdateRequest struct {
	ID          int64  `json:"-"`
	Name        string `json:"name" validate:"required"`
	IsActive    bool   `json:"is_active"`
	Description string `json:"description"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required" db:"name"`
	IsActive    bool   `json:"is_active" db:"is_active"`
	Description string `json:"description"`
}

func (model RoleModel) ToEntity() RoleEntity {
	return RoleEntity{
		ID:          model.ID.Int64,
		Name:        model.Name.String,
		IsActive:    model.IsActive.Bool,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
		DeletedAt:   model.DeletedAt.Time,
	}
}

func (entity RoleEntity) ToResponse() RoleResponse {
	return RoleResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		IsActive:    entity.IsActive,
		Permissions: entity.Permissions,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
