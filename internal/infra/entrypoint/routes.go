package entrypoint

import (
	"github.com/gin-gonic/gin"

	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/rest"
)

func LoadRoutes(app *gin.Engine, pingController rest.PingController, userController rest.UserController, webMessageController rest.WebMessageController) {
	group := app.Group("/crm/core/v1")

	// miscellaneous
	app.GET("/ping", pingController.Pong)

	// user
	group.POST("/users", userController.CreateUser)
	group.GET("/users", userController.SearchUser)
	group.GET("/users/:userID", userController.GetUser)
	group.PUT("/users/:userID", userController.UpdateUser)
	group.DELETE("/users/:userID", userController.DeleteUser)

	// webMessage
	group.POST("/web/message", webMessageController.ReceiveMessage)
}
