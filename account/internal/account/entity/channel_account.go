package entity

import (
	"database/sql"
	"time"
)

type ChannelAccountModel struct {
	ID          sql.NullInt64  `db:"id"`
	ChannelID   sql.NullInt64  `db:"channel_id"`
	AccountID   sql.NullInt64  `db:"account_id"`
	RoleID      sql.NullInt64  `db:"role_id"`
	IsActive    sql.NullBool   `db:"is_active"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

type ChannelAccountEntity struct {
	ID             int64     `json:"id"`
	ChannelID      int64     `json:"channel_id"`
	AccountID      int64     `json:"account_id"`
	RoleID         int64     `json:"role_id"`
	CurrentBalance int64     `json:"current_balance"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}

type ChannelAccountRoles struct {
	AccountID int64 `json:"account_id"`
	RoleID    int64 `json:"role_id"`
}

type AddChannelAccountsRequest struct {
	ChannelId int64 `json:"channel_id" db:"channel_id"`
	AccountID int64 `json:"account_id" db:"account_id"`
	RoleID    int64 `json:"role_id" db:"role_id"`
	//AccountsRoles []ChannelAccountRoles `json:"accounts_roles" db:"-"`
	Description string `json:"description" db:"description"`
}

type AccountRole struct {
	AccountID sql.NullInt64 `db:"account_id"`
	RoleID    sql.NullInt64 `db:"role_id"`
}

type UpdateAccountsToChannelRequest struct {
	RelationID int64 `json:"-" db:"-"`
	//ChannelId   int64   `json:"-" db:"channel_id"`
	AccountID int64 `json:"account_id" db:"account_id"`
	RoleID    int64 `json:"role_id" db:"role_id"`
	//AccountIDs  []int64 `json:"account_ids" db:"-"`
	Description string `json:"description" db:"description"`
}

type DeleteAccountsFromChannelRequest struct {
	RelationID int64 `json:"relation_id" db:"-"`
	//ChannelId   int64   `json:"-" db:"channel_id"`
	AccountID int64 `json:"-" db:"account_id"`
	RoleID    int64 `json:"-" db:"role_id"`
	//AccountIDs  []int64 `json:"account_ids" db:"-"`
	Description string `json:"-" db:"-"`
}

type ChannelAccountsResponse struct {
	ID             int64     `json:"id" db:"id"`
	ChannelID      int64     `json:"channel_id" db:"channel_id"`
	Channel        string    `json:"channel" db:"-"`
	AccountID      int64     `json:"account_id" db:"account_id"`
	Account        string    `json:"account" db:"-"`
	CurrentBalance int64     `json:"current_balance" db:"-"`
	RoleID         int64     `json:"role_id" db:"role_id"`
	Role           string    `json:"role" db:"-"`
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

func (model ChannelAccountModel) ToEntity() ChannelAccountEntity {
	return ChannelAccountEntity{
		ID:          model.ID.Int64,
		ChannelID:   model.ChannelID.Int64,
		AccountID:   model.AccountID.Int64,
		RoleID:      model.RoleID.Int64,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
		DeletedAt:   model.DeletedAt.Time,
	}
}

func (entity ChannelAccountEntity) ToResponse() ChannelAccountsResponse {
	return ChannelAccountsResponse{
		ID:             entity.ID,
		ChannelID:      entity.ChannelID,
		AccountID:      entity.AccountID,
		CurrentBalance: entity.CurrentBalance,
		RoleID:         entity.RoleID,
		Description:    entity.Description,
		CreatedAt:      entity.CreatedAt,
	}
}
