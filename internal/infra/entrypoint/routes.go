package entrypoint

import (
	"github.com/gin-gonic/gin"

	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/middleware"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/rest"
)

func LoadRoutes(
	app *gin.Engine,
	pingController rest.PingController,
	userController rest.UserController,
	webMessageController rest.WebMessageController,
	partnerController rest.PartnerController,
	customerController rest.CustomerController,
	contractorController rest.ContractorController,
	authController rest.AuthController,
	authMiddleware middleware.AuthenticationMiddleware,
) {
	authGroup := app.Group("/crm/core/api/v1")
	publicGroup := authGroup.Group("")

	// miscellaneous
	app.GET("/ping", pingController.Pong)

	// user
	authGroup.POST("/users", userController.CreateUser)
	authGroup.GET("/users", userController.SearchUser)
	authGroup.GET("/users/:userID", userController.GetUser)
	authGroup.PUT("/users/:userID", userController.UpdateUser)
	authGroup.DELETE("/users/:userID", userController.DeleteUser)

	// partner
	authGroup.POST("/partners", partnerController.CreatePartner)
	authGroup.GET("/partners", partnerController.SearchPartners)
	authGroup.GET("/partners/:partnerID", partnerController.GetPartner)
	authGroup.PUT("/partners/:partnerID", partnerController.UpdatePartner)
	authGroup.DELETE("/partners/:partnerID", partnerController.DeletePartner)

	// customers
	authGroup.POST("/customers", customerController.CreateCustomer)
	authGroup.GET("/customers", customerController.SearchCustomers)
	authGroup.GET("/customers/:customerID", customerController.GetCustomer)
	authGroup.PUT("/customers/:customerID", customerController.UpdateCustomer)
	authGroup.DELETE("/customers/:customerID", customerController.DeleteCustomer)

	// contractors
	authGroup.POST("/contractors", contractorController.CreateContractor)
	authGroup.GET("/contractors", contractorController.SearchContractors)
	authGroup.GET("/contractors/:contractorID", contractorController.GetContractor)
	authGroup.PUT("/contractors/:contractorID", contractorController.UpdateContractor)
	authGroup.DELETE("/contractors/:contractorID", contractorController.DeleteContractor)

	// auth
	publicGroup.POST("/login", authController.Login)
	authGroup.POST("/logout", authController.Logout)

	// webMessage
	publicGroup.POST("/web/message", webMessageController.ReceiveMessage)

	authGroup.Use(authMiddleware.Authenticate())
}
