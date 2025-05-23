package auth

import (
	"kmem/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: store refresh tokens into db and trace it for better security
func Refresh(secretKey string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		refreshToken, err := ctx.Cookie(string(utils.REFRESH_TOKEN))
		if err != nil {
			ctx.String(http.StatusBadRequest, "refresh token not found")
			return
		}

		_, claims, err := utils.PasrseJwt(secretKey, refreshToken)
		if err != nil {
			ctx.String(http.StatusBadRequest, "invalid refresh token")
			return
		}

		username, ok := claims[utils.USERNAME_KEY].(string)
		if !ok {
			ctx.String(http.StatusBadRequest, "invalid refresh token claim")
			return
		}

		accessToken, err := utils.CreateJwt(utils.ACCESSTOKEN_MAX_AGE, username, secretKey)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "failed to create token")
			return
		}

		refreshToken, err = utils.CreateJwt(utils.REFRESHTOKEN_MAX_AGE, username, secretKey)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "failed to create token")
			return
		}

		ctx.SetCookie(string(utils.ACCESS_TOKEN), accessToken, int(utils.ACCESSTOKEN_MAX_AGE), "/", "", true, true)
		ctx.SetCookie(string(utils.REFRESH_TOKEN), refreshToken, int(utils.REFRESHTOKEN_MAX_AGE), "/", "", true, true)

		ctx.String(http.StatusOK, "refresh success")
	}
}
