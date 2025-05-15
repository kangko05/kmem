package router

import (
	"kmem/internal/event"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(store *event.Store) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", ping)

	auth := router.Group("auth")
	{
		auth.POST("signup", signup(store))
	}

	return router
}

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}
