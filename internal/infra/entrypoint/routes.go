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
	caseController rest.CaseController,
	productController rest.ProductController,
	commentController rest.CommentController,
	transactionController rest.TransactionController,
	caseActionController rest.CaseActionController,
) {
	authGroup := app.Group("/crm/core/api/v1")
	authGroup.Use(authMiddleware.Authenticate())

	publicGroup := app.Group("/crm/core/api/v1")

	// miscellaneous
	app.GET("/ping", pingController.Pong)

	// user
	publicGroup.POST("/users", userController.CreateUser)
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
	authGroup.POST("/partners/batch", partnerController.CreateBatch)

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

	// cases
	authGroup.POST("/cases", caseController.CreateCase)
	authGroup.GET("/cases/:caseID", caseController.GetCase)
	authGroup.PUT("/cases/:caseID", caseController.UpdateCase)
	authGroup.GET("/cases", caseController.SearchCases)
	authGroup.GET("/cases/full", caseController.SearchCasesFull)
	authGroup.POST("/cases/batch", caseController.CreateBatch)
	authGroup.GET("/cases/:caseID/full", caseController.GetCaseFull)

	// products
	authGroup.GET("/products/:productID", productController.GetProductByID)
	authGroup.POST("/products", productController.CreateProduct)
	authGroup.PUT("/products/:productID", productController.UpdateProduct)

	// comments
	authGroup.GET("/comments/:commentID", commentController.GetByID)
	authGroup.POST("/cases/:caseID/comments", commentController.CreateComment)
	authGroup.GET("/cases/:caseID/comments", commentController.GetByCaseID)

	// transactions
	authGroup.POST("/cases/:caseID/transactions", transactionController.CreateTransaction)
	authGroup.GET("/transactions/:transactionID", transactionController.GetTransaction)
	authGroup.PUT("/transactions/:transactionID", transactionController.UpdateTransaction)
	authGroup.GET("/transactions", transactionController.SearchTransactions)
	authGroup.POST("/cases/:caseID/transactions/batch", transactionController.CreateTransactionBatch)

	// case actions
	authGroup.PATCH("/cases/:caseID/owner", caseActionController.ChangeOwner)
	authGroup.PATCH("/cases/:caseID/status", caseActionController.ChangeStatus)
	authGroup.PATCH("/cases/:caseID/partner", caseActionController.ChangePartner)
	authGroup.GET("/cases/:caseID/report", caseActionController.DownloadReport)
}
