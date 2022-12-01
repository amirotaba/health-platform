package grpc

import (
	"context"
	"fmt"
	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type server struct {
	usecase domain.ChannelRuleUsecase
}

func (s server) mustEmbedUnimplementedChannelRuleServiceServer() {
	//TODO implement me
	panic("implement me")
}

func New(usecase domain.ChannelRuleUsecase) ChannelRuleServiceServer {
	return server{usecase: usecase}
}

func (s server) GetChannelTags(ctx context.Context, request *ChannelsRuleRequest) (*ChannelRuleReply, error) {
	req := entity.GetChannelRuleRequest{}
	req.Filters = append(req.Filters, entity.Filter{
		Field: "channel_id",
		Type:  entity.EQ,
		Value: request.ChannelID,
	})

	results, err := s.usecase.FetchChannelRules(ctx, req)
	if err != nil {
		return nil, err
	}

	var res []*ChannelRule
	for _, result := range results {
		res = append(res, &ChannelRule{
			TagID:    result.TagId,
			Price:    float32(result.Price),
			IsActive: result.IsActive,
		})
	}

	return &ChannelRuleReply{Rules: res}, err
}

func (s server) GetChannelTag(ctx context.Context, request *ChannelRuleRequest) (*ChannelRule, error) {
	fmt.Println("****************************GetChannelTag*******************************")
	defer fmt.Println("****************************GetChannelTag*******************************")
	fmt.Println(request.ChannelID)
	fmt.Println(request.TagID)
	req := entity.GetChannelRuleRequest{}
	req.Filters = append(req.Filters, entity.Filter{
		Field: "channel_id",
		Type:  entity.EQ,
		Value: request.ChannelID,
	})

	req.Filters = append(req.Filters, entity.Filter{
		Field: "tag_id",
		Type:  entity.EQ,
		Value: request.TagID,
	})
	results, err := s.usecase.FetchChannelRule(ctx, req)
	if err != nil {
		return nil, err
	}

	fmt.Println(results)
	fmt.Println(results.ID)
	fmt.Println(results.ChannelId)
	fmt.Println(results.TagId)
	fmt.Println(results.Price)
	return &ChannelRule{
		ChannelID: results.ChannelId,
		TagID:     results.TagId,
		Price:     float32(results.Price),
		IsActive:  results.IsActive,
	}, err
}
