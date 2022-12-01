package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type Delivery struct {
	engine     *echo.Echo
	usecase    domain.PermissionUsecase
	middleware domain.Middleware
}

func New(echo *echo.Echo, usecase domain.PermissionUsecase, middleware domain.Middleware) Delivery {
	handler := Delivery{
		engine:     echo,
		usecase:    usecase,
		middleware: middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.UserAuthentication)
		{
			v1.POST("/permission", handler.SavePermission, handler.middleware.AccountAuthentication(handler.SavePermission))
			v1.GET("/permission", handler.GetPermission, handler.middleware.AccountAuthentication(handler.GetPermission))
			v1.PATCH("/permission/:id", handler.UpdatePermission, handler.middleware.AccountAuthentication(handler.UpdatePermission))
		}
	}

	return handler
}

func (d *Delivery) SavePermission(c echo.Context) error {
	var permission entity.CreatePermissionRequest
	err := utils.CustomDecoder(c.Request().Body, &permission)
	//err := c.Bind(&permission)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(permission)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	_, err = d.usecase.AddPermission(ctx, permission)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (d *Delivery) GetPermission(c echo.Context) error {
	var searches entity.FilterSearchRequest
	var err error
	searches.Searches, searches.Filters = utils.
		GetSearchAndFilter(c,
			[]string{"name"}, []string{"id"})

	if c.QueryParam(domain.KeyPage) != domain.KeyEmptyString {
		searches.Page, err = strconv.Atoi(c.QueryParam(domain.KeyPage))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
		}
	} else {
		searches.Page = 1
	}

	if c.QueryParam(domain.KeyPerPage) != domain.KeyEmptyString {
		searches.PerPage, err = strconv.Atoi(c.QueryParam(domain.KeyPerPage))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
		}
	} else {
		searches.PerPage = 10
	}

	//value := c.QueryParams()
	//ids, ok := value["id"]
	//if ok {
	//	var accountID []interface{}
	//	for _, id := range ids {
	//		accountID = append(accountID, id)
	//	}
	//	searches.Filters = append(searches.Filters, entity.Filter{Field: "id", Value: accountID})
	//}
	//
	//name, ok := value["name"]
	//if ok {
	//	searches.Searches = append(searches.Searches, entity.Search{Field: "name", Value: []interface{}{name[0]}})
	//}

	var result domain.PermissionPaginateResponseApp
	result.Result, result.Pagination, err = d.usecase.FetchPermissionWithPaginate(c.Request().Context(), searches)
	switch err.(type) {
	case utils.NotFoundError:
		return c.JSON(http.StatusNotFound, err)
	case utils.MysqlInternalServerError:
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "something is wrong")
	default:
		return c.JSON(http.StatusOK, result)
	}
}

func (d *Delivery) UpdatePermission(c echo.Context) error {
	id := c.Param("id")
	var permission entity.PermissionUpdateRequest
	err := utils.CustomDecoder(c.Request().Body, &permission)
	//err := c.Bind(&permission)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(permission)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	permission.ID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.usecase.PatchPermission(ctx, permission)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}
