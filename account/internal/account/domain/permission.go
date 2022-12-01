package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type PermissionUsecase interface {
	AddPermission(ctx context.Context, account entity.CreatePermissionRequest) (int64, error)
	FetchPermission(ctx context.Context, id []string) ([]entity.PermissionResponse, error)
	FetchPermissions(ctx context.Context) ([]entity.PermissionResponse, error)
	FetchPermissionWithPaginate(ctx context.Context, request entity.FilterSearchRequest) ([]entity.PermissionResponse, PaginateApp, error)
	PatchPermission(ctx context.Context, account entity.PermissionUpdateRequest) error
}

type PermissionRepository interface {
	CreatePermission(ctx context.Context, request entity.CollectionRequest) (int64, error)
	ReadPermission(ctx context.Context, id int64) (entity.PermissionEntity, error)
	NewReadPermission(ctx context.Context, request entity.CollectionRequest, result interface{}) error
	ReadManyPermission(ctx context.Context, request entity.CollectionRequest) ([]entity.PermissionEntity, error)
	ReadPermissionServices(ctx context.Context, request entity.CollectionRequest, permissionID int64) ([]entity.PermissionServicesResponse, error)
	ReadPermissions(ctx context.Context) ([]entity.PermissionEntity, error)
	UpdatePermission(ctx context.Context, account entity.PermissionUpdateRequest) error
	DeletePermission(ctx context.Context, request entity.CollectionRequest) error
}
