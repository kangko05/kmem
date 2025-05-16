package middleware

import (
	"fmt"
	"kmem/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkToken(ctx *gin.Context, tokenType utils.TokenType) (string, error) {
	tokenStr, err := ctx.Cookie(string(tokenType))
	if err != nil {
		return "", utils.TOKEN_NOT_FOUND
	}

	token, claims, err := utils.PasrseJwt(tokenStr)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", utils.INVALID_TOKEN
	}

	claim, ok := claims[utils.USERNAME_KEY]
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	username, ok := claim.(string)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	ctx.Set(utils.USERNAME_KEY, username)

	return username, nil
}

func abort(ctx *gin.Context) {
	ctx.String(http.StatusUnauthorized, "invalid token")
	ctx.Abort()
}

func AuthenticateJwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, accessErr := checkToken(ctx, utils.ACCESS_TOKEN)

		if accessErr == nil {
			ctx.Next()
			return
		}

		switch accessErr {
		case utils.INVALID_TOKEN, utils.TOKEN_NOT_FOUND:
			username, err := checkToken(ctx, utils.REFRESH_TOKEN)
			if err != nil {
				abort(ctx)
				return
			}

			accessToken, err := utils.CreateJwt(utils.ACCESSTOKEN_MAX_AGE, username)
			if err != nil {
				ctx.String(http.StatusInternalServerError, "failed to create token")
				ctx.Abort()
				return
			}

			refreshToken, err := utils.CreateJwt(utils.REFRESHTOKEN_MAX_AGE, username)
			if err != nil {
				ctx.String(http.StatusInternalServerError, "failed to create token")
				ctx.Abort()
				return
			}

			ctx.SetCookie(string(utils.ACCESS_TOKEN), accessToken, int(utils.ACCESSTOKEN_MAX_AGE.Seconds()), "/", "", true, true)
			ctx.SetCookie(string(utils.REFRESH_TOKEN), refreshToken, int(utils.REFRESHTOKEN_MAX_AGE.Seconds()), "/", "", true, true)

			ctx.Next()
			return
		default:
			abort(ctx)
			return
		}
	}
}
