package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type CAItem struct {
	Data   entity.ChannelAccountEntity `json:"data"`
	Status string                      `json:"status"`
	Err    string                      `json:"err"`
}
type CAInsertResult struct {
	Items []CAItem `json:"items"`
}

type ChannelAccountUsecase interface {
	NewChannelAccounts(ctx context.Context, request entity.AddChannelAccountsRequest) (CAInsertResult, error)
	FetchChannelAccounts(ctx context.Context, req entity.FilterSearchRequest) ([]entity.ChannelAccountsResponse, PaginateApp, error)
	PatchChannelAccounts(ctx context.Context, request entity.UpdateAccountsToChannelRequest) (CAInsertResult, error)
	RemoveChannelAccounts(ctx context.Context, request entity.DeleteAccountsFromChannelRequest) (CAInsertResult, error)
}

type ChannelAccountRepository interface {
	CreateChannelAccounts(ctx context.Context, req entity.InsertCollectionRequest) (int64, error)
	ReadChannelAccounts(ctx context.Context, req entity.CollectionRequest) ([]entity.ChannelAccountEntity, error)
	UpdateChannelAccounts(ctx context.Context, request entity.UpdateCollectionRequest) error
	DeleteAccountsFromChannel(ctx context.Context, request entity.DeleteCollectionRequest) error
}
