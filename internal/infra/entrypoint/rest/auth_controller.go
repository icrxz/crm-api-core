package rest

import (
	"net/http"

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
	var credentials *CredentialsDTO
	err := ctx.BindJSON(&credentials)
	if err != nil {
		ctx.Error(err)
		return
	}

	token, user, err := c.authService.Login(ctx.Request.Context(), credentials.Email, credentials.Password, ctx.ClientIP())
	if err != nil {
		ctx.Error(err)
		return
	}

	authResponseDTO := mapUserToAuthResponseDTO(token, *user)

	ctx.JSON(http.StatusOK, authResponseDTO)
}

func (c *AuthController) Logout(ctx *gin.Context) {
}
