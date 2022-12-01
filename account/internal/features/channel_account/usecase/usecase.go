package usecase

import (
	"context"

	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type usecase struct {
	caRepo         domain.ChannelAccountRepository
	channelRepo    domain.ChannelRepository
	channelUsecase domain.ChannelUsecase
	accountRepo    domain.AccountRepository
	roleUsecase    domain.RoleUsecase
	logger         *log.Logger
}

func New(
	rpRepo domain.ChannelAccountRepository, channelRepo domain.ChannelRepository,
	accountRepo domain.AccountRepository, channelUsecase domain.ChannelUsecase,
	roleUsecase domain.RoleUsecase, logger *log.Logger) domain.ChannelAccountUsecase {
	return &usecase{
		caRepo:         rpRepo,
		logger:         logger,
		channelRepo:    channelRepo,
		channelUsecase: channelUsecase,
		accountRepo:    accountRepo,
		roleUsecase:    roleUsecase,
	}
}

func (u usecase) NewChannelAccounts(ctx context.Context, requests entity.AddChannelAccountsRequest) (domain.CAInsertResult, error) {
	var items []domain.CAItem
	//for _, accountRoles := range requests.AccountsRoles {
	item := domain.CAItem{
		Data: entity.ChannelAccountEntity{
			ChannelID:   requests.ChannelId,
			AccountID:   requests.AccountID,
			RoleID:      requests.RoleID,
			Description: requests.Description,
		},
	}

	assignRequest := entity.NewInsertCollectionRequest(domain.ChannelAccountTable, requests)
	prID, err := u.caRepo.CreateChannelAccounts(ctx, assignRequest)
	if err != nil {
		u.logger.Println("error: ", err)
		item.Status = "failed"
		item.Err = err.Error()
		items = append(items, item)
		return domain.CAInsertResult{Items: items}, nil
	} else {
		item.Data.ID = prID
		item.Status = "success"
		items = append(items, item)
		return domain.CAInsertResult{Items: items}, nil
	}

}

func (u usecase) FetchChannelAccounts(ctx context.Context, request entity.FilterSearchRequest) (access []entity.ChannelAccountsResponse, paginate domain.PaginateApp, err error) {
	//request.TableName = domain.ChannelAccountTable
	//request.Model = entity.ChannelAccountsResponse{}
	req := entity.NewCollectionRequest(domain.ChannelAccountTable).SetModel(entity.ChannelAccountModel{})
	req.Filters = request.Filters
	req.Searches = request.Searches
	var resp []entity.ChannelAccountsResponse
	channelAccount, err := u.caRepo.ReadChannelAccounts(ctx, req)
	if err != nil {
		return access, paginate, err
	}

	for _, response := range channelAccount {
		ca := response.ToResponse()
		getAccountReq := entity.NewCollectionRequest(domain.AccountTable).SetModel(entity.AccountModel{}).EQ("id", response.AccountID)
		account, err := u.accountRepo.ReadOneAccount(ctx, getAccountReq)
		if err != nil {
			return nil, domain.PaginateApp{}, err
		}
		u.logger.Println(account)
		ca.Account = account.FirstName + account.LastName
		u.logger.Println(account.FirstName + account.LastName)
		//getChannelReq := entity.NewCollectionRequest(domain.ChannelTable).SetModel(entity.ChannelModel{}).EQ("id", response.ChannelID)
		//channel, _ := u.channelRepo.ReadOneChannel(ctx, getChannelReq)
		channel, _ := u.channelUsecase.FetchChannelsWithID(ctx, response.ChannelID)
		u.logger.Println(channel)
		ca.Channel = channel.DisplayName
		ca.CurrentBalance = channel.CurrentBalance
		u.logger.Println(channel.DisplayName)
		role, err := u.roleUsecase.FetchRoleByID(ctx, ca.RoleID)
		if err != nil {
			return nil, domain.PaginateApp{}, err
		}

		ca.Role = role.Name
		resp = append(resp, ca)
	}
	//
	//paginate.Total = len(channelAccount)
	//paginate.Page = request.Page
	//paginate.PerPage = request.PerPage
	//paginate.TotalPage = paginate.Total / paginate.PerPage
	//u.logger.Println(paginate.Total)
	//u.logger.Println(paginate.TotalPage)
	//if paginate.Total%request.PerPage > 0 {
	//	paginate.TotalPage += 1
	//}
	//
	//if paginate.TotalPage >= paginate.Page {
	//	paginate.HasNext = true
	//} else {
	//	paginate.HasNext = false
	//}

	log.Println(resp)
	return resp, paginate, err
}

func (u usecase) PatchChannelAccounts(ctx context.Context, requests entity.UpdateAccountsToChannelRequest) (domain.CAInsertResult, error) {
	req := entity.NewCollectionRequest(domain.ChannelAccountTable).SetModel(requests).EQ("id", requests.RelationID)
	accounts, err := u.caRepo.ReadChannelAccounts(ctx, req)
	if err != nil {
		log.Println(err)
		return domain.CAInsertResult{}, err
	}

	u.logger.Println(accounts)
	//var exist []int64
	//for _, account := range accounts {
	//	exist = append(exist, account.ID)
	//}
	//
	//var all []int64
	//all = append(all, exist...)
	//all = append(all, requests.AccountIDs...)
	//all = utils.RemoveDuplicateInt(all)
	//
	var items []domain.CAItem
	//u.logger.Println(all)
	//for _, i := range all {
	//	item := domain.CAItem{
	//		Data: entity.ChannelAccountEntity{
	//			ChannelID: requests.ChannelId,
	//			AccountID: i,
	//		},
	//	}

	//	requests.AccountId = i
	//	if !utils.IntInSlice(exist, i) && utils.IntInSlice(requests.AccountIDs, i) {
	//		log.Println("created it ", i)
	//		assignRequest := entity.NewInsertCollectionRequest(domain.ChannelAccountTable, requests)
	//		prID, err := u.caRepo.CreateChannelAccounts(ctx, assignRequest)
	//		if err != nil {
	//			u.logger.Println("error: ", err)
	//			item.Status = "failed"
	//			item.Err = err.Error()
	//			//return domain.CAInsertResult{}, err
	//			items = append(items, item)
	//			continue
	//		}
	//
	//		item.Data.ID = prID
	//		item.Status = "success"
	//		items = append(items, item)
	//	} else if utils.IntInSlice(exist, i) && !utils.IntInSlice(requests.AccountIDs, i) {
	//		log.Println("deleted ", i)
	//		deleteRequest := entity.NewDeleteCollectionRequest(domain.ChannelAccountTable)
	//		deleteRequest.Filters = deleteRequest.AddFilter("role_id", requests.ChannelId)
	//		deleteRequest.Filters = deleteRequest.AddFilter("account_id", i)
	//		err = u.caRepo.DeleteAccountsFromChannelN(ctx, deleteRequest)
	//		if err != nil {
	//			u.logger.Println("error: ", err)
	//			item.Status = "failed"
	//			item.Err = err.Error()
	//			items = append(items, item)
	//			continue
	//			//return domain.CAInsertResult{}, err
	//		}
	//
	//		item.Data.ID = i
	//		item.Status = "success"
	//		items = append(items, item)
	//		continue
	//	}
	//}
	//
	//return domain.CAInsertResult{Items: items}, nil
	update := entity.NewUpdateCollectionRequest(domain.ChannelAccountTable).SetModel(requests).EQ("id", requests.RelationID)
	err = u.caRepo.UpdateChannelAccounts(ctx, update)
	item := domain.CAItem{
		Data: entity.ChannelAccountEntity{
			ChannelID: requests.AccountID,
			AccountID: requests.RoleID,
		},
	}

	if err != nil {
		u.logger.Println("error: ", err)
		item.Status = "failed"
		item.Err = err.Error()
		items = append(items, item)
		return domain.CAInsertResult{Items: items}, err
	} else {
		item.Data.ID = requests.RelationID
		item.Status = "success"
		items = append(items, item)
		return domain.CAInsertResult{Items: items}, err
	}
}

func (u usecase) RemoveChannelAccounts(ctx context.Context, requests entity.DeleteAccountsFromChannelRequest) (domain.CAInsertResult, error) {
	req := entity.NewCollectionRequest(domain.ChannelAccountTable).SetModel(requests).EQ("id", requests.RelationID)
	accounts, err := u.caRepo.ReadChannelAccounts(ctx, req)
	if err != nil {
		log.Println(err)
		return domain.CAInsertResult{}, err
	}

	u.logger.Println(accounts)
	var items []domain.CAItem
	deleteReq := entity.NewDeleteCollectionRequest(domain.ChannelAccountTable)
	err = u.caRepo.DeleteAccountsFromChannel(ctx, deleteReq)
	item := domain.CAItem{
		Data: entity.ChannelAccountEntity{
			ChannelID: requests.AccountID,
			AccountID: requests.RoleID,
		},
	}

	if err != nil {
		u.logger.Println("error: ", err)
		item.Status = "failed"
		item.Err = err.Error()
		items = append(items, item)
		return domain.CAInsertResult{Items: items}, err
	} else {
		item.Data.ID = requests.RelationID
		item.Status = "success delete"
		items = append(items, item)
		return domain.CAInsertResult{Items: items}, err
	}
}
