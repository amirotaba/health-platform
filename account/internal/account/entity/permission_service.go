package entity

import (
	"database/sql"
	"time"
)

type PermissionServicesModel struct {
	ID           sql.NullInt64  `db:"id"`
	PermissionID sql.NullInt64  `db:"permission_id"`
	ServiceID    sql.NullInt64  `db:"service_id"`
	IsActive     sql.NullBool   `db:"is_active"`
	Description  sql.NullString `db:"description"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	DeletedAt    sql.NullTime   `db:"deleted_at"`
}

type PermissionServicesEntity struct {
	ID           int64     `json:"id"`
	PermissionID int64     `json:"permission_id"`
	ServiceID    int64     `json:"service_id"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type AssignServicesToPermission struct {
	PermissionID int64   `json:"permission_id" db:"permission_id"`
	ServiceID    int64   `json:"service_id" db:"service_id"`
	Services     []int64 `json:"services" db:"-"`
	Description  string  `json:"description" db:"description"`
}

type UpdateAssignedServicesToPermission struct {
	PermissionID int64   `json:"-" db:"permission_id"`
	ServiceID    int64   `json:"-" db:"service_id"`
	Services     []int64 `json:"services" db:"-"`
	Description  string  `json:"description" db:"description"`
}

type PermissionServicesResponse struct {
	ID           int64     `json:"ps_id" db:"id"`
	PermissionID int64     `json:"permission_id" db:"permission_id"`
	ServiceID    int64     `json:"service_id" db:"service_id"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type PermissionServicesResponseModel struct {
	ID           sql.NullInt64  `db:"id"`
	PermissionID sql.NullInt64  `db:"permission_id"`
	ServiceID    sql.NullInt64  `db:"service_id"`
	Description  sql.NullString `db:"description"`
	CreatedAt    sql.NullTime   `db:"created_at"`
}
