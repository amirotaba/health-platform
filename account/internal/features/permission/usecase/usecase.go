package usecase

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
)

type usecase struct {
	permissionRepo domain.PermissionRepository
	serviceRepo    domain.ServiceRepository
	prRepo         domain.PermissionServiceRepository
	logger         *log.Logger
}

func New(repo domain.PermissionRepository, serviceRepo domain.ServiceRepository, prRepo domain.PermissionServiceRepository, logger *log.Logger) domain.PermissionUsecase {
	return &usecase{permissionRepo: repo, serviceRepo: serviceRepo, prRepo: prRepo, logger: logger}
}

func (u usecase) AddPermission(ctx context.Context, request entity.CreatePermissionRequest) (id int64, err error) {
	req := entity.NewCollectionRequest(domain.PermissionTable).SetModel(request)
	id, err = u.permissionRepo.CreatePermission(ctx, req)
	if err != nil {
		return
	}

	//services := request.ServicesID
	//var serviceID int64
	//for i := 0; i < len(services); i++ {
	//	assignRequest := entity.NewCollectionRequest(domain.PermissionServiceTable, entity.AssignServicesToPermission{PermissionID: id, ServiceID: services[i]})
	//	serviceID, err = u.permissionRepo.CreatePermission(ctx, assignRequest)
	//	// do roll back when transaction failed
	//	if err != nil {
	//		deleteRequest := entity.CollectionRequest{TableName: domain.PermissionServiceTable}
	//		deleteRequest.Filters = append(deleteRequest.Filters, entity.Filter{Field: "permission_id", Value: []interface{}{id}})
	//		err = u.permissionRepo.DeletePermission(ctx, deleteRequest)
	//		if err != nil {
	//			return 0, err
	//		}
	//
	//		u.logger.Println("error: ", err)
	//		return
	//	}
	//
	//	log.Printf("a service with id %d created \n", serviceID)
	//}
	return
}

func (u usecase) FetchPermission(ctx context.Context, ids []string) ([]entity.PermissionResponse, error) {
	var permissions []entity.PermissionResponse
	for _, id := range ids {
		permissionId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}

		permission, err := u.permissionRepo.ReadPermission(ctx, permissionId)
		if err != nil {
			return nil, err
		}

		var services []entity.ServiceResponse
		var permissionServices []entity.PermissionServicesResponse
		req := entity.NewCollectionRequest(domain.PermissionServiceTable).SetModel(entity.PermissionServicesResponse{}).EQ("permission_id", permission.ID)
		permissionServices, err = u.prRepo.ReadPermissionServices(ctx, req)
		if err != nil {
			return nil, err
		}

		for _, response := range permissionServices {
			var service entity.ServiceResponse
			var getServiceRequest entity.CollectionRequest
			//(domain.ServicesTable, entity.ServiceResponse{})
			getServiceRequest.TableName = domain.ServicesTable
			getServiceRequest.Model = entity.ServiceResponse{}
			getServiceRequest.Filters = append(getServiceRequest.Filters, entity.Filter{Field: "id", Value: []interface{}{response.ServiceID}})
			service, err = u.serviceRepo.ReadOneService(ctx, getServiceRequest)
			if err != nil {
				return nil, err
			}

			services = append(services, service)
		}
		p := newResponse(permission)
		p.Services = services
		permissions = append(permissions, p)
	}

	return permissions, nil
}

func (u usecase) FetchPermissions(ctx context.Context) ([]entity.PermissionResponse, error) {
	permissions, err := u.permissionRepo.ReadPermissions(ctx)
	var resp []entity.PermissionResponse
	if err != nil {
		return resp, err
	}

	for _, permission := range permissions {
		var services []entity.ServiceResponse
		var permissionServices []entity.PermissionServicesResponse
		req := entity.NewCollectionRequest(domain.PermissionServiceTable).
			SetModel(entity.PermissionServicesResponse{}).
			EQ("permission_id", permission.ID)
		permissionServices, err = u.prRepo.ReadPermissionServices(ctx, req)
		if err != nil {
			return nil, err
		}

		for _, response := range permissionServices {
			var service entity.ServiceResponse
			getServiceRequest := entity.NewCollectionRequest(domain.ServicesTable).
				SetModel(entity.ServiceResponse{}).
				EQ("id", response.ServiceID)
			service, err = u.serviceRepo.ReadOneService(ctx, getServiceRequest)
			if err != nil {
				return nil, err
			}

			services = append(services, service)
		}

		p := newResponse(permission)
		p.Services = services
		resp = append(resp, p)
	}

	return resp, nil
}

func (u usecase) FetchPermissionWithPaginate(ctx context.Context, request entity.FilterSearchRequest) ([]entity.PermissionResponse, domain.PaginateApp, error) {
	req := entity.CollectionRequest{Filters: request.Filters, Searches: request.Searches}
	req.TableName = domain.PermissionTable
	req.Model = entity.PermissionModel{}
	permissions, err := u.permissionRepo.ReadManyPermission(ctx, req)
	var resp []entity.PermissionResponse
	var paginate domain.PaginateApp
	if err != nil {
		return resp, paginate, err
	}

	paginate.Total = len(permissions)
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

	for _, permission := range permissions {
		p := permission.ToResponse()
		getPrReq := entity.NewCollectionRequest(domain.PermissionServiceTable).SetModel(entity.PermissionServicesResponse{}).EQ("permission_id", p.ID)
		services, err := u.prRepo.ReadPermissionServices(ctx, getPrReq)
		if err != nil {
			return nil, domain.PaginateApp{}, err
		}

		for _, service := range services {
			getServiceReq := entity.NewCollectionRequest(domain.ServicesTable).SetModel(entity.ServiceResponse{}).EQ("id", service.ServiceID)
			oneService, err := u.serviceRepo.ReadOneService(ctx, getServiceReq)
			if err != nil {
				return nil, domain.PaginateApp{}, err
			}

			p.Services = append(p.Services, oneService)
		}

		resp = append(resp, p)
	}

	return resp, paginate, nil
}

func (u usecase) PatchPermission(ctx context.Context, permission entity.PermissionUpdateRequest) error {
	_, err := u.permissionRepo.ReadPermission(ctx, permission.ID)
	if err != nil {
		return err
	}

	//req := entity.NewCollectionRequest(domain.PermissionServiceTable, entity.PermissionServicesResponse{})
	//req.Filters = append(req.Filters, entity.Filter{Field: "permission_id", Value: []interface{}{permission.ID}})
	//if err != nil {
	//	return err
	//}

	return u.permissionRepo.UpdatePermission(ctx, permission)
}

func newResponse(per entity.PermissionEntity) entity.PermissionResponse {
	return entity.PermissionResponse{
		ID:          per.ID,
		Name:        per.Name,
		IsActive:    per.IsActive,
		Description: per.Description,
		CreatedAt:   per.CreatedAt,
		UpdatedAt:   per.UpdatedAt,
	}
}

func newServiceResponse(services []entity.ServiceModel) (resp []entity.ServiceResponse) {
	for _, service := range services {
		resp = append(resp, entity.ServiceResponse{
			ID:          service.ID.Int64,
			Name:        service.Name.String,
			Path:        service.Name.String,
			Function:    service.Function.String,
			Method:      service.Method.String,
			IsActive:    service.IsActive.Bool,
			Description: service.Description.String,
			CreatedAt:   service.CreatedAt.Time,
		})
	}

	return
}

func Merge(channels ...<-chan map[string]interface{}) <-chan map[string]interface{} {
	switch len(channels) {
	case 0:
		c := make(chan map[string]interface{})
		close(c)
		return c
	case 1:
		return channels[0]
	case 2:
		return mergeTwo(channels[0], channels[1])
	default:
		m := len(channels) / 2
		return mergeTwo(
			Merge(channels[:m]...),
			Merge(channels[m:]...))
	}
}

func mergeTwo(a, b <-chan map[string]interface{}) <-chan map[string]interface{} {
	c := make(chan map[string]interface{})
	go func() {
		defer close(c)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}

func removeDuplicateInt(a, b []int64) []int64 {
	intSlice := a
	intSlice = append(intSlice, b...)
	allKeys := make(map[int64]bool)
	var list []int64
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func intInSlice(arr []int64, n int64) bool {
	for _, a := range arr {
		if a == n {
			return true
		}
	}

	return false
}
