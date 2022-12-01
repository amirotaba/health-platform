package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type RoleUsecase interface {
	AddRole(ctx context.Context, role entity.CreateRoleRequest) (int64, error)
	FetchRole(ctx context.Context, id []string) ([]entity.RoleResponse, error)
	FetchRoles(ctx context.Context) ([]entity.RoleResponse, error)
	FetchRoleByID(ctx context.Context, id int64) (entity.RoleResponse, error)
	FetchRoleWithPaginate(ctx context.Context, request entity.FilterSearchRequest) ([]entity.RoleResponse, PaginateApp, error)
	PatchRole(ctx context.Context, role entity.RoleUpdateRequest) error
}

type RoleRepository interface {
	CreateRole(ctx context.Context, request entity.InsertCollectionRequest) (int64, error)
	ReadRole(ctx context.Context, id int64) (entity.RoleEntity, error)
	ReadRoles(ctx context.Context) ([]entity.RoleEntity, error)
	ReadManyRole(ctx context.Context, request entity.CollectionRequest) ([]entity.RoleEntity, error)
	ReadOneRole(ctx context.Context, req entity.CollectionRequest) (entity.RoleEntity, error)
	UpdateRole(ctx context.Context, role entity.RoleUpdateRequest) error
	DeleteRole(ctx context.Context, request entity.DeleteCollectionRequest) error
	//ReadPermissionsRoles(ctx context.Context, roleId int64) ([]entity.RoleEntity, error)
}
