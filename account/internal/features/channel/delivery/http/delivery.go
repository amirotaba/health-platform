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
	engine         *echo.Echo
	channelUsecase domain.ChannelUsecase
	middleware     domain.Middleware
}

func New(echo *echo.Echo, usecase domain.ChannelUsecase, middleware domain.Middleware) {
	handler := &delivery{
		engine:         echo,
		channelUsecase: usecase,
		middleware:     middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.AccountAuthentication)
		{
			v1.POST("/channel", handler.SaveChannel, handler.middleware.AccountAuthentication(handler.SaveChannel))
			v1.GET("/channel", handler.GetChannel, handler.middleware.AccountAuthentication(handler.GetChannel))
			v1.PATCH("/channel/:id", handler.UpdateChannel, handler.middleware.AccountAuthentication(handler.UpdateChannel))
		}
	}
}

func (d *delivery) SaveChannel(c echo.Context) error {
	var channel entity.CreateChannelRequest
	err := utils.CustomDecoder(c.Request().Body, &channel) // c.Bind(&channel)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(channel)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	_, err = d.channelUsecase.AddChannel(ctx, channel)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (d *delivery) GetChannel(c echo.Context) error {
	// todo handle type error like 400 500
	searches, err := utils.
		GetSearchAndFilterNew(c,
			[]entity.QuerySearch{{
				TableName: domain.ChannelTable,
				Field:     "name",
			}, {
				TableName: domain.ChannelTable,
				Field:     "owner_phone_number",
			}}, []entity.QueryFilter{{
				TableName: domain.ChannelTable,
				Field:     "id",
			}})

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var result domain.ChannelPaginateResponseApp
	result.Result, result.Pagination, err = d.channelUsecase.FetchChannelsWithPaginate(c.Request().Context(), searches)
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

func (d *delivery) UpdateChannel(c echo.Context) error {
	id := c.Param("id")
	var channel entity.UpdateChannelRequest
	err := utils.CustomDecoder(c.Request().Body, &channel) // c.Bind(&channel)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(channel)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	channel.ID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.channelUsecase.PatchChannel(ctx, channel)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}
