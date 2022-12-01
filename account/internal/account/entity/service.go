package entity

import (
	"database/sql"
	"time"
)

type Service struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Path        string    `json:"path"`
	Function    string    `json:"function"`
	Method      string    `json:"method"`
	Active      bool      `json:"active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type ServiceModel struct {
	ID          sql.NullInt64  `db:"id"`
	Name        sql.NullString `db:"name"`
	Code        sql.NullString `db:"code"`
	Path        sql.NullString `db:"path"`
	Function    sql.NullString `db:"func"`
	Method      sql.NullString `db:"method"`
	IsActive    sql.NullBool   `db:"is_active"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

type ServiceResponse struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Path        string    `json:"path" db:"path"`
	Function    string    `json:"function" db:"func"`
	Method      string    `json:"method" db:"method"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateServiceRequest struct {
	Name        string `json:"name" db:"name"`
	Code        string `json:"code" db:"code"`
	Path        string `json:"path" db:"path"`
	Function    string `json:"function" db:"func"`
	Method      string `json:"method" db:"method"`
	IsActive    bool   `json:"is_active" db:"is_active"`
	Description string `json:"description" db:"description"`
	//UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

type InternalUpdateServiceRequest struct {
	ID          int64     `json:"-" db:"-"`
	Name        string    `json:"name" db:"name"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Function    string    `json:"function" db:"func"`
	Description string    `json:"description" db:"description"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

type UpdateServiceRequest struct {
	ID       int64  `json:"-" db:"-"`
	Name     string `json:"name" db:"name"`
	IsActive bool   `json:"is_active" db:"is_active"`
	// Function    string    `json:"-" db:"func"`
	Description string    `json:"description" db:"description"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}
