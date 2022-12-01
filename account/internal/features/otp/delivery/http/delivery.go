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
	engine     *echo.Echo
	otpUsecase domain.OtpUsecase
	middleware domain.Middleware
}

func New(echo *echo.Echo, otpUsecase domain.OtpUsecase, middleware domain.Middleware) {
	handler := &delivery{
		engine:     echo,
		otpUsecase: otpUsecase,
		middleware: middleware,
	}

	v1 := handler.engine.Group("/v1")
	{
		//v1.Use(handler.middleware.AccountAuthentication)
		auth := v1.Group("/otp")
		{
			auth.POST("/verify", handler.VerifyOtp)
			auth.POST("/new", handler.Otp)
		}
	}
}

func (d *delivery) VerifyOtp(c echo.Context) error {
	var request entity.OtpVerifyRequest
	log.Println(request)
	if err := c.Bind(&request); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	token, err := d.otpUsecase.OtpVerify(ctx, request)
	if err != nil {
		switch err.(type) {
		case utils.WrongOtpError:
			return c.JSON(http.StatusForbidden, utils.JSONError{Message: err.Error()})
		case utils.ExpireError:
			return c.JSON(http.StatusExpectationFailed, utils.JSONError{Message: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, utils.JSONError{Message: err.Error()})

		}
	}

	return c.JSON(http.StatusOK, token)
}

func (d *delivery) Otp(c echo.Context) error {
	var request entity.OtpRequest
	log.Println(request)
	if err := c.Bind(&request); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	otp, err := d.otpUsecase.NewOtp(ctx, request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, utils.JSONError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": otp})
}
