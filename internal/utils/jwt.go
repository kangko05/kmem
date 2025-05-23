package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwt(dur time.Duration, username, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		USERNAME_KEY: username,
		"exp":        time.Now().Add(dur).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func PasrseJwt(secretKey, tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, INVALID_TOKEN
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("claim is not map")
	}

	return token, claim, nil
}
