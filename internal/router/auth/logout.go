package auth

import (
	"kmem/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout() func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.SetCookie(string(utils.ACCESS_TOKEN), "", -1, "/", "", true, true)
		ctx.SetCookie(string(utils.REFRESH_TOKEN), "", -1, "/", "", true, true)
		ctx.String(http.StatusOK, "ok")
	}
}
