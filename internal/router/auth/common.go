package auth

import (
	"fmt"
	"kmem/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkUserPostBody(ctx *gin.Context, user *models.User) {
	if err := ctx.Bind(user); err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("failed to get post body: %v", err))
		return
	}

	if len(user.Username) < 4 || len(user.Password) < 8 {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("username must be longer than 4, password must be longer than 8"))
		return
	}
}
