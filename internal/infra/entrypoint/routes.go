package entrypoint

import (
	"github.com/gin-gonic/gin"

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
) {
	group := app.Group("/crm/core/api/v1")

	// miscellaneous
	app.GET("/ping", pingController.Pong)

	// user
	group.POST("/users", userController.CreateUser)
	group.GET("/users", userController.SearchUser)
	group.GET("/users/:userID", userController.GetUser)
	group.PUT("/users/:userID", userController.UpdateUser)
	group.DELETE("/users/:userID", userController.DeleteUser)

	// partner
	group.POST("/partners", partnerController.CreatePartner)
	group.GET("/partners", partnerController.SearchPartners)
	group.GET("/partners/:partnerID", partnerController.GetPartner)
	group.PUT("/partners/:partnerID", partnerController.UpdatePartner)
	group.DELETE("/partners/:partnerID", partnerController.DeletePartner)

	// customers
	group.POST("/customers", customerController.CreateCustomer)
	group.GET("/customers", customerController.SearchCustomers)
	group.GET("/customers/:customerID", customerController.GetCustomer)
	group.PUT("/customers/:customerID", customerController.UpdateCustomer)
	group.DELETE("/customers/:customerID", customerController.DeleteCustomer)

	// contractors
	group.POST("/contractors", contractorController.CreateContractor)
	group.GET("/contractors", contractorController.SearchContractors)
	group.GET("/contractors/:contractorID", contractorController.GetContractor)
	group.PUT("/contractors/:contractorID", contractorController.UpdateContractor)
	group.DELETE("/contractors/:contractorID", contractorController.DeleteContractor)

	// auth
	group.POST("/login", authController.Login)
	group.POST("/logout", authController.Logout)

	// webMessage
	group.POST("/web/message", webMessageController.ReceiveMessage)
}
