package usecase

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type usecase struct {
	serviceRepo domain.ServiceRepository
	logger      *log.Logger
}

func New(serviceRepo domain.ServiceRepository, logger *log.Logger) domain.ServiceUsecase {
	return &usecase{
		serviceRepo: serviceRepo,
		logger:      logger,
	}
}

func (u usecase) NewService(ctx context.Context, service entity.CreateServiceRequest) (int64, error) {
	log.Println("-------------------------------- NewService Usecase ---------------------------------------")
	defer log.Println("-------------------------------- NewService Usecase ---------------------------------------")
	// todo: check if service exist or not
	var r entity.ServiceResponse
	getRequest := entity.NewCollectionRequest(domain.ServicesTable).SetModel(r).EQ("code", service.Code)
	r, err := u.serviceRepo.ReadOneService(ctx, getRequest)
	log.Println(r)
	log.Println(err)
	switch err.(type) {
	case utils.NotFoundError:
		req := entity.NewCollectionRequest(domain.ServicesTable).SetModel(service)
		return u.serviceRepo.CreateService(ctx, req)
	case utils.MysqlInternalServerError:
		return 0, err
	default:
		return 0, errors.New("service exist")
	}
}

func (u usecase) UpsertService(ctx context.Context, service entity.CreateServiceRequest) (int64, error) {
	log.Println("-------------------------------- UpsertService Usecase ---------------------------------------")
	defer log.Println("-------------------------------- UpsertService Usecase ---------------------------------------")
	// todo: check if service exist or not
	var r entity.ServiceResponse
	getRequest := entity.NewCollectionRequest(domain.ServicesTable).SetModel(r).EQ("path", service.Path).EQ("method", service.Method)
	r, err := u.serviceRepo.ReadOneService(ctx, getRequest)
	log.Println("r: ", r)
	log.Println("err: ", err)
	switch err.(type) {
	case utils.NotFoundError:
		req := entity.NewCollectionRequest(domain.ServicesTable).SetModel(service)
		return u.serviceRepo.CreateService(ctx, req)
	case utils.MysqlInternalServerError:
		return 0, err
	default:
		u.logger.Infof("Name: %s", r.Name)
		update := entity.InternalUpdateServiceRequest{
			Name:        r.Name,
			IsActive:    r.IsActive,
			Function:    r.Function,
			Description: r.Description,
			UpdatedAt:   time.Now(),
		}

		//service.UpdatedAt = time.Now()
		u.logger.Infof("service with id %d updated", r.ID)
		req := entity.NewUpdateCollectionRequest(domain.ServicesTable).SetModel(update).EQ("id", r.ID)
		return r.ID, u.serviceRepo.UpdateService(ctx, req)
	}
}

func (u usecase) FetchServices(ctx context.Context, request entity.FilterSearchRequest) (services []entity.ServiceResponse, paginate domain.PaginateApp, err error) {
	req := entity.NewCollectionRequest(domain.ServicesTable).SetModel(entity.ServiceResponse{})
	req.Filters = request.Filters
	req.Searches = request.Searches

	var service entity.ServiceResponse
	services, err = u.serviceRepo.ReadManyServices(ctx, req)
	if err != nil {
		return services, paginate, err
	}

	paginate.Total = len(services)
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
	log.Println(service)
	return services, paginate, err
}

func (u usecase) PatchService(ctx context.Context, req entity.UpdateServiceRequest) error {
	request := entity.NewUpdateCollectionRequest(domain.ServicesTable).SetModel(req).EQ("id", req.ID)
	req.UpdatedAt = time.Now()
	return u.serviceRepo.UpdateService(ctx, request)
}
