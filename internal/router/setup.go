package router

import (
	"kmem/internal/cache"
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/queue"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(pg *db.Postgres, conf *config.Config, q *queue.Queue, cache *cache.Cache) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://192.168.50.251:5173", "http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("ping", ping) // for test & health check
	router.StaticFS("static", http.Dir(conf.UploadPath()))

	setupAuth(router, pg, conf)
	setupFiles(router, pg, conf, q, cache)

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

func setupFiles(router *gin.Engine, pg *db.Postgres, conf *config.Config, q *queue.Queue, cache *cache.Cache) {
	gr := router.Group("files")
	gr.Use(authMiddleware(conf))
	{
		gr.GET("", servFiles(pg, cache))
		gr.POST("upload", upload(pg, conf, q, cache))
	}
}
