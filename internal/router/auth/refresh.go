package auth

import (
	"fmt"
	"kmem/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: store refresh tokens into db and trace it for better security
func Refresh() func(*gin.Context) {
	return func(ctx *gin.Context) {
		refreshToken, err := ctx.Cookie(utils.REFRESH_TOKEN)
		if err != nil {
			ctx.String(http.StatusBadRequest, "refresh token not found")
			return
		}

		token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(utils.SECRET_KEY), nil
		})
		if err != nil {
			ctx.String(http.StatusUnauthorized, fmt.Sprintf("failed to parse token: %v", err))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.String(http.StatusUnauthorized, "invalid token")
			return
		}

		claim, ok := claims["username"]
		if !ok {
			ctx.String(http.StatusUnauthorized, "invalid token claims")
			return
		}

		username, ok := claim.(string)
		if !ok {
			ctx.String(http.StatusUnauthorized, "invalid token claims")
			return
		}

		accessToken, err := utils.CreateJwt(utils.ACCESSTOKEN_MAX_AGE, username)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "failed to create token")
			return
		}

		refreshToken, err = utils.CreateJwt(utils.REFRESHTOKEN_MAX_AGE, username)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "failed to create token")
			return
		}

		ctx.SetCookie(utils.ACCESS_TOKEN, accessToken, int(utils.ACCESSTOKEN_MAX_AGE), "/", "", true, true)
		ctx.SetCookie(utils.REFRESH_TOKEN, refreshToken, int(utils.REFRESHTOKEN_MAX_AGE), "/", "", true, true)

		ctx.String(http.StatusOK, "refresh success")
	}
}
