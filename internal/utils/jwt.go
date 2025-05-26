package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenTokenString(jwtSecret, username string, dur time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(dur).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}

func ParseToken(jwtSecret, tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("wrong signing algorithm")
		}

		return []byte(jwtSecret), nil
	})
}
