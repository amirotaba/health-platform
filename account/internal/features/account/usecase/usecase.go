package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/command"
	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/account/event"
	"git.paygear.ir/giftino/account/internal/utils"
)

type usecase struct {
	accountRepo     domain.AccountRepository
	roleRepo        domain.RoleRepository
	permitRepo      domain.PermissionRepository
	accountTypeRepo domain.AccountTypeRepository
	eventStore      event.Store
	eventBus        event.Bus
	CommandChan     chan command.Command
	logger          *log.Logger
}

func New(
	accountRepo domain.AccountRepository, roleRepo domain.RoleRepository,
	permitRepo domain.PermissionRepository,
	accountTypeRepo domain.AccountTypeRepository,
	eventStore event.Store, eventBus event.Bus, logger *log.Logger) domain.AccountUsecase {
	return &usecase{
		accountRepo:     accountRepo,
		roleRepo:        roleRepo,
		permitRepo:      permitRepo,
		accountTypeRepo: accountTypeRepo,
		eventStore:      eventStore,
		eventBus:        eventBus,
		CommandChan:     make(chan command.Command),
		logger:          logger,
	}
}

func (u usecase) AddAccount(ctx context.Context, account entity.CreateAccountRequest) (id int64, err error) {
	phone, err := utils.CheckPhone(account.PhoneNumber)
	if err != nil {
		return
	}

	getAccountReq := entity.NewCollectionRequest(domain.AccountTable).SetModel(entity.AccountModel{}).EQ("phone_number", phone)
	//account, err := u.accountRepo.ReadAccountByOwnerPhoneNumber(ctx, phone)
	_, err = u.accountRepo.ReadOneAccount(ctx, getAccountReq)
	switch err.(type) {
	case utils.NotFoundError:
		if account.RoleID == 0 {
			roleReq := entity.NewCollectionRequest(domain.RoleTable).SetModel(entity.RoleResponse{}).EQ("name", "guest")
			var role entity.RoleEntity
			role, err = u.roleRepo.ReadOneRole(ctx, roleReq)
			if err != nil {
				return 0, err
			}

			account.RoleID = role.ID
		}

		if account.TypeID == 0 {
			//roleReq := domain.NewCollectionRequest(domain.AccountTypeTable, {})
			typesAccount, err := u.accountTypeRepo.ReadAccountTypes(ctx)
			if err != nil {
				return 0, err
			}

			for _, typeEntity := range typesAccount {
				if typeEntity.Name == "guest" {
					account.TypeID = typeEntity.ID
				}
			}
		}

		account.PhoneNumber = phone
		account.UUID = uuid.New().String()
		account.LastLogin = time.Now()
		req := entity.NewInsertCollectionRequest(domain.AccountTable, account)
		return u.accountRepo.CreateAccountNew(ctx, req)
	case utils.MysqlInternalServerError:
		return 0, err
	default:
		return 0, errors.New("user Exist")
	}
}

func (u usecase) FetchAccountByOwnerPhoneNumber(ctx context.Context, phoneNumber string) (entity.AccountResponse, error) {
	getAccountReq := entity.NewCollectionRequest(domain.AccountTable).SetModel(entity.AccountModel{}).EQ("phone_number", phoneNumber)
	account, err := u.accountRepo.ReadOneAccount(ctx, getAccountReq)
	//account, err := u.accountRepo.ReadAccountByOwnerPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return entity.AccountResponse{}, err
	}

	return account.ToResponse(), err
}

func (u usecase) FetchAccountByID(ctx context.Context, id int64) (entity.AccountResponse, error) {
	getAccountReq := entity.NewCollectionRequest(domain.AccountTable).SetModel(entity.AccountModel{}).EQ("id", id)
	account, err := u.accountRepo.ReadOneAccount(ctx, getAccountReq)
	//account, err := u.accountRepo.ReadAccountByOwnerPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return entity.AccountResponse{}, err
	}

	a := account.ToResponse()
	accountType, err := u.accountTypeRepo.ReadAccountType(ctx, a.TypeId)
	if err != nil {
		return entity.AccountResponse{}, err
	}

	a.Type = accountType.Name
	roleReq := entity.NewCollectionRequest(domain.RoleTable).SetModel(entity.RoleResponse{}).EQ("id", a.RoleId)

	var role entity.RoleEntity
	role, err = u.roleRepo.ReadOneRole(ctx, roleReq)
	//if err != nil {
	//	return entity.AccountResponse{}, err
	//}

	a.Role = role.Name
	return a, err
}

func (u usecase) FetchAccounts(ctx context.Context) ([]entity.AccountResponse, error) {
	accounts, err := u.accountRepo.ReadAccounts(ctx)
	var resp []entity.AccountResponse
	if err != nil {
		return resp, err
	}

	for _, account := range accounts {
		resp = append(resp, account.ToResponse())
	}
	return resp, nil
}

func (u usecase) FetchAccountsWithPaginate(ctx context.Context, request entity.FilterSearchRequest) ([]entity.AccountResponse, domain.PaginateApp, error) {
	req := entity.CollectionRequest{Filters: request.Filters, Searches: request.Searches}
	req.TableName = domain.AccountTable
	req.Model = entity.AccountModel{}
	accounts, err := u.accountRepo.ReadManyAccounts(ctx, req)
	var resp []entity.AccountResponse
	var paginate domain.PaginateApp
	if err != nil {
		return resp, paginate, err
	}

	paginate.Total = len(accounts)
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

	for _, account := range accounts {
		a := account.ToResponse()
		accountType, err := u.accountTypeRepo.ReadAccountType(ctx, a.TypeId)
		if err != nil {
			return nil, domain.PaginateApp{}, err
		}

		a.Type = accountType.Name

		roleReq := entity.NewCollectionRequest(domain.RoleTable).SetModel(entity.RoleResponse{}).EQ("id", a.RoleId)
		var role entity.RoleEntity
		role, err = u.roleRepo.ReadOneRole(ctx, roleReq)
		if err != nil {
			return nil, domain.PaginateApp{}, err
		}

		a.Role = role.Name
		resp = append(resp, a)
	}

	return resp, paginate, nil
}

func (u usecase) PatchAccount(ctx context.Context, account entity.UpdateAccountRequest) error {
	//_, err := u.accountRepo.ReadAccount(ctx, account.ID)
	getAccountReq := entity.NewCollectionRequest(domain.AccountTable).SetModel(entity.AccountModel{}).EQ("id", account.ID)
	//account, err := u.accountRepo.ReadAccountByOwnerPhoneNumber(ctx, phone)
	_, err := u.accountRepo.ReadOneAccount(ctx, getAccountReq)
	if err != nil {
		return err
	}

	if account.UpdatedAt.IsZero() {
		account.UpdatedAt = time.Now()
	}
	//return u.accountRepo.UpdateAccount(ctx, account)
	req := entity.NewUpdateCollectionRequest(domain.AccountTable).SetModel(account).EQ("id", account.ID)
	return u.accountRepo.UpdateAccountNew(ctx, req)
}

func (u usecase) PatchAccountProfile(ctx context.Context, account entity.UpdateAccountProfileRequest) error {
	//_, err := u.accountRepo.ReadAccount(ctx, account.ID)
	getAccountReq := entity.NewCollectionRequest(domain.AccountTable).SetModel(entity.AccountModel{}).EQ("id", account.ID)
	//account, err := u.accountRepo.ReadAccountByOwnerPhoneNumber(ctx, phone)
	_, err := u.accountRepo.ReadOneAccount(ctx, getAccountReq)
	if err != nil {
		return err
	}

	if account.UpdatedAt.IsZero() {
		account.UpdatedAt = time.Now()
	}
	//return u.accountRepo.UpdateAccount(ctx, account)
	req := entity.NewUpdateCollectionRequest(domain.AccountTable).SetModel(account).EQ("id", account.ID)
	return u.accountRepo.UpdateAccountNew(ctx, req)
}

//func (u usecase) PatchAccountPass(ctx context.Context, req entity.ResetPassRequest, id int64) error {
//	account, err := u.accountRepo.ReadAccount(ctx, id)
//	if err != nil {
//		return err
//	}
//
//	if account.Password == req.OldPass {
//		return errors.New("invalid password")
//	}
//
//	return u.accountRepo.UpdateAccountPassword(ctx, req, id)
//}
