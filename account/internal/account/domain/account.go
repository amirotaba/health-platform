package domain

import (
	"context"
	"time"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type AccountUsecase interface {
	AddAccount(ctx context.Context, account entity.CreateAccountRequest) (int64, error)
	FetchAccountByOwnerPhoneNumber(ctx context.Context, phoneNumber string) (entity.AccountResponse, error)
	FetchAccountByID(ctx context.Context, id int64) (entity.AccountResponse, error)
	FetchAccounts(ctx context.Context) ([]entity.AccountResponse, error)
	FetchAccountsWithPaginate(ctx context.Context, request entity.FilterSearchRequest) ([]entity.AccountResponse, PaginateApp, error)
	PatchAccount(ctx context.Context, account entity.UpdateAccountRequest) error
	PatchAccountProfile(ctx context.Context, account entity.UpdateAccountProfileRequest) error
}

type AccountRepository interface {
	CreateAccountNew(ctx context.Context, request entity.InsertCollectionRequest) (int64, error)
	ReadOneAccount(ctx context.Context, request entity.CollectionRequest) (entity.AccountEntity, error)
	ReadAccounts(ctx context.Context) ([]entity.AccountEntity, error)
	ReadManyAccounts(ctx context.Context, request entity.CollectionRequest) ([]entity.AccountEntity, error)
	UpdateAccount(ctx context.Context, account entity.UpdateAccountRequest) error
	UpdateAccountNew(ctx context.Context, request entity.UpdateCollectionRequest) error
	UpdateAccountPassword(ctx context.Context, account entity.ResetPassRequest, id int64) error
	UpdateAccountLastLogin(ctx context.Context, phoneNumber string, lastLogin time.Time) error
}
