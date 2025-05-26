package handlers

import (
	"fmt"
	"kmem/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		if err := ctx.Bind(&user); err != nil {
			ctx.String(http.StatusBadRequest, "failed to receive user info: %v", err)
			return
		}

		// insert user data into db
		fmt.Println(user)

		ctx.String(http.StatusOK, "ok")
	}
}
