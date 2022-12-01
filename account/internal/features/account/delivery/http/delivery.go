package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"git.paygear.ir/giftino/account/internal/account/aggregate"
	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type delivery struct {
	engine     *echo.Echo
	usecase    domain.AccountUsecase
	middleware domain.Middleware
}

func New(echo *echo.Echo, usecase domain.AccountUsecase, middleware domain.Middleware) {
	handler := &delivery{
		engine:     echo,
		usecase:    usecase,
		middleware: middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.AccountAuthentication)
		{
			v1.POST("/account", handler.SaveAccount, handler.middleware.AccountAuthentication(handler.SaveAccount))
			v1.GET("/account", handler.GetAccount, handler.middleware.AccountAuthentication(handler.GetAccount))
			v1.GET("/account/profile/:id", handler.GetAccountProfile, handler.middleware.AccountAuthentication(handler.GetAccountProfile))
			v1.PATCH("/account/:id", handler.UpdateAccount, handler.middleware.AccountAuthentication(handler.UpdateAccount))
			v1.PATCH("/account/profile/:id", handler.UpdateAccountProfile, handler.middleware.AccountAuthentication(handler.UpdateAccountProfile))
		}
	}
}

func (d *delivery) SaveAccount(c echo.Context) error {
	var account entity.CreateAccountRequest
	err := utils.CustomDecoder(c.Request().Body, &account) // c.Bind(&account)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(account)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	_, err = d.usecase.AddAccount(ctx, account)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (d *delivery) SearchAccount(c echo.Context) error {
	var req entity.FilterSearchRequest
	err := utils.CustomDecoder(c.Request().Body, &req) // c.Bind(&account)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	account, paginate, err := d.usecase.FetchAccountsWithPaginate(ctx, req)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, domain.AccountPaginateResponseApp{
		Result:     account,
		Pagination: paginate,
	})
}

func (d *delivery) GetAccount(c echo.Context) error {
	// todo handle type error like 400 500
	var searches entity.FilterSearchRequest
	var err error
	searches.Searches, searches.Filters = utils.
		GetSearchAndFilter(c,
			[]string{"last_name", "first_name", "phone_number"}, []string{"id"})

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

	var result domain.AccountPaginateResponseApp
	result.Result, result.Pagination, err = d.usecase.FetchAccountsWithPaginate(c.Request().Context(), searches)
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

func (d *delivery) GetAccountProfile(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	identity := c.Get("identity").(aggregate.Account)
	if identity.AccountData.ID != id {
		return c.JSON(http.StatusForbidden, "you have not access to this user profile")
	}

	ctx := c.Request().Context()
	account, err := d.usecase.FetchAccountByID(ctx, id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, account)
}

func (d *delivery) UpdateAccount(c echo.Context) error {
	id := c.Param("id")
	var account entity.UpdateAccountRequest
	err := utils.CustomDecoder(c.Request().Body, &account) // c.Bind(&account)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(account)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	account.ID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.usecase.PatchAccount(ctx, account)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (d *delivery) UpdateAccountProfile(c echo.Context) error {
	id := c.Param("id")
	var account entity.UpdateAccountProfileRequest
	err := utils.CustomDecoder(c.Request().Body, &account) // c.Bind(&account)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(account)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	account.ID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	identity := c.Get("identity").(aggregate.Account)
	if identity.AccountData.ID != account.ID {
		return c.JSON(http.StatusForbidden, "you have not access to modify this user")
	}

	ctx := c.Request().Context()
	err = d.usecase.PatchAccountProfile(ctx, account)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}
