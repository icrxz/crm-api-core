package rest

import (
	"fmt"
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
	userID := ctx.GetString("user_id")
	err := c.authService.Logout(ctx.Request.Context(), userID)
	if err != nil {
		if ctxErr := ctx.Error(err); ctxErr != nil {
			fmt.Printf("failed to add error in context: %v\n", ctxErr.Error())
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
