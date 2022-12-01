package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claim struct {
	ID            int64    `json:"id"`
	RolesID       []int64  `json:"roles_id"`
	PermissionsID []int64  `json:"permissions_id"`
	IPs           []string `json:"ips"`
	jwt.RegisteredClaims
}

type JWTService interface {
	Generate(id int64, roleId, permissionsID []int64, secretKey string, time time.Time) (string, error)
	GenerateN(claim Claim, secretKey string) (string, error)
	Validate(token, secretKey string) bool
	Extract(token, secretKey string) (Claim, error)
}
