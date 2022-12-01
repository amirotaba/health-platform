package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"git.paygear.ir/giftino/account/internal/account/domain"
)

type jwtService struct{}

func NewJWTService() domain.JWTService {
	return &jwtService{}
}

func (s *jwtService) GenerateN(claims domain.Claim, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *jwtService) Generate(id int64, rolesId, permissionsID []int64, secretKey string, exp time.Time) (string, error) {
	claim := &domain.Claim{
		ID:            id,
		RolesID:       rolesId,
		PermissionsID: permissionsID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    "",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *jwtService) Validate(token, secretKey string) bool {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid otp: %v", token)
		}

		return []byte(secretKey), nil
	})

	return err == nil
}

func (s *jwtService) Extract(token, secretKey string) (domain.Claim, error) {
	var c domain.Claim
	parser := new(jwt.Parser)
	_, _, err := parser.ParseUnverified(token, &c)
	if err != nil {
		return domain.Claim{}, err
	}

	return c, nil
}
