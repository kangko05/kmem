package router

import (
	"kmem/internal/config"
	"kmem/internal/db"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(pg *db.Postgres, conf *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("ping", ping) // for test & health check

	setupAuth(router, pg, conf)
	setupFiles(router, conf)

	return router
}

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong\n")
}

func setupAuth(router *gin.Engine, pg *db.Postgres, conf *config.Config) {
	gr := router.Group("auth")
	{
		gr.POST("signup", signup(pg))
		gr.POST("login", login(pg, conf))
		gr.GET("logout", logout())
		gr.GET("me", authMiddleware(conf), me())
	}
}

func setupFiles(router *gin.Engine, conf *config.Config) {
	gr := router.Group("files")
	// gr.Use(authMiddleware(conf))
	{
		gr.POST("upload", upload(conf))
	}
}
