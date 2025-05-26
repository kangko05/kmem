package router

import (
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(pg *db.Postgres, conf *config.Config) *gin.Engine {
	router := gin.Default()

	router.GET("ping", ping) // for test & health check

	setupAuth(router, pg, conf)

	return router
}

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong\n")
}

func setupAuth(router *gin.Engine, pg *db.Postgres, conf *config.Config) {
	gr := router.Group("auth")
	{
		gr.POST("signup", handlers.Signup(pg))
		gr.POST("login", handlers.Login(pg, conf))
		gr.GET("logout", handlers.Logout())
	}
}
