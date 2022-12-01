package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
)

type jwtNService struct {
}

func NewJWTNService() domain.JWTService {
	return &jwtNService{}
}

func (s *jwtNService) Generate(id int64, rolesID, permissionsID []int64, secretKey string, t time.Time) (string, error) {
	// Create the claims
	claims := domain.Claim{
		ID:            id,
		IPs:           nil,
		RolesID:       rolesID,
		PermissionsID: permissionsID,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(t), //jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	log.Printf("%v %v \n", ss, err)
	return ss, nil
}

func (s *jwtNService) GenerateN(claims domain.Claim, secretKey string) (string, error) {
	// Create the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	log.Printf("%v %v \n", ss, err)
	return ss, nil
}

func (s *jwtNService) Validate(token, secretKey string) bool {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if t.Valid {
		fmt.Println("You look nice today")
		return true
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		fmt.Println("That's not even a token")
		return false
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		fmt.Println("Timing is everything")
		return false
	} else {
		fmt.Println("Couldn't handle this token:", err)
		return false
	}
}

func (s *jwtNService) Extract(token, secretKey string) (domain.Claim, error) {
	var c domain.Claim
	t, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if claims, ok := t.Claims.(*domain.Claim); ok && t.Valid {
		return *claims, nil
		//fmt.Printf("%v %v", claims.Foo, claims.RegisteredClaims.Issuer)
	} else {
		fmt.Println(err)
		return domain.Claim{}, err
	}
}
