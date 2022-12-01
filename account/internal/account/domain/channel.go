package domain

import (
	"context"

	"github.com/nats-io/nats.go"

	"git.paygear.ir/giftino/account/internal/account/entity"
)

type ChannelUsecase interface {
	AddChannel(ctx context.Context, channel entity.CreateChannelRequest) (int64, error)
	FetchChannelsWithPaginate(ctx context.Context, request entity.FilterSearchRequest) ([]entity.ChannelResponse, PaginateApp, error)
	FetchChannelsWithID(ctx context.Context, id int64) (entity.ChannelResponse, error)
	PatchChannel(ctx context.Context, channel entity.UpdateChannelRequest) error
	ChargeChannel(ctx context.Context, channel entity.ChargeChannelRequest) error
	EventHandler(conn *nats.Conn)
}

type ChannelRepository interface {
	CreateChannel(ctx context.Context, request entity.InsertCollectionRequest) (int64, error)
	ReadOneChannel(ctx context.Context, request entity.CollectionRequest) (entity.ChannelEntity, error)
	ReadManyChannels(ctx context.Context, request entity.CollectionRequest) ([]entity.ChannelEntity, error)
	ReadTotal(ctx context.Context, request entity.CollectionRequest) entity.CollectionRequest
	UpdateChannelNew(ctx context.Context, request entity.UpdateCollectionRequest) error
}
