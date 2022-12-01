package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type ServiceUsecase interface {
	NewService(ctx context.Context, service entity.CreateServiceRequest) (int64, error)
	UpsertService(ctx context.Context, service entity.CreateServiceRequest) (int64, error)
	FetchServices(ctx context.Context, request entity.FilterSearchRequest) ([]entity.ServiceResponse, PaginateApp, error)
	PatchService(ctx context.Context, request entity.UpdateServiceRequest) error
}

type ServiceRepository interface {
	CreateService(ctx context.Context, request entity.CollectionRequest) (int64, error)
	ReadOneService(ctx context.Context, request entity.CollectionRequest) (entity.ServiceResponse, error)
	ReadManyServices(ctx context.Context, request entity.CollectionRequest) ([]entity.ServiceResponse, error)
	UpdateService(ctx context.Context, request entity.UpdateCollectionRequest) error
	TotalServices(ctx context.Context, request entity.CollectionRequest) (int64, error)
}
