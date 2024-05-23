package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
)

type AuthenticationMiddleware struct {
	authService application.AuthService
}

func NewAuthenticationMiddleware(authService application.AuthService) AuthenticationMiddleware {
	return AuthenticationMiddleware{authService: authService}
}

func (a *AuthenticationMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authentication header"})
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(tokenString, " ")

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			ctx.Abort()
			return
		}

		tokenString = tokenParts[1]

		claims, err := a.authService.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			ctx.Abort()
			return
		}

		userID := claims["user_id"].(string)
		clientIP := ctx.ClientIP()
		err = a.authService.VerifyUserSession(ctx.Request.Context(), userID, clientIP)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
