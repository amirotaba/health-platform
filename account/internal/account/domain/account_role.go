package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type URItem struct {
	Data   entity.AccountRoleEntity `json:"data"`
	Status string                   `json:"status"`
	Err    string                   `json:"err"`
}
type URInsertResult struct {
	Items []URItem `json:"items"`
}

type AccountRoleUsecase interface {
	NewAccountRoles(ctx context.Context, services []entity.AccountRolesResponse) (URInsertResult, error)
	NewAccountRolesN(ctx context.Context, request entity.CreateRolesToAccountRequest) (URInsertResult, error)
	PatchAccountRolesN(ctx context.Context, request entity.UpdateRolesToAccountRequest) (URInsertResult, error)
	FetchAccountRoles(ctx context.Context, req entity.CollectionRequest) ([]entity.AccountRolesResponse, PaginateApp, error)
	FetchAccountRolesN(ctx context.Context, req entity.CollectionRequest) ([]entity.AccountRolesResponse, PaginateApp, error)
	AssignRolesToAccount(ctx context.Context, request entity.CreateRolesToAccountRequest) error
	UnAssignRolesToAccount(ctx context.Context, request entity.CreateRolesToAccountRequest) error
}

type AccountRoleRepository interface {
	CreateAccountRoles(ctx context.Context, req entity.InsertCollectionRequest) (int64, error)
	ReadAccountRoles(ctx context.Context, req entity.CollectionRequest) ([]entity.AccountRolesResponse, error)
	TotalAccountRoles(ctx context.Context, request entity.CollectionRequest) (int64, error)
	UpdateAccountRolesN(ctx context.Context, request entity.CollectionRequest) error
	DeleteRolesFromAccountN(ctx context.Context, request entity.DeleteCollectionRequest) error
	UpdateAccountRoles(ctx context.Context, request entity.CreateRolesToAccountRequest) error
	DeleteRolesFromAccount(ctx context.Context, request entity.CreateRolesToAccountRequest) error
}
