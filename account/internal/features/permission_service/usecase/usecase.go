package usecase

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type usecase struct {
	prRepo domain.PermissionServiceRepository
	logger *log.Logger
}

func New(prRepo domain.PermissionServiceRepository, logger *log.Logger) domain.PermissionServiceUsecase {
	return &usecase{
		prRepo: prRepo,
		logger: logger,
	}
}

func (u usecase) NewPermissionServices(ctx context.Context, services entity.AssignServicesToPermission) (domain.PSInsertResult, error) {
	var items []domain.PSItem
	for i := range services.Services {
		item := domain.PSItem{
			Data: entity.PermissionServicesEntity{
				PermissionID: services.PermissionID,
				ServiceID:    services.Services[i],
				Description:  services.Description,
			},
		}

		checkExistsReq := entity.NewCollectionRequest(domain.PermissionServiceTable).
			SetModel(entity.PermissionServicesResponseModel{}).
			EQ("permission_id", services.PermissionID).
			EQ("service_id", services.Services[i])

		permissionServices, err := u.prRepo.ReadPermissionServices(ctx, checkExistsReq)
		if err != nil {
			u.logger.Println("error: ", err)
			item.Status = "failed"
			item.Err = err.Error()
			items = append(items, item)
			continue
		} else if len(permissionServices) > 0 {
			u.logger.Println("error: ", err)
			item.Status = "failed"
			item.Err = "this relation exists"
			item.Data.ID = permissionServices[0].ID
			items = append(items, item)
			continue
		}

		services.ServiceID = services.Services[i]
		assignRequest := entity.NewInsertCollectionRequest(domain.PermissionServiceTable, services)
		prID, err := u.prRepo.CreatePermissionServices(ctx, assignRequest)
		if err != nil {
			u.logger.Println("error: ", err)
			item.Status = "failed"
			item.Err = err.Error()
			item.Data.ID = prID
			//return domain.PRInsertResult{}, err
		} else {
			item.Data.ID = prID
			item.Status = "success"
		}

		items = append(items, item)
	}

	return domain.PSInsertResult{Items: items}, nil
}

func (u usecase) FetchPermissionServices(ctx context.Context, request entity.FilterSearchRequest) (prServices []entity.PermissionServicesResponse, paginate domain.PaginateApp, err error) {
	req := entity.NewCollectionRequest(domain.PermissionServiceTable).SetModel(entity.PermissionServicesModel{})
	req.Filters = request.Filters
	req.Searches = request.Searches
	// todo here we need from frontend to send at least service_id, service_name, permission_id, permission_name as valid filter
	req.JoinWithOptions = entity.JoinWithOption{
		RootTable: domain.PermissionServiceTable,
		Joins: []entity.Join{
			{
				JoinType:   string(sqlbuilder.LeftJoin),
				LeftTable:  domain.PermissionTable,
				RightTable: domain.PermissionServiceTable,
				ON: entity.ON{
					RightON: "permission_id",
					LeftON:  "id",
				},
			},
			{
				JoinType:   string(sqlbuilder.LeftJoin),
				LeftTable:  domain.ServicesTable,
				RightTable: domain.PermissionServiceTable,
				ON: entity.ON{
					RightON: "service_id",
					LeftON:  "id",
				},
			},
		},
	}

	var service entity.ServiceResponse
	if err != nil {
		return prServices, paginate, err
	}

	prServices, err = u.prRepo.ReadPermissionServices(ctx, req)
	if err != nil {
		return prServices, paginate, err
	}

	//paginate.Total = len(prServices)
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

	log.Println(service)
	return prServices, paginate, err
}

func (u usecase) PatchPermissionServices(ctx context.Context, requests entity.UpdateAssignedServicesToPermission) (domain.PSInsertResult, error) {
	// todo : fetch all service of permission id
	req := entity.NewCollectionRequest(domain.PermissionServiceTable).
		SetModel(entity.PermissionServicesResponse{}).
		EQ("permission_id", requests.PermissionID)

	permissions, err := u.prRepo.ReadPermissionServices(ctx, req)
	if err != nil {
		return domain.PSInsertResult{}, err
	}

	var exist []int64
	for _, permission := range permissions {
		exist = append(exist, permission.ServiceID)
	}

	var all []int64
	all = append(all, exist...)
	all = append(all, requests.Services...)
	all = utils.RemoveDuplicateInt(all)

	var items []domain.PSItem
	for _, i := range all {
		item := domain.PSItem{
			Data: entity.PermissionServicesEntity{
				PermissionID: requests.PermissionID,
				ServiceID:    i,
				Description:  requests.Description,
			},
		}

		requests.ServiceID = i

		if !utils.IntInSlice(exist, i) && utils.IntInSlice(requests.Services, i) {
			log.Println("created it ", i)
			assignRequest := entity.NewInsertCollectionRequest(domain.PermissionServiceTable, requests)
			item.Data.ID, err = u.prRepo.CreatePermissionServices(ctx, assignRequest)
			if err != nil {
				u.logger.Println("error: ", err)
				item.Status = "failed"
				item.Err = err.Error()
				//return err
			}

			//item.Status = "success"
			//items = append(items, item)
		} else if utils.IntInSlice(exist, i) && !utils.IntInSlice(requests.Services, i) {
			log.Println("deleted ", i)
			deleteRequest := entity.NewDeleteCollectionRequest(domain.PermissionServiceTable).
				EQ("permission_id", requests.PermissionID).
				EQ("service_id", i)
			err = u.prRepo.DeletePermissionServices(ctx, deleteRequest)
			if err != nil {
				u.logger.Println("error: ", err)
				item.Status = "failed"
				item.Err = err.Error()
				//return err
			}

			//item.Data.ID = i
			//item.Status = "success"
			items = append(items, item)
		}
	}

	return domain.PSInsertResult{Items: items}, nil
	//return domain.RPInsertResult{}, nil
}
