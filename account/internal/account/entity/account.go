package entity

import (
	"database/sql"
	"time"
)

type AccountModel struct {
	ID          sql.NullInt64  `db:"id"`
	UUID        sql.NullString `db:"uuid"`
	FirstName   sql.NullString `db:"first_name"`
	LastName    sql.NullString `db:"last_name"`
	DisplayName sql.NullString `db:"display_name"`
	Password    sql.NullString `db:"password"`
	Email       sql.NullString `db:"email"`
	PhoneNumber sql.NullString `db:"phone_number"`
	Address     sql.NullString `db:"address"`
	ImageUrl    sql.NullString `db:"image_url"`
	TypeId      sql.NullInt64  `db:"type_id"`
	RoleID      sql.NullInt64  `db:"role_id"`
	IsActive    sql.NullBool   `db:"is_active"`
	ExpireAt    sql.NullTime   `db:"expire_at"`
	LastLogin   sql.NullTime   `db:"last_login"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

type AccountEntity struct {
	ID          int64     `json:"id"`
	UUID        string    `json:"uuid"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DisplayName string    `json:"display_name"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number" validate:"required"`
	Address     string    `json:"address"`
	ImageUrl    string    `json:"image_url"`
	Type        string    `json:"type"` // account, sales channel
	TypeId      int64     `json:"type_id"`
	Role        string    `json:"role"` // account, sales channel
	RoleId      int64     `json:"role_id"`
	IsActive    bool      `json:"is_active"`
	ExpireAt    time.Time `json:"expire_at"`
	LastLogin   time.Time `json:"last_login"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type AccountResult struct {
	Accounts    []AccountEntity `json:"accounts"`
	TotalRecord int64           `json:"total_record"`
}

type CreateAccountRequest struct {
	UUID        string    `json:"-"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Password    string    `json:"password" db:"password"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber string    `json:"phone_number" validate:"required" db:"phone_number"`
	Address     string    `json:"address" db:"address"`
	ImageUrl    string    `json:"image_url" db:"image_url"`
	TypeID      int64     `json:"type_id" db:"type_id"`
	RoleID      int64     `json:"role_id" db:"role_id"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	ExpireAt    time.Time `json:"expire_at" db:"expire_at"`
	LastLogin   time.Time `json:"last_login" db:"last_login"`
	Description string    `json:"description" db:"description"`
}

type AccountResponse struct {
	ID          int64     `json:"id"`
	UUID        string    `json:"uuid"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	ImageUrl    string    `json:"image_url"`
	Type        string    `json:"type"` // account, sales channel
	TypeId      int64     `json:"type_id"`
	Role        string    `json:"role"` // account, sales channel
	RoleId      int64     `json:"role_id"`
	IsActive    bool      `json:"is_active"`
	ExpireAt    time.Time `json:"expire_at"`
	LastLogin   time.Time `json:"last_login"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ResetPassRequest struct {
	OldPass string `json:"old_pass" validate:"required"`
	NewPass string `json:"new_pass" validate:"required"`
}

type UpdateAccountRequest struct {
	ID          int64     `json:"-" db:"-"`
	FirstName   string    `json:"first_name" validate:"required" db:"first_name"`
	LastName    string    `json:"last_name" validate:"required" db:"last_name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber string    `json:"phone_number" validate:"required" db:"phone_number"`
	TypeID      int64     `json:"type_id" validate:"required" db:"type_id"`
	RoleID      int64     `json:"role_id" validate:"required" db:"role_id"`
	Address     string    `json:"address" db:"address"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Description string    `json:"description" db:"description"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

type UpdateAccountProfileRequest struct {
	ID        int64  `json:"-" db:"-"`
	FirstName string `json:"first_name" validate:"required" db:"first_name"`
	LastName  string `json:"last_name" validate:"required" db:"last_name"`
	//DisplayName string    `json:"display_name" db:"display_name"`
	Email string `json:"email" db:"email"`
	//PhoneNumber string    `json:"phone_number" validate:"required" db:"phone_number"`
	Address     string    `json:"address" db:"address"`
	Description string    `json:"description" db:"description"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

type UpdateAccountRole struct {
	UUID      string `json:"uuid"`
	OldRoleID int64  `json:"old_role_id"`
	NewRoleID int64  `json:"new_role_id"`
}

type UpdateAccountType struct {
	UUID      string `json:"uuid"`
	OldTypeID int64  `json:"old_type_id"`
	NewTypeID int64  `json:"new_type_id"`
}

func (a AccountModel) ToEntity() AccountEntity {
	return AccountEntity{
		ID:          a.ID.Int64,
		UUID:        a.UUID.String,
		FirstName:   a.FirstName.String,
		LastName:    a.LastName.String,
		DisplayName: a.DisplayName.String,
		Password:    a.Password.String,
		Email:       a.Email.String,
		PhoneNumber: a.PhoneNumber.String,
		Address:     a.Address.String,
		ImageUrl:    a.ImageUrl.String,
		TypeId:      a.TypeId.Int64,
		RoleId:      a.RoleID.Int64,
		IsActive:    a.IsActive.Bool,
		ExpireAt:    a.ExpireAt.Time,
		LastLogin:   a.LastLogin.Time,
		Description: a.Description.String,
		CreatedAt:   a.CreatedAt.Time,
		UpdatedAt:   a.UpdatedAt.Time,
		DeletedAt:   a.DeletedAt.Time,
	}
}

func (a AccountEntity) ToResponse() AccountResponse {
	return AccountResponse{
		ID:          a.ID,
		UUID:        a.UUID,
		FirstName:   a.FirstName,
		LastName:    a.LastName,
		DisplayName: a.DisplayName,
		Email:       a.Email,
		PhoneNumber: a.PhoneNumber,
		Address:     a.Address,
		ImageUrl:    a.ImageUrl,
		Type:        a.Type,
		TypeId:      a.TypeId,
		Role:        a.Role,
		RoleId:      a.RoleId,
		IsActive:    a.IsActive,
		ExpireAt:    a.ExpireAt,
		LastLogin:   a.LastLogin,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func (uar UpdateAccountRequest) ToEntity(entity AccountEntity) AccountEntity {
	entity.FirstName = uar.FirstName
	entity.LastName = uar.LastName
	entity.DisplayName = uar.DisplayName
	entity.Email = uar.Email
	entity.PhoneNumber = uar.PhoneNumber
	entity.TypeId = uar.TypeID
	entity.Address = uar.Address
	entity.IsActive = uar.IsActive
	entity.Description = uar.Description
	entity.UpdatedAt = uar.UpdatedAt
	return entity
}

func (car CreateAccountRequest) ToEntity() AccountEntity {
	return AccountEntity{
		UUID:        car.UUID,
		FirstName:   car.FirstName,
		LastName:    car.LastName,
		DisplayName: car.DisplayName,
		Password:    car.Password,
		Email:       car.Email,
		PhoneNumber: car.PhoneNumber,
		Address:     car.Address,
		TypeId:      car.TypeID,
		IsActive:    car.IsActive,
		ExpireAt:    car.ExpireAt,
		LastLogin:   car.LastLogin,
		Description: car.Description,
	}
}
