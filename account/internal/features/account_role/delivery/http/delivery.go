package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type delivery struct {
	engine     *echo.Echo
	usecase    domain.AccountRoleUsecase
	middleware domain.Middleware
	logger     *log.Logger
}

func New(echo *echo.Echo, usecase domain.AccountRoleUsecase, middleware domain.Middleware) {
	handler := &delivery{
		engine:     echo,
		usecase:    usecase,
		middleware: middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.AccountAuthentication)
		{
			v1.POST("/account/roles", handler.AssignAccountRoles, handler.middleware.AccountAuthentication(handler.AssignAccountRoles))
			v1.GET("/account/roles/", handler.ShowAccountRoles, handler.middleware.AccountAuthentication(handler.ShowAccountRoles))
			v1.PATCH("/account/roles/:id", handler.AssignAccountRoles, handler.middleware.AccountAuthentication(handler.AssignAccountRoles))
			v1.DELETE("/account/roles/:id", handler.UnAssignAccountRoles, handler.middleware.AccountAuthentication(handler.UnAssignAccountRoles))
		}
	}
}

func (d *delivery) AssignAccountRoles(c echo.Context) error {
	id := c.Param("id")
	var request entity.CreateRolesToAccountRequest
	err := utils.CustomDecoder(c.Request().Body, &request)
	//err := c.Bind(&request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	request.AccountId, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.usecase.AssignRolesToAccount(ctx, request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (d *delivery) ShowAccountRoles(c echo.Context) error {
	paging := new(entity.CollectionRequest)
	var err error
	if c.QueryParam(domain.KeyPage) != domain.KeyEmptyString {
		paging.Page, err = strconv.Atoi(c.QueryParam(domain.KeyPage))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
		}
	} else {
		paging.Page = 1
	}

	if c.QueryParam(domain.KeyPerPage) != domain.KeyEmptyString {
		paging.PerPage, err = strconv.Atoi(c.QueryParam(domain.KeyPerPage))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
		}
	} else {
		paging.PerPage = 10
	}

	if c.QueryParam(domain.KeyFilters) != domain.KeyEmptyString {
		paging.Filters, err = prepareFilter(c.QueryParam(domain.KeyFilters))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
		}
	} else {
		paging.Filters = nil
	}

	ctx := c.Request().Context()
	access, paginate, err := d.usecase.FetchAccountRoles(ctx, *paging)
	if err != nil {
		d.logger.Info(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, domain.AccountRolePaginateResponseApp{
		Result:     access,
		Pagination: paginate,
	})
}

func (d *delivery) UnAssignAccountRoles(c echo.Context) error {
	id := c.Param("id")
	var request entity.CreateRolesToAccountRequest
	err := utils.CustomDecoder(c.Request().Body, &request)
	//err := c.Bind(&request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	request.AccountId, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.usecase.UnAssignRolesToAccount(ctx, request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (d *delivery) UpdateAccess(c echo.Context) error {
	id := c.Param("id")
	var request entity.UpdateRolesToAccountRequest
	err := utils.CustomDecoder(c.Request().Body, &request)
	//err := c.Bind(&request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	request.AccountId, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	result, err := d.usecase.PatchAccountRolesN(ctx, request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	_ = result
	return c.JSON(http.StatusOK, "OK")
}

func prepareFilter(base string) (filters []entity.Filter, err error) {
	fs := strings.Split(base, ".")
	for _, s := range fs {
		var filter entity.Filter
		fields := strings.Split(s, ":")
		if fields[1] == "[]" {
			continue
		} else {
			args := strings.Split(fields[1][1:len(fields[1])-1], ",")
			sts := make([]interface{}, len(args))
			filter.Field = fields[0]
			if filter.Field == domain.KeyToDate || filter.Field == domain.KeyFromDate {
				for key, value := range args {
					filter.Type = "date"
					sts[key] = value
				}
			} else {
				for key, value := range args {
					sts[key] = value
				}
			}

			filter.Value = sts
		}

		filters = append(filters, filter)
	}
	return
}
