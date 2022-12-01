package entity

import (
	"database/sql"
	"time"
)

type TokenEntity struct {
	ID           int64     `json:"id" db:"-"`
	UUID         string    `json:"uuid"`
	AccountID    int64     `json:"account_id" db:"account_id"`
	AuthToken    string    `json:"auth_token" db:"auth_token"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	Type         string    `json:"type" db:"type"`
	ExpireAt     time.Time `json:"expire_at" db:"expire_at"`
	CreatedAt    time.Time `json:"created_at" db:"-"`
}

type TokenModel struct {
	ID           sql.NullInt64  `db:"id"`
	UUID         sql.NullString `db:"uuid"`
	AccountID    sql.NullInt64  `db:"account_id"`
	AuthToken    sql.NullString `db:"auth_token"`
	RefreshToken sql.NullString `db:"refresh_token"`
	Type         sql.NullString `db:"type"`
	ExpireAt     sql.NullTime   `db:"expire_at"`
	CreatedAt    sql.NullTime   `db:"created_at"`
}

type CreateTokenRequest struct {
	AccountID     int64     `json:"account_id"`
	RolesID       []int64   `json:"roles"`
	PermissionsID []int64   `json:"permissions"`
	SecretKey     string    `json:"secret_key"`
	ExpireAt      time.Time `json:"expire_at"`
}

type RefreshTokenRequest struct {
	OldAuthToken string `json:"old_auth_token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateTokenRequest struct {
	AuthToken    string `json:"auth_token" db:"auth_token"`
	RefreshToken string `json:"refresh_token" db:"-"`
}
