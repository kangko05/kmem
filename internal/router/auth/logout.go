package auth

import (
	"kmem/internal/event"
	"kmem/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(store *event.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		val, exists := ctx.Get("username")
		if !exists {
			ctx.String(http.StatusUnauthorized, "user not found")
			return
		}

		username, ok := val.(string)
		if !ok {
			ctx.String(http.StatusBadRequest, "user not found")
		}

		rchan := make(chan event.Result, 1)
		defer close(rchan)

		store.Register(event.UserLoggedOut(username, event.WithResultChan(rchan)))
		result := <-rchan

		// try removing tokens no matter what
		ctx.SetCookie(string(utils.ACCESS_TOKEN), "", -1, "/", "", true, true)
		ctx.SetCookie(string(utils.REFRESH_TOKEN), "", -1, "/", "", true, true)

		if result.Status() == utils.SUCCESS {
			ctx.String(http.StatusOK, "ok")
		} else {
			ctx.String(http.StatusInternalServerError, "something went wrong")
		}
	}
}
