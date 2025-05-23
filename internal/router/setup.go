package router

import (
	"kmem/internal/config"
	"kmem/internal/database"
	"kmem/internal/event"
	"kmem/internal/router/auth"
	"kmem/internal/router/files"
	"kmem/internal/router/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(store *event.Store, conf *config.Config, pg *database.Postgres, cache *database.Cache, secretKey string) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"}, // TODO: need to set this properly on production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    true,
	}))

	router.GET("ping", ping)

	setupAuth(router, store, pg, secretKey)
	setupFiles(router, conf, store, pg, cache, secretKey)

	return router
}

func setupAuth(router *gin.Engine, store *event.Store, pg *database.Postgres, secretKey string) *gin.RouterGroup {
	authGroup := router.Group("auth")
	{
		authGroup.POST("signup", auth.Signup(store, pg))
		authGroup.POST("login", auth.Login(store, pg, secretKey))
		authGroup.GET("logout", middleware.AuthenticateJwt(secretKey), auth.Logout(store))
		authGroup.GET("me", middleware.AuthenticateJwt(secretKey), auth.Me())
	}

	return authGroup
}

func setupFiles(router *gin.Engine, conf *config.Config, store *event.Store, pg *database.Postgres, cache *database.Cache, secretKey string) {
	filesGroup := router.Group("files")
	filesGroup.Use(middleware.AuthenticateJwt(secretKey))
	{
		filesGroup.POST("upload", files.Upload(conf, store))
		filesGroup.GET("items", files.GetItems(pg, cache))
		filesGroup.GET("search", files.Search(store, pg, cache))

		filesGroup.Static("static", "/home/kang/Downloads")
	}
}

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}
