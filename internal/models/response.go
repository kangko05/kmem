package models

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (r APIResponse) Send(ctx *gin.Context) {
	ctx.JSON(r.Status, r)
}
