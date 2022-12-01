package domain

import (
	"context"
	"database/sql"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"time"
)

type AccountTypeModel struct {
	ID          sql.NullInt64  `db:"id"`
	Name        sql.NullString `db:"name"`
	IsActive    sql.NullBool   `db:"is_active"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

type AccountTypeEntity struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type AccountTypeResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AccountTypeUpdateRequest struct {
	ID          int64  `json:"-"`
	Name        string `json:"name" validate:"required"`
	IsActive    bool   `json:"is_active"`
	Description string `json:"description"`
}

type AccountTypeUsecase interface {
	AddAccountType(ctx context.Context, accountType *AccountTypeEntity) error
	FetchAccountType(ctx context.Context, ids []string) ([]AccountTypeResponse, error)
	FetchAccountTypes(ctx context.Context) ([]AccountTypeResponse, error)
	PatchAccountType(ctx context.Context, role *AccountTypeUpdateRequest) error
}

type AccountTypeRepository interface {
	CreateAccountType(ctx context.Context, role *AccountTypeEntity) error
	ReadAccountType(ctx context.Context, id int64) (AccountTypeEntity, error)
	ReadAccountTypes(ctx context.Context) ([]AccountTypeEntity, error)
	ReadOneAccountType(ctx context.Context, req entity.CollectionRequest) (AccountTypeEntity, error)
	UpdateAccountType(ctx context.Context, role *AccountTypeUpdateRequest) error
}
