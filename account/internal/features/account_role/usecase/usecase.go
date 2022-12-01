package usecase

import (
	"context"

	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type usecase struct {
	rpRepo domain.AccountRoleRepository
	logger *log.Logger
}

func New(rpRepo domain.AccountRoleRepository, logger *log.Logger) domain.AccountRoleUsecase {
	return &usecase{
		rpRepo: rpRepo,
		logger: logger,
	}
}

func (u usecase) NewAccountRoles(ctx context.Context, services []entity.AccountRolesResponse) (domain.URInsertResult, error) {
	var items []domain.URItem
	for _, service := range services {
		item := domain.URItem{
			Data: entity.AccountRoleEntity{
				AccountID:   service.AccountID,
				RoleID:      service.RoleID,
				Description: service.Description,
			},
		}

		assignRequest := entity.NewInsertCollectionRequest(domain.AccountRoleTable, service)
		prID, err := u.rpRepo.CreateAccountRoles(ctx, assignRequest)
		if err != nil {
			u.logger.Println("error: ", err)
			item.Status = "failed"
			item.Err = err.Error()
			return domain.URInsertResult{}, err
		}

		item.Data.ID = prID
		item.Status = "success"
		items = append(items, item)
	}

	return domain.URInsertResult{Items: items}, nil
}

func (u usecase) NewAccountRolesN(ctx context.Context, requests entity.CreateRolesToAccountRequest) (domain.URInsertResult, error) {
	var items []domain.URItem
	for _, role := range requests.Roles {
		item := domain.URItem{
			Data: entity.AccountRoleEntity{
				AccountID:   requests.AccountId,
				RoleID:      role,
				Description: requests.Description,
			},
		}

		requests.RoleId = role
		assignRequest := entity.NewInsertCollectionRequest(domain.AccountRoleTable, requests)
		prID, err := u.rpRepo.CreateAccountRoles(ctx, assignRequest)
		if err != nil {
			u.logger.Println("error: ", err)
			item.Status = "failed"
			item.Err = err.Error()
			return domain.URInsertResult{}, err
		}

		item.Data.ID = prID
		item.Status = "success"
		items = append(items, item)
	}

	return domain.URInsertResult{Items: items}, nil
}

func (u usecase) FetchAccountRoles(ctx context.Context, request entity.CollectionRequest) (access []entity.AccountRolesResponse, paginate domain.PaginateApp, err error) {
	request.TableName = domain.AccountRoleTable
	request.Model = entity.AccountRoleModel{}
	total, err := u.rpRepo.TotalAccountRoles(ctx, request)
	//service = make([]entity.Service, total)
	var service entity.ServiceResponse
	if err != nil {
		return access, paginate, err
	}

	paginate.Page = request.Page
	paginate.PerPage = request.PerPage
	paginate.Total = int(total)
	if total > int64(request.Page*request.PerPage) {
		paginate.HasNext = true
	}

	access, err = u.rpRepo.ReadAccountRoles(ctx, request)
	if err != nil {
		return access, paginate, err
	}

	log.Println(service)
	return access, paginate, err
	//return u.rpRepo.ReadAccountRoles(ctx, req)
}

func (u usecase) FetchAccountRolesN(ctx context.Context, request entity.CollectionRequest) (access []entity.AccountRolesResponse, paginate domain.PaginateApp, err error) {
	request.TableName = domain.AccountRoleTable
	request.Model = entity.AccountRoleModel{}
	// todo here we need from frontend to send at least service_id, service_name, role_id, role_name as valid filter
	total, err := u.rpRepo.TotalAccountRoles(ctx, request)
	//service = make([]entity.Service, total)
	if err != nil {
		return access, paginate, err
	}

	paginate.Page = request.Page
	paginate.PerPage = request.PerPage
	paginate.Total = int(total)
	if total > int64(request.Page*request.PerPage) {
		paginate.HasNext = true
	}

	access, err = u.rpRepo.ReadAccountRoles(ctx, request)
	if err != nil {
		return access, paginate, err
	}

	log.Println(access)
	return access, paginate, err
	//return u.rpRepo.ReadAccountRoles(ctx, req)
}

func (u usecase) PatchAccountRolesN(ctx context.Context, requests entity.UpdateRolesToAccountRequest) (domain.URInsertResult, error) {
	// todo : fetch all role of account id
	req := entity.NewCollectionRequest(domain.AccountRoleTable).SetModel(entity.AccountRolesResponse{}).EQ("account_id", requests.AccountId)
	roles, err := u.rpRepo.ReadAccountRoles(ctx, req)
	if err != nil {
		return domain.URInsertResult{}, err
	}

	var exist []int64
	for _, role := range roles {
		exist = append(exist, role.ID)
	}

	var all []int64
	all = append(all, exist...)
	all = append(all, requests.Roles...)
	all = utils.RemoveDuplicateInt(all)

	var items []domain.URItem
	for _, i := range all {
		item := domain.URItem{
			Data: entity.AccountRoleEntity{
				AccountID: requests.AccountId,
				RoleID:    i,
			},
		}

		requests.RoleId = i

		if !utils.IntInSlice(exist, i) && utils.IntInSlice(requests.Roles, i) {
			log.Println("created it ", i)
			assignRequest := entity.NewInsertCollectionRequest(domain.AccountRoleTable, requests)
			prID, err := u.rpRepo.CreateAccountRoles(ctx, assignRequest)
			if err != nil {
				u.logger.Println("error: ", err)
				item.Status = "failed"
				item.Err = err.Error()
				return domain.URInsertResult{}, err
			}

			item.Data.ID = prID
			item.Status = "success"
			items = append(items, item)
		} else if utils.IntInSlice(exist, i) && !utils.IntInSlice(requests.Roles, i) {
			log.Println("deleted ", i)
			deleteRequest := entity.NewDeleteCollectionRequest(domain.AccountRoleTable).
				EQ("account_id", requests.AccountId).
				EQ("role_id", i)
			err = u.rpRepo.DeleteRolesFromAccountN(ctx, deleteRequest)
			if err != nil {
				u.logger.Println("error: ", err)
				item.Status = "failed"
				item.Err = err.Error()
				return domain.URInsertResult{}, err
			}

			item.Data.ID = i
			item.Status = "success"
			items = append(items, item)
		}
	}

	return domain.URInsertResult{Items: items}, nil
	//return domain.URInsertResult{}, nil
}

func (u usecase) AssignRolesToAccount(ctx context.Context, roles entity.CreateRolesToAccountRequest) error {
	//domain.NewCollectionRequest(ctx, roles)
	for _, role := range roles.Roles {
		roles.RoleId = role
		err := u.rpRepo.UpdateAccountRoles(ctx, roles)
		if err != nil {
			return err
		}

	}
	return nil //u.rpRepo.UpdateAccountRoles(ctx, roles)
}

func (u usecase) UnAssignRolesToAccount(ctx context.Context, roles entity.CreateRolesToAccountRequest) error {
	return u.rpRepo.DeleteRolesFromAccount(ctx, roles)
}
