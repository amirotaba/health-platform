package domain

import (
	"context"

	"git.paygear.ir/giftino/account/internal/account/aggregate"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type ChannelRuleUsecase interface {
	AddChannelRule(ctx context.Context, request entity.CreateChannelRuleRequest) (int64, error)
	FetchChannelRules(ctx context.Context, request entity.GetChannelRuleRequest) ([]entity.ChannelRuleResponse, error)
	FetchChannelRulesCategory(ctx context.Context, request entity.GetChannelRuleRequest) ([]aggregate.CategoryChannelRule, error)
	FetchChannelRule(ctx context.Context, request entity.GetChannelRuleRequest) (entity.ChannelRuleResponse, error)
	PatchChannelRule(ctx context.Context, account entity.ChannelRuleUpdateRequest) error
	RemoveChannelRuleByID(ctx context.Context, id int64) error
}

type ChannelRuleRepository interface {
	CreateChannelRule(ctx context.Context, request entity.InsertCollectionRequest) (int64, error)
	ReadChannelRule(ctx context.Context, id int64) (entity.ChannelRuleEntity, error)
	ReadChannelRules(ctx context.Context, request entity.CollectionRequest) ([]entity.ChannelRuleEntity, error)
	UpdateChannelRule(ctx context.Context, account entity.ChannelRuleUpdateRequest) error
	DeleteChannelRule(ctx context.Context, request entity.DeleteCollectionRequest) error
}
