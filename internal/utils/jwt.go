package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwt(dur time.Duration, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(dur),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}
