package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingController struct{}

func NewPingController() PingController {
	return PingController{}
}

func (c *PingController) Pong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{"result": "pong"})
}
