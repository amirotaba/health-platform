package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/reactivex/rxgo/v2"

	"git.paygear.ir/giftino/account/internal/account/aggregate"
	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type usecase struct {
	repo        domain.ChannelRuleRepository
	channelRepo domain.ChannelRepository
	caRepo      domain.ChannelAccountRepository
	inventory   domain.InventoryGrpcPort
}

func New(repo domain.ChannelRuleRepository, inventory domain.InventoryGrpcPort, channelRepo domain.ChannelRepository) domain.ChannelRuleUsecase {
	return &usecase{
		repo:        repo,
		channelRepo: channelRepo,
		inventory:   inventory,
	}
}

func (u usecase) AddChannelRule(ctx context.Context, channelRule entity.CreateChannelRuleRequest) (int64, error) {
	tag, err := u.inventory.TagExist(ctx, channelRule.TagId)
	if err != nil {
		return 0, err
	}

	if !tag.Exist {
		return 0, errors.New("tag inactive")
	}

	req := entity.NewInsertCollectionRequest(domain.ChannelRuleTable, channelRule)
	return u.repo.CreateChannelRule(ctx, req)
}

func (u usecase) FetchChannelRules(ctx context.Context, request entity.GetChannelRuleRequest) ([]entity.ChannelRuleResponse, error) {
	if request.CheckAdmin {
		checkReq := entity.NewCollectionRequest(domain.ChannelAccountTable).SetModel(entity.ChannelAccountsResponse{}).
			EQ("account_id", request.AccountID).
			EQ("channel_id", request.ChannelID)
		_, _ = u.caRepo.ReadChannelAccounts(ctx, checkReq)
	}

	var req entity.CollectionRequest
	req.TableName = domain.ChannelRuleTable
	req.Model = entity.ChannelRuleResponse{}
	req.Filters = request.Filters
	req.PaginateData = request.Pagination
	log.Println(req)
	channelRules, err := u.repo.ReadChannelRules(ctx, req)
	var resp []entity.ChannelRuleResponse
	if err != nil {
		return resp, err
	}

	log.Println(channelRules)
	for _, channelRule := range channelRules {
		tag, err := u.inventory.TagExist(ctx, channelRule.TagId)
		if err != nil {
			return resp, err
		}

		getChannelReq := entity.NewCollectionRequest(domain.ChannelTable).
			SetModel(entity.ChannelModel{}).
			EQ("id", channelRule.ChannelId)
		channel, _ := u.channelRepo.ReadOneChannel(ctx, getChannelReq)
		channelRule.Channel = channel.Name
		channelRule.Tag = tag.Name
		resp = append(resp, newResponse(channelRule))
	}

	return resp, nil
}

func (u usecase) FetchChannelRulesCategory(ctx context.Context, request entity.GetChannelRuleRequest) ([]aggregate.CategoryChannelRule, error) {
	if request.CheckAdmin {
		checkReq := entity.NewCollectionRequest(domain.ChannelAccountTable).SetModel(entity.ChannelAccountsResponse{}).
			EQ("account_id", request.AccountID).
			EQ("channel_id", request.ChannelID)
		_, _ = u.caRepo.ReadChannelAccounts(ctx, checkReq)
	}

	var req entity.CollectionRequest
	req.TableName = domain.ChannelRuleTable
	req.Model = entity.ChannelRuleResponse{}
	req.Filters = request.Filters
	req.PaginateData = request.Pagination
	log.Println(req)
	channelRules, err := u.repo.ReadChannelRules(ctx, req)
	var resp []aggregate.CategoryChannelRule
	if err != nil {
		return resp, err
	}

	log.Println(channelRules)
	/*
		for _, channelRule := range channelRules {
			tag, err := u.inventory.TagExist(ctx, channelRule.TagId)
			if err != nil {
				return resp, err
			}

			getChannelReq := entity.NewCollectionRequest(domain.ChannelTable).
				SetModel(entity.ChannelModel{}).
				EQ("id", channelRule.ChannelId)
			channel, _ := u.channelRepo.ReadOneChannel(ctx, getChannelReq)
			channelRule.Channel = channel.Name
			channelRule.Tag = tag.Name

			resp = append(resp, )
		}

	*/

	err = rxgo.Just(channelRules)().Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		channelRule := i.(entity.ChannelRuleEntity)
		tag, err := u.inventory.TagExist(ctx, channelRule.TagId)
		if err != nil {
			return nil, err
		}

		log.Println(tag)
		log.Println(tag.CategoryID)
		log.Println(tag.CategoryName)
		log.Println(tag)
		log.Println(channelRule)
		log.Println(reflect.TypeOf(i))
		getChannelReq := entity.NewCollectionRequest(domain.ChannelTable).
			SetModel(entity.ChannelModel{}).
			EQ("id", channelRule.ChannelId)
		channel, _ := u.channelRepo.ReadOneChannel(ctx, getChannelReq)
		channelRule.Channel = channel.Name
		channelRule.Tag = tag.Name
		log.Println("1111111111111111111111111111111111111111111111111111111111111111111")
		log.Println(resp)
		for index, rule := range resp {
			if rule.ChannelRuleCategory.CategoryID == tag.CategoryID {
				resp[index].ChannelRules = append(resp[index].ChannelRules, channelRule)
				return nil, nil
			}
		}

		log.Println("22222222222222222222222222222222222222222222222222222222222222222")
		resp = append(resp, aggregate.CategoryChannelRule{
			ChannelRuleCategory: entity.ChannelRuleCategory{
				CategoryID:   tag.CategoryID,
				CategoryName: tag.CategoryName,
			},
			ChannelRules: []entity.ChannelRuleEntity{channelRule},
		})

		log.Println("333333333333333333333333333333333333333333333333333333333333333333333")
		log.Println(resp)
		return nil, nil
	}).Error()

	log.Println("4444444444444444444444444444444444444444444444444444444444444444444444444")
	log.Println(resp)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func (u usecase) FetchChannelRule(ctx context.Context, request entity.GetChannelRuleRequest) (entity.ChannelRuleResponse, error) {
	fmt.Println("************************FetchChannelRule**********************")
	defer fmt.Println("************************FetchChannelRule**********************")
	var req entity.CollectionRequest
	req.TableName = domain.ChannelRuleTable
	req.Model = entity.ChannelRuleResponse{}
	req.Filters = request.Filters
	req.PaginateData = request.Pagination
	log.Println(req)
	channelRules, err := u.repo.ReadChannelRules(ctx, req)
	if err != nil {
		return entity.ChannelRuleResponse{}, err
	}

	if len(channelRules) == 0 {
		return entity.ChannelRuleResponse{}, err
	}

	return newResponse(channelRules[0]), nil
}

func (u usecase) PatchChannelRule(ctx context.Context, channelRule entity.ChannelRuleUpdateRequest) error {
	_, err := u.repo.ReadChannelRule(ctx, channelRule.ID)
	if err != nil {
		return err
	}

	// todo check new tag in update exist
	tag, err := u.inventory.TagExist(ctx, channelRule.TagId)
	if !tag.Exist || err != nil {
		return err
	}

	return u.repo.UpdateChannelRule(ctx, channelRule)
}

func (u usecase) RemoveChannelRuleByID(ctx context.Context, id int64) error {
	req := entity.NewDeleteCollectionRequest(domain.ChannelRuleTable).EQ("id", id)
	return u.repo.DeleteChannelRule(ctx, req)
}

func newResponse(cr entity.ChannelRuleEntity) entity.ChannelRuleResponse {
	return entity.ChannelRuleResponse{
		ID:          cr.ID,
		Channel:     cr.Channel,
		ChannelId:   cr.ChannelId,
		Tag:         cr.Tag,
		TagId:       cr.TagId,
		Price:       cr.Price,
		IsActive:    cr.IsActive,
		Description: cr.Description,
		CreatedAt:   cr.CreatedAt,
		UpdatedAt:   cr.UpdatedAt,
	}
}
