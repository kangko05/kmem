package middlewares

import (
	"fmt"
	"kmem/internal/config"
	"kmem/internal/models"
	"kmem/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func validateToken(ctx *gin.Context, cookieName, jwtSecret string) (string, error) {
	tokenStr, err := ctx.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("failed to find cookie: %v", err)
	}

	token, err := utils.ParseToken(jwtSecret, tokenStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to get claims from token")
	}

	val, ok := claims["username"]
	if !ok {
		return "", fmt.Errorf("failed to get username from token")
	}

	username, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("invalid token claim")
	}

	return username, nil
}

// check access token & refresh token from cookies
// if access token is expired, refresh
func Auth(conf *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtSecret := conf.JwtSecretKey()

		username, err := validateToken(ctx, utils.ACCESS_TOKEN_KEY, jwtSecret)
		if err != nil {
			log.Printf("failed to validate access token: %v\n", err)
			log.Println("getting refresh token...")

			username, refreshErr := validateToken(ctx, utils.REFRESH_TOKEN_KEY, jwtSecret)
			if refreshErr != nil {
				models.APIResponse{
					Status:  http.StatusUnauthorized,
					Message: fmt.Sprintf("failed to validate tokens: %v", refreshErr),
				}.Send(ctx)

				ctx.Abort()
				return
			}

			// renew tokens
			accessToken, err := utils.GenTokenString(conf.JwtSecretKey(), username, utils.ACCESS_TOKEN_DUR)
			if err != nil {
				models.APIResponse{
					Status:  http.StatusInternalServerError,
					Message: fmt.Sprintf("failed to gen access token: %v", err),
				}.Send(ctx)

				ctx.Abort()
				return
			}

			refreshToken, err := utils.GenTokenString(conf.JwtSecretKey(), username, utils.REFRESH_TOKEN_DUR)
			if err != nil {
				models.APIResponse{
					Status:  http.StatusInternalServerError,
					Message: fmt.Sprintf("failed to gen refresh token: %v", err),
				}.Send(ctx)

				ctx.Abort()
				return
			}

			ctx.SetCookie(utils.ACCESS_TOKEN_KEY, accessToken, int(utils.ACCESS_TOKEN_DUR), "/", "", true, true)
			ctx.SetCookie(utils.REFRESH_TOKEN_KEY, refreshToken, int(utils.REFRESH_TOKEN_DUR), "/", "", true, true)
		}

		ctx.Set("username", username)
		ctx.Next()
	}
}
