package router

import (
	"kmem/internal/handlers"

	"github.com/gin-gonic/gin"
)

func setupAuth(router *gin.Engine) {
	gr := router.Group("auth")
	{
		gr.POST("signup", handlers.Signup())
	}
}
