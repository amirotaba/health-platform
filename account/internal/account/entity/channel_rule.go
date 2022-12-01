package entity

import (
	"database/sql"
	"time"
)

type ChannelRuleModel struct {
	ID          sql.NullInt64   `db:"id"`
	ChannelId   sql.NullInt64   `db:"channel_id"`
	TagId       sql.NullInt64   `db:"tag_id"`
	Price       sql.NullFloat64 `db:"price"`
	To          sql.NullString  `db:"destination"`
	IsActive    sql.NullBool    `db:"is_active"`
	Description sql.NullString  `db:"description"`
	CreatedAt   sql.NullTime    `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
	DeletedAt   sql.NullTime    `db:"deleted_at"`
}

type ChannelRuleEntity struct {
	ID          int64     `json:"id"`
	Channel     string    `json:"channel"`
	ChannelId   int64     `json:"channel_id" validate:"required"`
	Tag         string    `json:"product"`
	TagId       int64     `json:"tag_id" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	To          string    `json:"destination"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type ChannelRuleResponse struct {
	ID          int64     `json:"id" db:"id"`
	Channel     string    `json:"channel" db:"-"`
	ChannelId   int64     `json:"channel_id" db:"channel_id"`
	Tag         string    `json:"tag" db:"-"`
	TagId       int64     `json:"tag_id" db:"tag_id"`
	Price       float64   `json:"price" db:"price"`
	To          string    `json:"to" db:"destination"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ChannelRuleCategory struct {
	CategoryID   int64  `json:"category_id" db:"-"`
	CategoryName string `json:"category_name" db:"-"`
}

type ChannelRuleUpdateRequest struct {
	ID          int64   `json:"-"`
	ChannelId   int64   `json:"channel_id" validate:"required"`
	TagId       int64   `json:"tag_id" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	To          string  `json:"to"`
	IsActive    bool    `json:"is_active"`
	Description string  `json:"description"`
}

type CreateChannelRuleRequest struct {
	ChannelId   int64   `json:"channel_id" validate:"required" db:"channel_id"`
	TagId       int64   `json:"tag_id" validate:"required" db:"tag_id"`
	Price       float64 `json:"price" validate:"required" db:"price"`
	To          string  `json:"to" db:"destination"`
	IsActive    bool    `json:"is_active" db:"is_active"`
	Description string  `json:"description" db:"description"`
}

type GetChannelRuleRequest struct {
	ChannelID       int64    `json:"channel_id"`
	ChannelIDs      []int64  `json:"channel_ids"`
	ChannelRulesIDs []int64  `json:"channel_rules_ids"`
	ChannelRulesID  int64    `json:"channel_rules_id"`
	AccountID       int64    `json:"account_id"`
	CheckAdmin      bool     `json:"check_admin"`
	Pagination      Paginate `json:"pagination"`
	Filters         []Filter `json:"filters"`
}

func (model ChannelRuleModel) ToEntity() ChannelRuleEntity {
	return ChannelRuleEntity{
		ID:          model.ID.Int64,
		ChannelId:   model.ChannelId.Int64,
		TagId:       model.TagId.Int64,
		Price:       model.Price.Float64,
		To:          model.To.String,
		IsActive:    model.IsActive.Bool,
		Description: model.Description.String,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
		DeletedAt:   model.DeletedAt.Time,
	}
}

func (entity ChannelRuleEntity) ToResponse() ChannelRuleResponse {
	return ChannelRuleResponse{
		ID:          entity.ID,
		ChannelId:   entity.ChannelId,
		TagId:       entity.TagId,
		Price:       entity.Price,
		To:          entity.To,
		IsActive:    entity.IsActive,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
