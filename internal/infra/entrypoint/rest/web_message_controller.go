package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebMessageController struct{}

func NewWebMessageController() WebMessageController {
	return WebMessageController{}
}

func (c *WebMessageController) ReceiveMessage(ctx *gin.Context) {
	ctx.String(http.StatusNoContent, "")
}
