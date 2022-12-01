package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/account/event"
	pbWallet "git.paygear.ir/giftino/account/internal/account/proto/wallet"
)

type usecase struct {
	repo           domain.ChannelRepository
	accountUsecase domain.AccountUsecase
	channelAccount domain.ChannelAccountRepository
	wallet         domain.WalletGrpcPort
	eventBus       event.Bus
	logger         *log.Logger
}

func New(repo domain.ChannelRepository, accountUsecase domain.AccountUsecase, channelAccount domain.ChannelAccountRepository,
	wallet domain.WalletGrpcPort, eventBus event.Bus, logger *log.Logger) domain.ChannelUsecase {
	u := &usecase{
		repo:           repo,
		accountUsecase: accountUsecase,
		channelAccount: channelAccount,
		wallet:         wallet,
		logger:         logger,
	}

	//u.eventBus.Register("account.channel")
	return u
}

func (u usecase) AddChannel(ctx context.Context, channel entity.CreateChannelRequest) (int64, error) {
	// todo add all of it stuff in trx
	channel.UUID = uuid.New().String()
	wallet, err := u.wallet.CreateWallet(ctx, channel.UUID, channel.CurrentBalance)
	if err != nil {
		return 0, err
	}

	channel.ExpireAt = time.Now().AddDate(0, 1, 0)
	req := entity.NewInsertCollectionRequest(domain.ChannelTable, channel)
	id, err := u.repo.CreateChannel(ctx, req)
	if err != nil {
		return 0, err
	}

	u.logger.Println(wallet.Created)
	u.logger.Println(wallet.StatusName)
	u.logger.Println(wallet.StatusCode)
	return id, nil
}

func (u usecase) FetchChannelsWithPaginate(ctx context.Context, request entity.FilterSearchRequest) ([]entity.ChannelResponse, domain.PaginateApp, error) {
	req := entity.NewCollectionRequest(domain.ChannelTable).SetModel(entity.ChannelModel{}).SQL("count(*)")
	req.Filters = request.Filters
	req.Searches = request.Searches
	req.OrderBy = request.Sort
	var total int64
	u.repo.ReadTotal(ctx, req).Scan(&total)
	req.PaginateData = entity.Paginate{
		Has:    true,
		Limit:  request.PerPage,
		Offset: request.PerPage * (request.Page - 1),
	}

	//joins := entity.JoinWithOption{RootTable: domain.ChannelTable}
	//joins.Joins = append(joins.Joins, entity.Join{
	//	JoinType:   string(sqlbuilder.LeftJoin),
	//	LeftTable:  domain.ChannelAccountTable,
	//	RightTable: domain.ChannelTable,
	//	ON: entity.ON{
	//		RightON: "id",
	//		LeftON:  "channel_id",
	//	},
	//})
	//
	//joins.Joins = append(joins.Joins, entity.Join{
	//	JoinType:   string(sqlbuilder.LeftJoin),
	//	LeftTable:  domain.AccountTable,
	//	RightTable: domain.ChannelAccountTable,
	//	ON: entity.ON{
	//		LeftON:  "id",
	//		RightON: "account_id",
	//	},
	//})

	var resp []entity.ChannelResponse
	var paginate domain.PaginateApp
	//req.JoinWithOptions = joins
	//req.GroupsBy = append(req.GroupsBy, entity.GroupBy{
	//	TableName: domain.AccountTable,
	//	Field:     "id",
	//})

	channels, err := u.repo.ReadManyChannels(ctx, req)
	if err != nil {
		return resp, paginate, err
	}

	paginate.Total = int(total)
	paginate.Page = request.Page
	paginate.PerPage = request.PerPage
	paginate.TotalPage = paginate.Total / paginate.PerPage
	u.logger.Println(paginate.Total)
	u.logger.Println(paginate.TotalPage)
	if paginate.Total%request.PerPage > 0 {
		paginate.TotalPage += 1
	}

	if paginate.TotalPage >= paginate.Page {
		paginate.HasNext = true
	} else {
		paginate.HasNext = false
	}

	for _, channel := range channels {
		ch := channel.ToResponse()
		w, _ := u.wallet.GetWallet(ctx, &pbWallet.GetWalletRequest{
			UserID:   ch.UUID,
			WalletID: ch.WalletID,
		})

		ch.CurrentBalance = w.Balance
		resp = append(resp, ch)
	}
	return resp, paginate, nil
}

func (u usecase) FetchChannelsWithID(ctx context.Context, id int64) (entity.ChannelResponse, error) {
	req := entity.CollectionRequest{}
	req.TableName = domain.ChannelTable
	req.Model = entity.ChannelModel{}
	req.Filters = append(req.Filters, entity.Filter{
		TableName: domain.ChannelTable,
		Field:     "id",
		Type:      entity.EQ,
		Value:     id,
	})
	channel, err := u.repo.ReadOneChannel(ctx, req)
	if err != nil {
		return channel.ToResponse(), err
	}

	w, gwErr := u.wallet.GetWallet(ctx, &pbWallet.GetWalletRequest{
		UserID:   channel.UUID,
		WalletID: channel.WalletID,
	})

	if gwErr == nil {
		channel.CurrentBalance = w.Balance
	}
	return channel.ToResponse(), nil
}

func (u usecase) PatchChannel(ctx context.Context, channel entity.UpdateChannelRequest) error {
	get := entity.NewCollectionRequest(domain.ChannelTable).SetModel(entity.ChannelModel{}).EQ("id", channel.ID)
	_, err := u.repo.ReadOneChannel(ctx, get)
	if err != nil {
		return err
	}

	update := entity.NewUpdateCollectionRequest(domain.ChannelTable).SetModel(channel).EQ("id", channel.ID)
	return u.repo.UpdateChannelNew(ctx, update)
}

func (u usecase) ChargeChannel(ctx context.Context, channelCharge entity.ChargeChannelRequest) error {
	channelsWithID, err := u.FetchChannelsWithID(ctx, channelCharge.ID)
	if err != nil {
		return err
	}

	channelCharge.CurrentBalance = channelsWithID.CurrentBalance + channelCharge.CurrentBalance
	updateReq := entity.NewUpdateCollectionRequest(domain.ChannelTable).SetModel(channelCharge).EQ("id", channelCharge.ID)
	err = u.repo.UpdateChannelNew(ctx, updateReq)
	if err != nil {
		u.logger.Println(err)
		return err
	}

	return nil
}

func (u usecase) EventHandler(conn *nats.Conn) {
	if _, err := conn.Subscribe("account.channel", func(m *nats.Msg) {
		var baseEvent event.BaseEvent
		err := json.Unmarshal(m.Data, &baseEvent)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(m.Data))
		log.Println(baseEvent.Data)
		switch baseEvent.GetEventType() {
		case event.ChannelChargedEventType:
			fmt.Println("**************************ChannelChargedEventType*************************")
			var charge entity.ChargeChannelRequest
			log.Println(json.Unmarshal([]byte(baseEvent.Data), &charge))
			log.Println(charge.ID)
			log.Println(charge.CurrentBalance)
			err = u.ChargeChannel(context.Background(), charge)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}); err != nil {
		return
	}
}
