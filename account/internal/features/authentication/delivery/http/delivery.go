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
	engine         *echo.Echo
	AuthUsecase    domain.AuthUsecase
	AccountUsecase domain.AccountUsecase
	middleware     domain.Middleware
}

func New(echo *echo.Echo, AuthUsecase domain.AuthUsecase,
	AccountUsecase domain.AccountUsecase, middleware domain.Middleware) {
	handler := &delivery{
		engine:         echo,
		AuthUsecase:    AuthUsecase,
		AccountUsecase: AccountUsecase,
		middleware:     middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.AccountAuthentication)
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", handler.SignUp)
			auth.POST("/signin", handler.SighIn)
		}
	}
}

func (d *delivery) SignUp(c echo.Context) error {
	request := new(entity.SingUpRequest)
	log.Println(request.Phone)

	if err := utils.CustomDecoder(c.Request().Body, &request); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	otp, err := d.AuthUsecase.SignUp(ctx, *request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, otp)
}

func (d *delivery) SighIn(c echo.Context) error {
	request := new(entity.SingInRequest)
	log.Println(request)
	if err := utils.CustomDecoder(c.Request().Body, &request); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	if request.Phone == "" || request.Password == "" {
		log.Println("phone or password fields is empty")
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: "phone or password fields is empty"})
	}

	ctx := c.Request().Context()
	token, err := d.AuthUsecase.SignIn(ctx, *request)
	if err != nil {
		// todo redirect to signup
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
