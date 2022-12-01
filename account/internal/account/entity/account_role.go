package entity

import (
	"database/sql"
	"time"
)

type AccountRoleModel struct {
	ID          sql.NullInt64  `db:"id"`
	AccountID   sql.NullInt64  `db:"account_id"`
	RoleID      sql.NullInt64  `db:"role_id"`
	IsActive    sql.NullBool   `db:"is_active"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

type AccountRoleEntity struct {
	ID          int64     `json:"id"`
	AccountID   int64     `json:"account_id"`
	RoleID      int64     `json:"role_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type CreateRolesToAccountRequest struct {
	AccountId   int64   `json:"account_id" db:"account_id"`
	RoleId      int64   `json:"-" db:"role_id"`
	Roles       []int64 `json:"permissions" db:"-"`
	Description string  `json:"description" db:"description"`
}

type UpdateRolesToAccountRequest struct {
	AccountId int64   `json:"-" db:"account_id"`
	RoleId    int64   `json:"-" db:"role_id"`
	Roles     []int64 `json:"permissions" db:"-"`
}

type AccountRolesResponse struct {
	ID          int64     `json:"id" db:"id"`
	AccountID   int64     `json:"account_id" db:"account_id"`
	RoleID      int64     `json:"role_id" db:"role_id"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
