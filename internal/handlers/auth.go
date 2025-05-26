package handlers

import (
	"fmt"
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/models"
	"kmem/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkUserPostBody(ctx *gin.Context, user *models.User) error {
	if err := ctx.Bind(&user); err != nil {
		models.APIResponse{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("failed to receive user info: %v", err),
		}.Send(ctx)

		return fmt.Errorf("failed to bind user body")
	}

	if len(user.Username) < 4 {
		models.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "short username",
		}.Send(ctx)

		return fmt.Errorf("short username")
	}

	if len(user.Password) < 8 {
		models.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "short password",
		}.Send(ctx)

		return fmt.Errorf("short password")
	}

	return nil
}

func Signup(pg *db.Postgres) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		// response to client handled already
		if err := checkUserPostBody(ctx, &user); err != nil {
			log.Println(err)
			return
		}

		if err := pg.InsertUser(user); err != nil {
			ctx.String(http.StatusInternalServerError, "failed to insert user")
			return
		}

		ctx.String(http.StatusOK, "ok")
	}
}

func Login(pg *db.Postgres, conf *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := checkUserPostBody(ctx, &user); err != nil {
			log.Println(err)
			return
		}

		// check password
		dbuser, err := pg.QueryUser(user.Username)
		if err != nil {
			models.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "incorrect username or password",
			}.Send(ctx)

			return
		}

		if res := utils.CheckPasswordHash(dbuser.Password, user.Password); !res {
			models.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "incorrect username or password",
			}.Send(ctx)

			return
		}

		accessToken, err := utils.GenTokenString(conf.JwtSecretKey(), user.Username, utils.ACCESS_TOKEN_DUR)
		if err != nil {
			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to gen access token: %v", err),
			}.Send(ctx)

			return
		}

		refreshToken, err := utils.GenTokenString(conf.JwtSecretKey(), user.Username, utils.REFRESH_TOKEN_DUR)
		if err != nil {
			models.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to gen refresh token: %v", err),
			}.Send(ctx)

			return
		}

		ctx.SetCookie(utils.ACCESS_TOKEN_KEY, accessToken, int(utils.ACCESS_TOKEN_DUR), "/", "", true, true)
		ctx.SetCookie(utils.REFRESH_TOKEN_KEY, refreshToken, int(utils.REFRESH_TOKEN_DUR), "/", "", true, true)

		models.APIResponse{
			Status:  http.StatusOK,
			Message: "ok",
		}.Send(ctx)
	}
}

func Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetCookie(utils.ACCESS_TOKEN_KEY, "", -1, "/", "", true, true)
		ctx.SetCookie(utils.REFRESH_TOKEN_KEY, "", -1, "/", "", true, true)
		ctx.String(http.StatusOK, "ok")

		models.APIResponse{
			Status:  http.StatusOK,
			Message: "ok",
		}.Send(ctx)
	}
}
