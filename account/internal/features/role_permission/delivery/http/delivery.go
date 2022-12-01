package http

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type delivery struct {
	engine     *echo.Echo
	usecase    domain.RolePermissionUsecase
	middleware domain.Middleware
	logger     *log.Logger
}

func New(echo *echo.Echo, usecase domain.RolePermissionUsecase, middleware domain.Middleware, logger *log.Logger) {
	handler := &delivery{
		engine:     echo,
		usecase:    usecase,
		middleware: middleware,
		logger:     logger,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.AccountAuthentication)
		{
			v1.POST("/role/permissions", handler.AssignRolePermissions, handler.middleware.AccountAuthentication(handler.AssignRolePermissions))
			v1.GET("/role/permissions", handler.ShowRolePermissions, handler.middleware.AccountAuthentication(handler.ShowRolePermissions))
			v1.PATCH("/role/permissions/:role_id", handler.UpdateRolePermission, handler.middleware.AccountAuthentication(handler.UpdateRolePermission))
		}
	}
}

func (d *delivery) AssignRolePermissions(c echo.Context) error {
	var request entity.CreatePermissionsToRoleRequest
	err := utils.CustomDecoder(c.Request().Body, &request)
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

	ctx := c.Request().Context()
	res, err := d.usecase.NewRolePermissions(ctx, request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (d *delivery) ShowRolePermissions(c echo.Context) error {
	var searches entity.FilterSearchRequest
	var err error
	searches.Searches, searches.Filters = utils.GetSearchAndFilter(c, []string{"name"}, []string{"id", "role_id", "permission_id"})
	ctx := c.Request().Context()
	access, err := d.usecase.FetchRolePermissions(ctx, searches)
	if err != nil {
		d.logger.Info(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, access)
}

func (d delivery) UpdateRolePermission(c echo.Context) error {
	var request entity.UpdateRolePermissionsRequest
	err := utils.CustomDecoder(c.Request().Body, &request)
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

	request.RoleId, err = strconv.ParseInt(c.Param("role_id"), 10, 64)
	log.Println("** : ", request.RoleId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	res, err := d.usecase.PatchRolePermissions(ctx, request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
