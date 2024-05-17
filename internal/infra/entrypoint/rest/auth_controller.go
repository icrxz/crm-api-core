package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
)

type AuthController struct {
	authService application.AuthService
}

func NewAuthController(authService application.AuthService) AuthController {
	return AuthController{
		authService: authService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
}

func (c *AuthController) Logout(ctx *gin.Context) {
}
