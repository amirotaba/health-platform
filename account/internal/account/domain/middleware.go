package domain

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Middleware interface {
	RateLimit(formatted string) (echo.HandlerFunc, error)
	AccountAuthentication(fn echo.HandlerFunc, roles ...string) echo.MiddlewareFunc
	HttpCache(defaultExpire time.Duration) echo.HandlerFunc
	AuditLog(c *echo.Context)
	CORS() echo.HandlerFunc
}
