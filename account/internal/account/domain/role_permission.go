package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type RPItem struct {
	Data   entity.RolePermissionEntity `json:"data"`
	Status string                      `json:"status"`
	Err    string                      `json:"err"`
}
type RPInsertResult struct {
	Items []RPItem `json:"items"`
}

type RolePermissionUsecase interface {
	NewRolePermissions(ctx context.Context, request entity.CreatePermissionsToRoleRequest) (RPInsertResult, error)
	FetchRolePermissions(ctx context.Context, req entity.FilterSearchRequest) ([]entity.RolePermissionsResponse, error)
	PatchRolePermissions(ctx context.Context, request entity.UpdateRolePermissionsRequest) (RPInsertResult, error)
}

type RolePermissionRepository interface {
	CreateRolePermissions(ctx context.Context, req entity.InsertCollectionRequest) (int64, error)
	ReadRolePermissions(ctx context.Context, req entity.CollectionRequest) ([]entity.RolePermissionsResponse, error)
	UpdateRolePermissions(ctx context.Context, request entity.UpdateCollectionRequest) error
	DeletePermissionsFromRoleN(ctx context.Context, request entity.DeleteCollectionRequest) error
}
