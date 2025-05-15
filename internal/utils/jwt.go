package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	SECRET_KEY = "TODO: SECRET KEY MUST BE HANDLED WITH CARE"
)

// tokens
const (
	ACCESS_TOKEN  = "accessToken"
	REFRESH_TOKEN = "refreshToken"

	ACCESSTOKEN_MAX_AGE  = time.Minute * 20
	REFRESHTOKEN_MAX_AGE = time.Hour * 24 * 7
)

func CreateJwt(dur time.Duration, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(dur),
	})

	return token.SignedString([]byte(SECRET_KEY))
}
