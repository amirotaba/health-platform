package http

import (
	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type delivery struct {
	engine     *echo.Echo
	usecase    domain.PermissionServiceUsecase
	middleware domain.Middleware
	logger     *log.Logger
}

func New(echo *echo.Echo, usecase domain.PermissionServiceUsecase, middleware domain.Middleware) {
	handler := &delivery{
		engine:     echo,
		usecase:    usecase,
		middleware: middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.UserAuthentication)
		{
			v1.POST("/permissions/services", handler.NewPermissionServices, handler.middleware.AccountAuthentication(handler.NewPermissionServices))
			v1.GET("/permissions/services", handler.ShowPermissionServices, handler.middleware.AccountAuthentication(handler.ShowPermissionServices))
			v1.PATCH("/permissions/services/:permission_id", handler.PatchPermissionServices, handler.middleware.AccountAuthentication(handler.PatchPermissionServices))
		}
	}
}

func (d *delivery) ShowPermissionServices(c echo.Context) error {
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

	var searches entity.FilterSearchRequest
	searches.Searches, searches.Filters = utils.
		GetSearchAndFilter(c,
			[]string{}, []string{"id", "permission_id", "service_id"})

	ctx := c.Request().Context()
	access, paginate, err := d.usecase.FetchPermissionServices(ctx, searches)
	if err != nil {
		d.logger.Info(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	_ = paginate
	//return c.JSON(http.StatusOK, domain.PermissionServicesPaginateResponseApp{
	//	Result:     access,
	//	Pagination: paginate,
	//})

	return c.JSON(http.StatusOK, access)
}

func (d *delivery) NewPermissionServices(c echo.Context) error {
	var ps entity.AssignServicesToPermission
	err := c.Bind(&ps)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	res, err := d.usecase.NewPermissionServices(ctx, ps)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (d *delivery) PatchPermissionServices(c echo.Context) error {
	var ps entity.UpdateAssignedServicesToPermission
	err := c.Bind(&ps)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ps.PermissionID, err = strconv.ParseInt(c.Param("permission_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	res, err := d.usecase.PatchPermissionServices(ctx, ps)
	if err != nil {
		d.logger.Info(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
