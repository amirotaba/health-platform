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

type delivery struct {
	engine     *echo.Echo
	usecase    domain.RoleUsecase
	middleware domain.Middleware
}

func New(echo *echo.Echo, usecase domain.RoleUsecase, middleware domain.Middleware) {
	handler := &delivery{
		engine:     echo,
		usecase:    usecase,
		middleware: middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.UserAuthentication)
		{
			v1.POST("/role", handler.SaveRole, handler.middleware.AccountAuthentication(handler.SaveRole))
			v1.GET("/role", handler.GetRole, handler.middleware.AccountAuthentication(handler.GetRole))
			v1.PATCH("/role/:id", handler.UpdateRole, handler.middleware.AccountAuthentication(handler.UpdateRole))
			//v1.PATCH("/role/permissions/:id", handler.AssignRoleRoles, handler.middleware.UserAuthentication(handler.AssignRoleRoles))
			//v1.DELETE("/role/permissions/:id", handler.UnAssignRoleRoles, handler.middleware.UserAuthentication(handler.UnAssignRoleRoles))
		}
	}
}

func (d *delivery) SaveRole(c echo.Context) error {
	var role entity.CreateRoleRequest
	err := c.Bind(&role)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(role)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	log.Println(err)
	ctx := c.Request().Context()
	_, err = d.usecase.AddRole(ctx, role)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (d *delivery) GetRole(c echo.Context) error {
	var searches entity.FilterSearchRequest
	var err error
	searches.Searches, searches.Filters = utils.GetSearchAndFilter(c, []string{"name"}, []string{"id"})

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

	var result domain.RolePaginateResponseApp
	result.Result, result.Pagination, err = d.usecase.FetchRoleWithPaginate(c.Request().Context(), searches)
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

func (d *delivery) UpdateRole(c echo.Context) error {
	id := c.Param("id")
	var role entity.RoleUpdateRequest
	err := utils.CustomDecoder(c.Request().Body, &role)
	//err := c.Bind(&role)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(role)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	role.ID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.usecase.PatchRole(ctx, role)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}
