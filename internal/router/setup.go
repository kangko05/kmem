package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {

	router := gin.Default()

	router.GET("/ping", ping)

	return router
}

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}
