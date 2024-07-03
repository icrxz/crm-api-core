package infra

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/infra/config"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/middleware"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/rest"
	"github.com/icrxz/crm-api-core/internal/infra/repository/database"
)

func RunApp() error {
	_ = context.Background()

	appConfig, err := config.Load()
	if err != nil {
		return err
	}

	// database
	sqlDB, err := database.NewDatabase(appConfig.Database)
	if err != nil {
		return err
	}
	defer func() {
		err := sqlDB.Close()
		panic(err)
	}()

	// repositories
	userRepository := database.NewUserRepository(sqlDB)
	partnerRepository := database.NewPartnerRepository(sqlDB)
	customerRepository := database.NewCustomerRepository(sqlDB)
	contractorRepository := database.NewContractorRepository(sqlDB)
	caseRepository := database.NewCaseRepository(sqlDB)
	productRepository := database.NewProductRepository(sqlDB)
	commentRepository := database.NewCommentRepository(sqlDB)
	transactionRepository := database.NewTransactionRepository(sqlDB)

	// services
	userService := application.NewUserService(userRepository)
	partnerService := application.NewPartnerService(partnerRepository)
	customerService := application.NewCustomerService(customerRepository)
	contractorService := application.NewContractorService(contractorRepository)
	authService := application.NewAuthService(userRepository, appConfig.SecretKey())
	productService := application.NewProductService(productRepository)
	caseService := application.NewCaseService(customerService, caseRepository, productService, userService)
	commentService := application.NewCommentService(commentRepository)
	transactionService := application.NewTransactionService(transactionRepository)
	reportService := application.NewReportService(appConfig.ReportFolder, caseService, productService, customerService, commentService, partnerService)
	caseActionService := application.NewCaseActionService(caseRepository, commentService, reportService)

	// controllers
	pingController := rest.NewPingController()
	userController := rest.NewUserController(userService)
	partnerController := rest.NewPartnerController(partnerService)
	customerController := rest.NewCustomerController(customerService)
	contractorController := rest.NewContractorController(contractorService)
	webMessageController := rest.NewWebMessageController()
	authController := rest.NewAuthController(authService)
	caseController := rest.NewCaseController(caseService)
	productController := rest.NewProductController(productService)
	commentController := rest.NewCommentController(commentService)
	transactionController := rest.NewTransactionController(transactionService)
	caseActionController := rest.NewCaseActionController(caseActionService)

	// middlewares
	authMiddleware := middleware.NewAuthenticationMiddleware(authService)

	router := gin.Default()
	router.Use(entrypoint.CustomErrorEncoder())

	entrypoint.LoadRoutes(
		router,
		pingController,
		userController,
		webMessageController,
		partnerController,
		customerController,
		contractorController,
		authController,
		authMiddleware,
		caseController,
		productController,
		commentController,
		transactionController,
		caseActionController,
	)

	return router.Run()
}
