package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"git.paygear.ir/giftino/account/internal/account/aggregate"
	"time"

	"github.com/golang-jwt/jwt"

	"git.paygear.ir/giftino/account/internal/account/domain"
)

func ValidateJWT(token string, secret string) (jwt.MapClaims, error) {
	decodedToken := make([]byte, base64.StdEncoding.DecodedLen(len(secret)))
	_, base64err := base64.StdEncoding.Decode(decodedToken, []byte(secret))
	if base64err != nil {
		return nil, base64err
	}

	// Parse takes the otp string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the otp to identify which key to use, but the parsed otp (head and claims) is provided
	// to the callback, providing flexibility.
	parser := new(jwt.Parser)
	parser.SkipClaimsValidation = true
	parsedToken, parseError := parser.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return decodedToken, nil
	})

	if parseError != nil {
		fmt.Println("jwt parse error", parseError)
		return nil, errors.New("parse error")
	}

	if _vErr := parsedToken.Claims.Valid(); _vErr != nil {
		if vErr, ok := _vErr.(*jwt.ValidationError); ok {
			if vErr.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
				fmt.Println("otp is expired")
				return nil, errors.New("otp is expired")
			} else if vErr.Errors&jwt.ValidationErrorIssuedAt == jwt.ValidationErrorIssuedAt {
				fmt.Println("Token used before issued")
				return nil, errors.New("otp used before issued")
			} else if vErr.Errors&jwt.ValidationErrorNotValidYet == jwt.ValidationErrorNotValidYet {
				fmt.Println("otp is not valid yet")
				// return nil, errors.New("otp is not valid yet")
			}
		} else {
			fmt.Println("jwt validation error")
			return nil, errors.New("jwt validation error")
		}
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid otp")
}

func CreateJWT(id string, account aggregate.Account) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["id"] = id
	claims["role"] = account.RolesData
	claims["exp"] = time.Now().AddDate(0, 1, 0)
	exDate := claims["exp"].(int64)
	claims["aud"] = "" //config.Aud
	token.Claims = claims
	base64Secret, err := Base64Decode(domain.GiftinoSecretKey) //base64Decode(config.GiftinoSecretKey)
	if err != nil {
		return "", "", err
	}

	tokenString, err := token.SignedString(base64Secret)
	if err != nil {
		return "", "", err
	}
	return tokenString, string(exDate), nil
}
