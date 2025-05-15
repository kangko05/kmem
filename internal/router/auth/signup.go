package auth

import (
	"fmt"
	"kmem/internal/database"
	"kmem/internal/database/query"
	"kmem/internal/event"
	"kmem/internal/models"
	"kmem/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Signup(store *event.Store, pg *database.Postgres) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var user models.User
		checkUserPostBody(ctx, &user)

		// check if username exists
		quser, _ := query.QueryUser(pg, user.Username)
		if quser.Username == user.Username && len(quser.Password) > 0 {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("username %s exists", quser.Username))
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("failed to hash password: %v", err))
			return
		}

		// get result from store
		resChan := make(chan event.Result, 1)
		defer close(resChan)

		store.Register(event.UserRegistered(
			models.User{Username: user.Username, Password: hashedPassword},
			event.WithResultChan(resChan)),
		)

		result := <-resChan

		if result.Status() == event.SUCCESS {
			ctx.String(http.StatusOK, result.Message())
		} else {
			ctx.String(http.StatusInternalServerError, result.Message())
		}
	}
}
