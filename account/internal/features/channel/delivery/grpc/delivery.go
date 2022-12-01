package grpc

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/domain"
	channelProto "git.paygear.ir/giftino/account/internal/account/proto/channel"
)

type server struct {
	channelProto.UnsafeChannelServiceServer
	usecase domain.ChannelUsecase
}

func New(usecase domain.ChannelUsecase) channelProto.ChannelServiceServer {
	return server{usecase: usecase}
}

func (s server) GetChannel(ctx context.Context, request *channelProto.ChannelRequest) (*channelProto.ChannelReply, error) {
	ch, err := s.usecase.FetchChannelsWithID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &channelProto.ChannelReply{
		ID:             ch.ID,
		UUID:           ch.UUID,
		Name:           ch.DisplayName,
		CurrentBalance: float32(ch.CurrentBalance),
		IsActive:       ch.IsActive,
	}, err
}
