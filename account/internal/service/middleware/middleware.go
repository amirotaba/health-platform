package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/utils"
)

type middleware struct {
	authManager domain.AuthUsecase
}

func New(authManager domain.AuthUsecase) domain.Middleware {
	return &middleware{
		authManager: authManager,
	}
}

func (m middleware) RateLimit(formatted string) (echo.HandlerFunc, error) {
	return nil, nil
}

func (m middleware) AccountAuthentication(fn echo.HandlerFunc, roles ...string) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// todo check expired is not valid return internal server error
			token := c.Request().Header.Get("Authorization")
			if token == "" || len(token) < 7 {
				return c.JSON(http.StatusForbidden, "invalid token")
			}

			if token[:7] == "bearer " || token[:7] == "Bearer " {
				token = token[7:]
			}

			ctx := c.Request().Context()
			identity, err := m.authManager.Authentication(ctx, token, domain.GiftinoSecretKey)
			log.Println("999999999999999999999999999999")
			log.Println("888888888888888888888888888888:w:")
			if err != nil {
				log.Println(err.Error())
				return c.JSON(http.StatusUnauthorized, "")
			}

			log.Println("/////////////", identity.RolesData)
			ok := false
			// todo get all account detail for get best decision
			log.Println(utils.RoleToString(identity.RolesData))
			if utils.StringInSlice(utils.RoleToString(identity.RolesData), "administrator") {
				log.Println("nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn", identity.RolesData)
				ok = true
			} else {
				log.Println("yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy", identity.RolesData)

				ok, err = m.authManager.Authorization(ctx, identity.RolesData, c.Path(), c.Request().Method)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, "you have not access to this resource")
				}
			}

			if !ok {
				return c.JSON(http.StatusUnauthorized, "you have not access to this resource")
			}

			c.Set("identity", identity)
			return h(c)
		}
	}
}

func (m middleware) HttpCache(defaultExpire time.Duration) echo.HandlerFunc {
	return nil
}

func (m middleware) AuditLog(c *echo.Context) {
}

func (m middleware) CORS() echo.HandlerFunc {
	return nil
}

func (m *middleware) extractTokenFromBearerToken(bearerToken string) (result string, err error) {
	extractedToken := strings.Split(bearerToken, " ")
	// Verify if the format of the otp is correct
	if len(extractedToken) == 2 {
		result = strings.TrimSpace(extractedToken[1])
	} else {
		err = domain.AuthError{
			Status:  http.StatusBadRequest,
			Message: "Incorrect Format of Authorization Token",
		}
		// c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Incorrect Format of Authorization Token "})
		// c.Abort()
		return "", err
	}
	return result, nil
}
