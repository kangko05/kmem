package router

import (
	"fmt"
	"kmem/internal/event"
	"kmem/internal/models"
	"kmem/internal/utils"
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

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("failed to hash password: %v", err))
			return
		}

		// get result from store
		resultCh := make(chan event.Result, 1)
		defer close(resultCh)

		store.Register(event.UserRegistered(
			models.User{Username: user.Username, Password: hashedPassword},
			event.WithResultChan(resultCh)),
		)

		result := <-resultCh

		if result.Status() == event.SUCCESS {
			ctx.String(http.StatusOK, result.Message())
		} else {
			ctx.String(http.StatusInternalServerError, result.Message())
		}
	}
}
