package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type PSItem struct {
	Data   entity.PermissionServicesEntity `json:"data"`
	Status string                          `json:"status"`
	Err    string                          `json:"err"`
}

type PSInsertResult struct {
	Items []PSItem `json:"items"`
}

type PermissionServiceUsecase interface {
	NewPermissionServices(ctx context.Context, services entity.AssignServicesToPermission) (PSInsertResult, error)
	FetchPermissionServices(ctx context.Context, req entity.FilterSearchRequest) ([]entity.PermissionServicesResponse, PaginateApp, error)
	PatchPermissionServices(ctx context.Context, req entity.UpdateAssignedServicesToPermission) (PSInsertResult, error)
}

type PermissionServiceRepository interface {
	CreatePermissionServices(ctx context.Context, req entity.InsertCollectionRequest) (int64, error)
	ReadPermissionServices(ctx context.Context, req entity.CollectionRequest) ([]entity.PermissionServicesResponse, error)
	DeletePermissionServices(ctx context.Context, req entity.DeleteCollectionRequest) error
	TotalPermissionServices(ctx context.Context, request entity.CollectionRequest) (int64, error)
}
