package entity

import (
	"database/sql"
	"time"
)

type NewOtpRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type OtpRequest struct {
	AccountID        int64     `json:"account_id" db:"account_id"`
	OwnerPhoneNumber string    `json:"phone_number" db:"phone_number"`
	OTP              string    `json:"otp" db:"otp"`
	ExpireAt         time.Time `json:"expire_at" db:"expire_at"`
}

type OtpEntity struct {
	ID          int64     `json:"id"`
	AccountID   int64     `json:"account_id"`
	PhoneNumber string    `json:"phone_number"`
	OTP         string    `json:"otp"`
	ExpireAt    time.Time `json:"expire_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type OtpModel struct {
	ID          sql.NullInt64  `db:"id"`
	AccountID   sql.NullInt64  `db:"account_id"`
	PhoneNumber sql.NullString `db:"phone_number"`
	OTP         sql.NullString `db:"otp"`
	ExpireAt    sql.NullTime   `db:"expire_at"`
	CreatedAt   sql.NullTime   `db:"created_at"`
}

type OtpVerifyRequest struct {
	PhoneNumber string `json:"phone_number"`
	Otp         string `json:"otp"`
}
