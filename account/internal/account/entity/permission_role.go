package entity

import (
	"database/sql"
	"time"
)

type RolePermissionModel struct {
	ID           sql.NullInt64  `db:"id"`
	RoleID       sql.NullInt64  `db:"role_id"`
	PermissionID sql.NullInt64  `db:"permission_id"`
	IsActive     sql.NullBool   `db:"is_active"`
	Description  sql.NullString `db:"description"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	DeletedAt    sql.NullTime   `db:"deleted_at"`
}

type RolePermissionEntity struct {
	ID           int64     `json:"id"`
	RoleID       int64     `json:"role_id"`
	PermissionID int64     `json:"permission_id"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type CreatePermissionsToRoleRequest struct {
	RoleId       int64   `json:"role_id" db:"role_id"`
	PermissionId int64   `json:"-" db:"permission_id"`
	Permissions  []int64 `json:"permissions" db:"-"`
	Description  string  `json:"description" db:"description"`
}

type UpdateRolePermissionsRequest struct {
	RoleId       int64   `json:"-" db:"role_id"`
	PermissionId int64   `json:"-" db:"permission_id"`
	Permissions  []int64 `json:"permissions" db:"-"`
}

type RolePermissionsResponse struct {
	ID           int64     `json:"id" db:"id"`
	RoleID       int64     `json:"role_id" db:"role_id"`
	PermissionID int64     `json:"permission_id" db:"permission_id"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
