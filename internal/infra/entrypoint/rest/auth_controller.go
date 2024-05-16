package rest

import (
	"github.com/gin-gonic/gin"
)

type AuthController struct {
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (c *AuthController) Login(ctx *gin.Context) {
}

func (c *AuthController) Logout(ctx *gin.Context) {
}
