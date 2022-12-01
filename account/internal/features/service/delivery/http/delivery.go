package http

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type delivery struct {
	engine         *echo.Echo
	serviceUsecase domain.ServiceUsecase
	middleware     domain.Middleware
	logger         *log.Logger
}

func New(echo *echo.Echo, serviceUsecase domain.ServiceUsecase, middleware domain.Middleware, logger *log.Logger) {
	handler := &delivery{
		engine:         echo,
		serviceUsecase: serviceUsecase,
		middleware:     middleware,
		logger:         logger,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.AccountAuthentication)
		service := v1.Group("/service")
		{
			//service.POST("/new", handler.New)
			service.GET("/list", handler.ShowServices, handler.middleware.AccountAuthentication(handler.ShowServices))
			service.PATCH("/:id", handler.UpdateServices, handler.middleware.AccountAuthentication(handler.UpdateServices))
		}
	}
}

//func (d *delivery) New(c echo.Context) error {
//	var request entity.CreateServiceRequest
//	log.Println(request)
//	if err := c.Bind(&request); err != nil {
//		log.Println(err)
//		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
//	}
//
//	ctx := c.Request().Context()
//	service, err := d.serviceUsecase.NewService(ctx, request)
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
//	}
//
//	return c.JSON(http.StatusOK, service)
//}

func (d *delivery) ShowServices(c echo.Context) error {
	//paging := new(entity.CollectionRequest)
	//var err error
	//if c.QueryParam(domain.KeyPage) != domain.KeyEmptyString {
	//	paging.Page, err = strconv.Atoi(c.QueryParam(domain.KeyPage))
	//	if err != nil {
	//		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	//	}
	//} else {
	//	paging.Page = 1
	//}
	//
	//if c.QueryParam(domain.KeyPerPage) != domain.KeyEmptyString {
	//	paging.PerPage, err = strconv.Atoi(c.QueryParam(domain.KeyPerPage))
	//	if err != nil {
	//		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	//	}
	//} else {
	//	paging.PerPage = 10
	//}
	//
	//if c.QueryParam(domain.KeyFilters) != domain.KeyEmptyString {
	//	paging.Filters, err = prepareFilter(c.QueryParam(domain.KeyFilters))
	//	if err != nil {
	//		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	//	}
	//} else {
	//	paging.PerPage = 10
	//}

	searches, err := utils.
		GetSearchAndFilterNew(c,
			[]entity.QuerySearch{}, []entity.QueryFilter{{
				TableName: domain.ServicesTable,
				Field:     "id",
			}})

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	services, paginate, err := d.serviceUsecase.FetchServices(ctx, searches)
	if err != nil {
		d.logger.Info(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, domain.ServicePaginateResponseApp{
		Result:     services,
		Pagination: paginate,
	})
}

func (d *delivery) UpdateServices(c echo.Context) error {
	var service entity.UpdateServiceRequest
	var err error
	if err = c.Bind(&service); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	service.ID, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.serviceUsecase.PatchService(ctx, service)
	if err != nil {
		d.logger.Info(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}
