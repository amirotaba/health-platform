package usecase

import (
	"context"

	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type usecase struct {
	rpRepo domain.RolePermissionRepository
	logger *log.Logger
}

func New(rpRepo domain.RolePermissionRepository, logger *log.Logger) domain.RolePermissionUsecase {
	return &usecase{
		rpRepo: rpRepo,
		logger: logger,
	}
}

func (u usecase) NewRolePermissions(ctx context.Context, requests entity.CreatePermissionsToRoleRequest) (domain.RPInsertResult, error) {
	var items []domain.RPItem
	for _, permission := range requests.Permissions {
		item := domain.RPItem{
			Data: entity.RolePermissionEntity{
				RoleID:       requests.RoleId,
				PermissionID: permission,
				Description:  requests.Description,
			},
		}

		requests.PermissionId = permission
		assignRequest := entity.NewInsertCollectionRequest(domain.RolePermissionTable, requests)
		prID, err := u.rpRepo.CreateRolePermissions(ctx, assignRequest)
		if err != nil {
			u.logger.Println("error: ", err)
			item.Err = err.Error()
			item.Status = "failed"
			items = append(items, item)
			//return domain.RPInsertResult{}, err
		} else {
			item.Data.ID = prID
			item.Status = "success"
			items = append(items, item)
		}

	}

	return domain.RPInsertResult{Items: items}, nil
}

func (u usecase) FetchRolePermissions(ctx context.Context, req entity.FilterSearchRequest) (access []entity.RolePermissionsResponse, err error) {
	request := entity.NewCollectionRequest(domain.RolePermissionTable).SetModel(entity.RolePermissionModel{})
	request.Filters = req.Filters
	request.Searches = req.Searches
	// todo here we need from frontend to send at least service_id, service_name, permission_id, permission_name as valid filter
	//var service entity.ServiceResponse
	access, err = u.rpRepo.ReadRolePermissions(ctx, request)
	if err != nil {
		return access, err
	}

	//log.Println(service)
	return access, err
	//return u.rpRepo.ReadRolePermissions(ctx, req)
}

func (u usecase) PatchRolePermissions(ctx context.Context, requests entity.UpdateRolePermissionsRequest) (domain.RPInsertResult, error) {
	// todo : fetch all permission of role id
	req := entity.NewCollectionRequest(domain.RolePermissionTable).SetModel(entity.RolePermissionsResponse{}).IN("role_id", requests.RoleId)
	permissions, err := u.rpRepo.ReadRolePermissions(ctx, req)
	if err != nil {
		return domain.RPInsertResult{}, err
	}

	var exist []int64
	for _, permission := range permissions {
		exist = append(exist, permission.PermissionID)
	}

	var all []int64
	all = append(all, exist...)
	all = append(all, requests.Permissions...)
	all = utils.RemoveDuplicateInt(all)

	var items []domain.RPItem
	for _, i := range all {
		item := domain.RPItem{
			Data: entity.RolePermissionEntity{
				RoleID:       requests.RoleId,
				PermissionID: i,
			},
		}

		requests.PermissionId = i
		if !utils.IntInSlice(exist, i) && utils.IntInSlice(requests.Permissions, i) {
			log.Println("created it ", i)
			assignRequest := entity.NewInsertCollectionRequest(domain.RolePermissionTable, requests)
			prID, err := u.rpRepo.CreateRolePermissions(ctx, assignRequest)
			if err != nil {
				u.logger.Println("error: ", err)
				item.Status = "failed"
				item.Err = err.Error()
				items = append(items, item)
				//return domain.RPInsertResult{}, err
			} else {
				item.Data.ID = prID
				item.Status = "success"
				items = append(items, item)
			}

		} else if utils.IntInSlice(exist, i) && !utils.IntInSlice(requests.Permissions, i) {
			log.Println("deleted ", i)
			deleteRequest := entity.NewDeleteCollectionRequest(domain.RolePermissionTable).EQ("role_id", requests.RoleId).EQ("permission_id", i)
			err = u.rpRepo.DeletePermissionsFromRoleN(ctx, deleteRequest)
			if err != nil {
				u.logger.Println("error: ", err)
				item.Status = "failed"
				item.Err = err.Error()
				items = append(items, item)
				//return domain.RPInsertResult{}, err
			} else {
				item.Data.ID = i
				item.Status = "success"
				items = append(items, item)
			}

		}
	}

	return domain.RPInsertResult{Items: items}, nil
}
