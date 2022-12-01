package entity

import (
	"database/sql"
	"time"
)

type ChannelModel struct {
	ID               sql.NullInt64  `db:"id"`
	UUID             sql.NullString `db:"uuid"`
	Name             sql.NullString `db:"name"`
	DisplayName      sql.NullString `db:"display_name"`
	WalletID         sql.NullString `db:"wallet_id"`
	Password         sql.NullString `db:"password"`
	Email            sql.NullString `db:"email"`
	ImageUrl         sql.NullString `db:"image_url"`
	MembershipTypeId sql.NullInt64  `db:"membership_type_id"`
	CurrentBalance   sql.NullInt64  `db:"current_balance"`
	IsActive         sql.NullBool   `db:"is_active"`
	OwnerPhoneNumber sql.NullString `db:"owner_phone_number"`
	Description      sql.NullString `db:"description"`
	ExpireAt         sql.NullTime   `db:"expire_at"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
	DeletedAt        sql.NullTime   `db:"deleted_at"`
}

type ChannelEntity struct {
	ID               int64     `json:"id"`
	UUID             string    `json:"uuid"`
	Name             string    `json:"name" validate:"required"`
	DisplayName      string    `json:"display_name"`
	WalletID         string    `json:"wallet_id"`
	Password         string    `json:"password"`
	Email            string    `json:"email" validate:"required"`
	ImageUrl         string    `json:"image_url"`
	MembershipTypeId int64     `json:"membership_type_id"`
	MembershipType   string    `json:"membership_type"`
	CurrentBalance   int64     `json:"current_balance"`
	IsActive         bool      `json:"is_active"`
	OwnerPhoneNumber string    `json:"owner_phone_number" validate:"required"`
	Description      string    `json:"description"`
	ExpireAt         time.Time `json:"expire_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type ChannelResponse struct {
	ID               int64     `json:"id"`
	UUID             string    `json:"uuid"`
	Name             string    `json:"name"`
	DisplayName      string    `json:"display_name"`
	WalletID         string    `json:"wallet_id"`
	Email            string    `json:"email"`
	ImageUrl         string    `json:"image_url"`
	MembershipTypeId int64     `json:"membership_type_id"`
	MembershipType   string    `json:"membership_type"`
	CurrentBalance   int64     `json:"current_balance" `
	IsActive         bool      `json:"is_active"`
	OwnerPhoneNumber string    `json:"owner_phone_number"`
	Description      string    `json:"description"`
	ExpireAt         time.Time `json:"expire_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type UpdateChannelRequest struct {
	ID               int64  `json:"-" db:"-"`
	Name             string `json:"name" validate:"required" db:"name"`
	DisplayName      string `json:"display_name" db:"display_name"`
	WalletID         string `json:"wallet_id" db:"wallet_id"`
	Email            string `json:"email" validate:"required" db:"email"`
	MembershipTypeId int64  `json:"membership_type_id" db:"membership_type_id"`
	IsActive         bool   `json:"is_active" db:"is_active"`
	OwnerPhoneNumber string `json:"owner_phone_number" validate:"required" db:"owner_phone_number"`
	Description      string `json:"description" db:"description"`
}

type CreateChannelRequest struct {
	UUID             string    `json:"-" db:"uuid"`
	Name             string    `json:"name" validate:"required" db:"name"`
	DisplayName      string    `json:"display_name" db:"display_name"`
	WalletID         string    `json:"wallet_id" db:"wallet_id"`
	Password         string    `json:"password" db:"password"`
	Email            string    `json:"email" validate:"required" db:"email"`
	MembershipTypeId int64     `json:"membership_type_id" db:"-"`
	MembershipType   string    `json:"membership_type" db:"-"`
	CurrentBalance   int64     `json:"current_balance" db:"current_balance"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	OwnerPhoneNumber string    `json:"owner_phone_number" validate:"required" db:"owner_phone_number"`
	Description      string    `json:"description" db:"description"`
	ExpireAt         time.Time `json:"expire_at" db:"expire_at"`
}

type ChargeChannelRequest struct {
	ID             int64 `json:"channel_id" db:"-"`
	CurrentBalance int64 `json:"amount" db:"current_balance"`
}

func (model ChannelModel) ToEntity() ChannelEntity {
	return ChannelEntity{
		ID:               model.ID.Int64,
		UUID:             model.UUID.String,
		Name:             model.Name.String,
		DisplayName:      model.DisplayName.String,
		WalletID:         model.WalletID.String,
		Password:         model.Password.String,
		Email:            model.Email.String,
		ImageUrl:         model.ImageUrl.String,
		MembershipTypeId: model.MembershipTypeId.Int64,
		CurrentBalance:   model.CurrentBalance.Int64,
		IsActive:         model.IsActive.Bool,
		OwnerPhoneNumber: model.OwnerPhoneNumber.String,
		Description:      model.Description.String,
		ExpireAt:         model.ExpireAt.Time,
		CreatedAt:        model.CreatedAt.Time,
		UpdatedAt:        model.UpdatedAt.Time,
		DeletedAt:        model.DeletedAt.Time,
	}
}

func (entity ChannelEntity) ToResponse() ChannelResponse {
	return ChannelResponse{
		ID:               entity.ID,
		UUID:             entity.UUID,
		Name:             entity.Name,
		DisplayName:      entity.DisplayName,
		WalletID:         entity.WalletID,
		Email:            entity.Email,
		ImageUrl:         entity.ImageUrl,
		MembershipTypeId: entity.MembershipTypeId,
		MembershipType:   entity.MembershipType,
		CurrentBalance:   entity.CurrentBalance,
		IsActive:         entity.IsActive,
		OwnerPhoneNumber: entity.OwnerPhoneNumber,
		Description:      entity.Description,
		ExpireAt:         entity.ExpireAt,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}
