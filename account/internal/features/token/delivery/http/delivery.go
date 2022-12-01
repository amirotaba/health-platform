package http

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/account/entity"
	"git.paygear.ir/giftino/account/internal/utils"
)

type delivery struct {
	engine       *echo.Echo
	tokenUsecase domain.TokenUsecase
	middleware   domain.Middleware
}

func New(echo *echo.Echo, tokenUsecase domain.TokenUsecase, middleware domain.Middleware) {
	handler := &delivery{
		engine:       echo,
		tokenUsecase: tokenUsecase,
		middleware:   middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.AccountAuthentication)
		auth := v1.Group("/token")
		{
			auth.POST("/refresh", handler.RefreshToken)
			//auth.POST("/new", handler.New)
		}
	}
}

func (d *delivery) RefreshToken(c echo.Context) error {
	var request entity.RefreshTokenRequest
	log.Println(request)
	if err := c.Bind(&request); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	token, err := d.tokenUsecase.RefreshToken(ctx, request)
	if err != nil {
		switch err.(type) {
		case utils.ExpireError:
			return c.JSON(http.StatusExpectationFailed, utils.JSONError{Message: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})

		}
	}

	return c.JSON(http.StatusOK, token)
}
