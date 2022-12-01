package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/utils"
)

type delivery struct {
	engine     *echo.Echo
	usecase    domain.AccountTypeUsecase
	middleware domain.Middleware
}

func New(echo *echo.Echo, usecase domain.AccountTypeUsecase, middleware domain.Middleware) {
	handler := &delivery{
		engine:     echo,
		usecase:    usecase,
		middleware: middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.AccountAuthentication)
		{
			v1.POST("/account_type", handler.SaveAccountType, handler.middleware.AccountAuthentication(handler.SaveAccountType))
			v1.GET("/account_type", handler.GetAccountType, handler.middleware.AccountAuthentication(handler.GetAccountType))
			v1.PATCH("/account_type/:id", handler.UpdateAccountType, handler.middleware.AccountAuthentication(handler.UpdateAccountType))
		}
	}
}

func (d *delivery) SaveAccountType(c echo.Context) error {
	var accountType domain.AccountTypeEntity
	err := utils.CustomDecoder(c.Request().Body, &accountType) // c.Bind(&accountType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(accountType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.usecase.AddAccountType(ctx, &accountType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (d *delivery) GetAccountType(c echo.Context) error {
	ctx := c.Request().Context()
	value := c.QueryParams()
	var accountTypes []domain.AccountTypeResponse
	var err error
	id, ok := value["id"]
	if ok {
		accountTypes, err = d.usecase.FetchAccountType(ctx, id)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
		}
	} else {
		accountTypes, err = d.usecase.FetchAccountTypes(ctx)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, accountTypes)
}

func (d *delivery) GetAccountTypes(c echo.Context) error {
	ctx := c.Request().Context()
	accountTypes, err := d.usecase.FetchAccountTypes(ctx)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, accountTypes)
}

func (d *delivery) UpdateAccountType(c echo.Context) error {
	id := c.Param("id")
	var accountType domain.AccountTypeUpdateRequest
	err := utils.CustomDecoder(c.Request().Body, &accountType) // c.Bind(&accountType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(accountType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	accountType.ID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.usecase.PatchAccountType(ctx, &accountType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}
