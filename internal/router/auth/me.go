package auth

import (
	"fmt"
	"kmem/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// check if user is logged in currently

func Me() func(*gin.Context) {
	return func(ctx *gin.Context) {
		username, exists := ctx.Get(utils.USERNAME_KEY)
		if !exists {
			ctx.String(http.StatusUnauthorized, "not authenticated")
			return
		}

		ctx.String(http.StatusOK, fmt.Sprintf("user %s exists", username))
	}
}
