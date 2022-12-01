package http

import (
	"fmt"
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
	usecase    domain.ChannelRuleUsecase
	middleware domain.Middleware
}

func New(echo *echo.Echo, usecase domain.ChannelRuleUsecase, middleware domain.Middleware) {
	handler := &delivery{
		engine:     echo,
		usecase:    usecase,
		middleware: middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.AccountAuthentication)
		{
			v1.POST("/channel_rule", handler.SaveChannelRule, handler.middleware.AccountAuthentication(handler.SaveChannelRule))
			v1.GET("/channel_rule", handler.GetChannelRule, handler.middleware.AccountAuthentication(handler.GetChannelRule))
			v1.GET("/channel_rule/:id", handler.GetChannelRules, handler.middleware.AccountAuthentication(handler.GetChannelRules))
			v1.GET("/my_channel/channel_rule/:id", handler.GetMyChannelChannelRules, handler.middleware.AccountAuthentication(handler.GetMyChannelChannelRules))
			v1.PATCH("/channel_rule/:id", handler.UpdateChannelRule, handler.middleware.AccountAuthentication(handler.UpdateChannelRule))
			v1.DELETE("/channel_rule/:id", handler.DeleteChannelRule, handler.middleware.AccountAuthentication(handler.DeleteChannelRule))
		}
	}
}

func (d *delivery) SaveChannelRule(c echo.Context) error {
	var channelRule entity.CreateChannelRuleRequest
	err := utils.CustomDecoder(c.Request().Body, &channelRule) // c.Bind(&channelRule)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(channelRule)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	log.Println(err)
	ctx := c.Request().Context()
	_, err = d.usecase.AddChannelRule(ctx, channelRule)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (d *delivery) GetChannelRule(c echo.Context) error {
	fmt.Println("********************************* Delivery GetChannelRule *************************************************")
	defer fmt.Println("********************************* Delivery GetChannelRule *************************************************")
	ctx := c.Request().Context()
	value := c.QueryParams()
	tagIds, ok := value["tag_id"]
	var err error
	var channelRules []entity.ChannelRuleResponse
	var req entity.GetChannelRuleRequest
	if ok {
		filterID := make([]interface{}, 0)
		for _, id := range tagIds {
			filterID = append(filterID, id)
		}

		req.Filters = append(req.Filters, entity.Filter{
			Field:  "tag_id",
			Type:   entity.IN,
			Values: filterID,
		})
	}

	ids, ok := value["id"]
	if ok {
		filterID := make([]interface{}, 0)
		for _, id := range ids {
			filterID = append(filterID, id)
		}

		req.Filters = append(req.Filters, entity.Filter{
			Field:  "id",
			Type:   entity.IN,
			Values: filterID,
		})
	}

	channelIDs, okay := value["channel_id"]
	if okay {
		channelFilter := make([]interface{}, 0)
		for _, channelID := range channelIDs {
			channelFilter = append(channelFilter, channelID)
		}
		req.Filters = append(req.Filters, entity.Filter{Field: "channel_id", Type: entity.IN, Values: channelFilter})
	}

	// todo check if account has admin role ignore account id
	identity := c.Get("identity").(aggregate.Account)
	log.Println(req)
	req.AccountID = identity.ID
	channelRules, err = d.usecase.FetchChannelRules(ctx, req)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, channelRules)
}

func (d *delivery) GetChannelRules(c echo.Context) error {
	fmt.Println("********************************* Delivery GetChannelRules *************************************************")
	defer fmt.Println("********************************* Delivery GetChannelRules *************************************************")
	ctx := c.Request().Context()
	//value := c.QueryParams()
	//ids, ok := value["id"]
	//var err error
	//var channelRules []entity.ChannelRuleResponse
	var req entity.GetChannelRuleRequest
	//if ok {
	//	filterID := make([]interface{}, 0)
	//	for _, id := range ids {
	//		filterID = append(filterID, id)
	//	}
	//
	//	req.Filters = append(req.Filters, domain.Filter{
	//		Field: "id",
	//		Value: filterID,
	//	})
	//}

	channelID := c.Param("id")
	if channelID != "" {
		req.Filters = append(req.Filters, entity.Filter{
			Field: "channel_id",
			Value: channelID,
			Type:  entity.EQ,
		})
	}

	identity := c.Get("identity").(aggregate.Account)
	log.Println(req)
	req.AccountID = identity.ID
	channelRules, err := d.usecase.FetchChannelRules(ctx, req)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, channelRules)
}

func (d *delivery) GetMyChannelChannelRules(c echo.Context) error {
	fmt.Println("********************************* Delivery GetChannelRules *************************************************")
	defer fmt.Println("********************************* Delivery GetChannelRules *************************************************")
	ctx := c.Request().Context()
	//value := c.QueryParams()
	//ids, ok := value["id"]
	//var err error
	//var channelRules []entity.ChannelRuleResponse
	var req entity.GetChannelRuleRequest
	//if ok {
	//	filterID := make([]interface{}, 0)
	//	for _, id := range ids {
	//		filterID = append(filterID, id)
	//	}
	//
	//	req.Filters = append(req.Filters, domain.Filter{
	//		Field: "id",
	//		Value: filterID,
	//	})
	//}

	channelID := c.Param("id")
	if channelID != "" {
		req.Filters = append(req.Filters, entity.Filter{
			Field: "channel_id",
			Value: channelID,
			Type:  entity.EQ,
		})
	}

	identity := c.Get("identity").(aggregate.Account)
	log.Println(req)
	req.AccountID = identity.ID
	channelRules, err := d.usecase.FetchChannelRulesCategory(ctx, req)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, channelRules)
}

func (d *delivery) UpdateChannelRule(c echo.Context) error {
	id := c.Param("id")
	var channelRule entity.ChannelRuleUpdateRequest
	err := utils.CustomDecoder(c.Request().Body, &channelRule) // c.Bind(&channelRule)

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(channelRule)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	channelRule.ID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.usecase.PatchChannelRule(ctx, channelRule)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}

func (d *delivery) DeleteChannelRule(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = d.usecase.RemoveChannelRuleByID(ctx, id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "OK")
}
