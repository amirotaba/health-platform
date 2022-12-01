package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/aggregate"
	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type delivery struct {
	engine     *echo.Echo
	usecase    domain.ChannelAccountUsecase
	middleware domain.Middleware
	logger     *log.Logger
}

func New(echo *echo.Echo, usecase domain.ChannelAccountUsecase, middleware domain.Middleware, logger *log.Logger) {
	handler := &delivery{
		engine:     echo,
		usecase:    usecase,
		middleware: middleware,
		logger:     logger,
	}

	v1 := handler.engine.Group("/v1")
	{
		{
			v1.POST("/channel/account", handler.AddChannelAccounts, handler.middleware.AccountAuthentication(handler.AddChannelAccounts))
			v1.GET("/account/channel", handler.ShowChannelAccounts, handler.middleware.AccountAuthentication(handler.ShowChannelAccounts))
			v1.GET("/account/my_channel", handler.ShowMyChannelAccounts, handler.middleware.AccountAuthentication(handler.ShowMyChannelAccounts))
			v1.GET("/account/channel/all", handler.ShowAllChannelAccounts, handler.middleware.AccountAuthentication(handler.ShowAllChannelAccounts, "admin"))
			v1.PATCH("/channel/account/:id", handler.UpdateChannelAccounts, handler.middleware.AccountAuthentication(handler.UpdateChannelAccounts))
			v1.DELETE("/channel/account/:id", handler.DeleteChannelAccounts, handler.middleware.AccountAuthentication(handler.DeleteChannelAccounts))
		}
	}
}

func (d *delivery) AddChannelAccounts(c echo.Context) error {
	var request entity.AddChannelAccountsRequest
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
	res, err := d.usecase.NewChannelAccounts(ctx, request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (d *delivery) ShowChannelAccounts(c echo.Context) error {
	searches, err := utils.
		GetSearchAndFilterNew(c,
			[]entity.QuerySearch{}, []entity.QueryFilter{{
				TableName: domain.ChannelAccountTable,
				Type:      entity.IN,
				Field:     "id",
			}, {
				TableName: domain.ChannelAccountTable,
				Field:     "channel_id",
				Type:      entity.EQ,
			}, {
				TableName: domain.ChannelAccountTable,
				Field:     "account_id",
				Type:      entity.EQ,
			}})

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	identity := c.Get("identity").(aggregate.Account)
	//if !identity.IsAdmin {
	//searches.Filters = append(searches.Filters, entity.Filter{
	//	Field: "account_id",
	//	Value: []interface{}{identity.AccountData.ID},
	//})
	//}

	log.Println(identity.RolesData)
	log.Println(identity.AccountData)
	ctx := c.Request().Context()
	access, paginate, err := d.usecase.FetchChannelAccounts(ctx, searches)
	if err != nil {
		d.logger.Info(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	_ = paginate
	//return c.JSON(http.StatusOK, domain.ChannelAccountPaginateResponseApp{
	//	Result:     access,
	//	Pagination: paginate,
	//})
	return c.JSON(http.StatusOK, access)
}

func (d *delivery) ShowMyChannelAccounts(c echo.Context) error {
	searches, err := utils.
		GetSearchAndFilterNew(c,
			[]entity.QuerySearch{}, []entity.QueryFilter{{
				TableName: domain.ChannelAccountTable,
				Field:     "id",
			}, {
				TableName: domain.ChannelAccountTable,
				Field:     "channel_id",
			}, {
				TableName: domain.ChannelAccountTable,
				Field:     "account_id",
			}})

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	identity := c.Get("identity").(aggregate.Account)
	//if !identity.IsAdmin {
	searches.Filters = append(searches.Filters, entity.Filter{
		Field: "account_id",
		Type:  entity.EQ,
		Value: identity.AccountData.ID,
	})
	//}

	log.Println(identity.RolesData)
	log.Println(identity.AccountData)
	ctx := c.Request().Context()
	access, paginate, err := d.usecase.FetchChannelAccounts(ctx, searches)
	if err != nil {
		d.logger.Info(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	_ = paginate
	//return c.JSON(http.StatusOK, domain.ChannelAccountPaginateResponseApp{
	//	Result:     access,
	//	Pagination: paginate,
	//})
	return c.JSON(http.StatusOK, access)
}

func (d *delivery) ShowAllChannelAccounts(c echo.Context) error {
	searches, err := utils.
		GetSearchAndFilterNew(c,
			[]entity.QuerySearch{},
			[]entity.QueryFilter{{
				TableName: domain.ChannelAccountTable,
				Field:     "id",
			}, {
				TableName: domain.ChannelAccountTable,
				Field:     "channel_id",
			}})

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	identity := c.Get("identity").(aggregate.Account)

	if !identity.IsAdmin {
		searches.Filters = append(searches.Filters,
			entity.Filter{
				TableName: domain.ChannelAccountTable,
				Field:     "account_id",
				Value:     identity.ID,
				Type:      entity.EQ,
			})
	}

	if identity.IsAdmin {
		accountFilter, _ := utils.GetOneFilterNew(c, entity.QueryFilter{
			TableName: domain.ChannelAccountTable,
			Field:     "account_id",
			Type:      entity.EQ,
		})

		searches.Filters = append(searches.Filters, accountFilter)
	}

	ctx := c.Request().Context()
	access, paginate, err := d.usecase.FetchChannelAccounts(ctx, searches)
	if err != nil {
		d.logger.Info(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, domain.ChannelAccountPaginateResponseApp{
		Result:     access,
		Pagination: paginate,
	})
}

func (d *delivery) UpdateChannelAccounts(c echo.Context) error {
	fmt.Println("**********************UpdateChannelAccounts*****************************")
	defer fmt.Println("**********************UpdateChannelAccounts*****************************")
	var request entity.UpdateAccountsToChannelRequest
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

	request.RelationID, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	res, err := d.usecase.PatchChannelAccounts(ctx, request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (d *delivery) DeleteChannelAccounts(c echo.Context) error {
	fmt.Println("**********************UpdateChannelAccounts*****************************")
	defer fmt.Println("**********************UpdateChannelAccounts*****************************")
	var request entity.DeleteAccountsFromChannelRequest
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

	request.RelationID, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	res, err := d.usecase.RemoveChannelAccounts(ctx, request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
