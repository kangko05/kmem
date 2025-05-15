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

func Login(store *event.Store, pg *database.Postgres) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var user models.User
		checkUserPostBody(ctx, &user)

		// 1. check user from db
		quser, err := query.QueryUser(pg, user.Username)
		if err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("error occurred while querying user: %v", err))
			return
		}

		// 2. check password
		ok := utils.CheckPasswordHash(user.Password, quser.Password)
		if !ok {
			ctx.String(http.StatusBadRequest, "wrong username or password")
			return
		}

		// 3. create jwt and store it in cookie - from store - register user logged in event
		resChan := make(chan event.Result, 1)

		store.Register(event.UserLoggedIn(
			user.Username,
			event.WithResultChan(resChan),
		))

		result := <-resChan

		if result.Status() == event.SUCCESS {
			tokens, ok := result.Payload().(map[string]string)
			if tokens == nil {
				ctx.String(http.StatusInternalServerError, "tokens not found")
				return
			}
			if !ok {
				ctx.String(http.StatusInternalServerError, "invalid tokens type")
				return
			}

			accessToken, atExists := tokens[utils.ACCESS_TOKEN]
			if !atExists {
				ctx.String(http.StatusInternalServerError, "accessToken does not exist")
				return
			}

			refreshToken, rtExists := tokens[utils.REFRESH_TOKEN]
			if !rtExists {
				ctx.String(http.StatusInternalServerError, "refreshToken does not exist")
				return
			}

			ctx.SetCookie(utils.ACCESS_TOKEN, accessToken, int(utils.ACCESSTOKEN_MAX_AGE), "/", "", true, true)
			ctx.SetCookie(utils.REFRESH_TOKEN, refreshToken, int(utils.REFRESHTOKEN_MAX_AGE), "/", "", true, true)

			ctx.String(http.StatusOK, result.Message())
		} else {
			ctx.String(http.StatusInternalServerError, result.Message())
		}
	}
}
