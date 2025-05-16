package router

import (
	"kmem/internal/database"
	"kmem/internal/event"
	"kmem/internal/router/auth"
	"kmem/internal/router/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(store *event.Store, pg *database.Postgres) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"}, // TODO: need to set this properly on production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ping", ping)

	setupAuth(router, store, pg)

	return router
}

func setupAuth(router *gin.Engine, store *event.Store, pg *database.Postgres) *gin.RouterGroup {
	authGroup := router.Group("auth")
	{
		authGroup.POST("signup", auth.Signup(store, pg))
		authGroup.POST("login", auth.Login(store, pg))
		authGroup.GET("logout", auth.Logout())
		authGroup.GET("me", middleware.AuthenticateJwt(), auth.Me())
	}

	return authGroup
}

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}
