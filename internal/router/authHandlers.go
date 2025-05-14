package router

import (
	"fmt"
	"kmem/internal/event"
	"kmem/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signup(store *event.Store) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var user models.User

		if err := ctx.Bind(&user); err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("failed to get post body: %v", err))
			return
		}

		if len(user.Username) < 4 || len(user.Password) < 8 {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("username must be longer than 4, password must be longer than 8"))
			return
		}
	}
}
