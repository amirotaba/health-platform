package echo

import (
	"context"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	accountDelivery "git.paygear.ir/giftino/account/internal/features/account/delivery/http"
	accountTypeDelivery "git.paygear.ir/giftino/account/internal/features/account_type/delivery/http"
	authDelivery "git.paygear.ir/giftino/account/internal/features/authentication/delivery/http"
	channelDelivery "git.paygear.ir/giftino/account/internal/features/channel/delivery/http"
	channelAccountDelivery "git.paygear.ir/giftino/account/internal/features/channel_account/delivery/http"
	channelRuleDelivery "git.paygear.ir/giftino/account/internal/features/channel_rule/delivery/http"
	otpDelivery "git.paygear.ir/giftino/account/internal/features/otp/delivery/http"
	permissionDelivery "git.paygear.ir/giftino/account/internal/features/permission/delivery/http"
	perServiceDelivery "git.paygear.ir/giftino/account/internal/features/permission_service/delivery/http"
	roleDelivery "git.paygear.ir/giftino/account/internal/features/role/delivery/http"
	rolePerDelivery "git.paygear.ir/giftino/account/internal/features/role_permission/delivery/http"
	serviceDelivery "git.paygear.ir/giftino/account/internal/features/service/delivery/http"
	tokenDelivery "git.paygear.ir/giftino/account/internal/features/token/delivery/http"
	"git.paygear.ir/giftino/account/internal/utils/wpool"
)

func New(services domain.AccountServices) {
	e := echo.New()
	e.Use(middleware.CORS())

	otpDelivery.New(e, services.OtpUsecase, services.Middleware)
	authDelivery.New(e, services.AuthUsecase, services.AccountUsecase, services.Middleware)
	tokenDelivery.New(e, services.TokenUsecase, services.Middleware)
	accountDelivery.New(e, services.AccountUsecase, services.Middleware)
	channelDelivery.New(e, services.ChannelUsecase, services.Middleware)
	channelRuleDelivery.New(e, services.ChannelRuleUsecase, services.Middleware)
	roleDelivery.New(e, services.RoleUsecase, services.Middleware)
	rolePerDelivery.New(e, services.RolePermission, services.Middleware, services.Logger)
	perServiceDelivery.New(e, services.PermissionServices, services.Middleware)
	serviceDelivery.New(e, services.ServiceUsecase, services.Middleware, services.Logger)
	accountTypeDelivery.New(e, services.AccountTypeUsecase, services.Middleware)
	permissionDelivery.New(e, services.PermissionUsecase, services.Middleware)
	channelAccountDelivery.New(e, services.ChannelAccount, services.Middleware, services.Logger)

	rs := e.Routes()
	resource := routeToService(rs)
	jobsCount := len(resource)
	workerCount := 5
	jobs := make([]wpool.Job, jobsCount)
	for i := 0; i < jobsCount; i++ {
		jobs[i] = wpool.Job{
			Descriptor: wpool.JobDescriptor{
				ID:       wpool.JobID(fmt.Sprintf("%v", i)),
				JType:    "anyType",
				Metadata: nil,
			},
			ExecFn: func(ctx context.Context, arg interface{}) error {
				input := arg.(entity.CreateServiceRequest)
				_, err := services.ServiceUsecase.UpsertService(ctx, input)
				if err != nil {
					log.Println("error: ", err)
					return err
				}

				//log.Printf("a service with id %d created \n", serviceID)
				return nil
			},
			Args: resource[i],
		}
	}

	wp := wpool.New(workerCount)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	go wp.GenerateFrom(jobs)
	go wp.Run(ctx)

	go func() {
		for {
			select {
			case r, ok := <-wp.Results():
				if !ok {
					continue
				}
				log.Println("error: ", r.Err)
			case <-wp.Done:
				return
			default:
			}
		}
	}()
	e.Logger.Fatal(e.Start(":" + services.HttpPort))
}

func routeToService(routes []*echo.Route) (services []entity.CreateServiceRequest) {
	for _, route := range routes {
		f := strings.Split(route.Name, ".")
		funcName := strings.Split(f[len(f)-1], "-")[0]
		services = append(services, entity.CreateServiceRequest{
			Name:     route.Method + "-" + route.Path,
			Code:     f[len(f)-1],
			Path:     route.Path,
			Function: funcName,
			Method:   route.Method,
			IsActive: true,
		})
	}

	return
}
