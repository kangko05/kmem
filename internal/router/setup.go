package router

import (
	"kmem/internal/database"
	"kmem/internal/event"
	"kmem/internal/router/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(store *event.Store, pg *database.Postgres) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", ping)

	setupAuth(router, store, pg)

	return router
}

func setupAuth(router *gin.Engine, store *event.Store, pg *database.Postgres) *gin.RouterGroup {
	authGroup := router.Group("auth")
	{
		authGroup.POST("signup", auth.Signup(store, pg))
		authGroup.POST("login", auth.Login(store, pg))
		authGroup.GET("refresh", auth.Refresh())
	}

	return authGroup
}

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}
