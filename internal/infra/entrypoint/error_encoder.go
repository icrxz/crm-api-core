package entrypoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/domain"
)

func CustomErrorEncoder() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case *domain.CustomError:
				c.AbortWithStatusJSON(e.StatusCode(), e)
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "unexpected error"})
			}
		}
	}
}
