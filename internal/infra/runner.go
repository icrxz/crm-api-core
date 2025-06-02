package infra

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/infra/config"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/middleware"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/rest"
	"github.com/icrxz/crm-api-core/internal/infra/repository/bucket"
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

	// bucket
	s3Client, err := bucket.NewS3Bucket(context.Background(), appConfig.AttachmentsBucket)
	if err != nil {
		return err
	}

	attachmentBucket := bucket.NewAttachmentBucket(s3Client, appConfig.AttachmentsBucket.Name)

	// repositories
	userRepository := database.NewUserRepository(sqlDB)
	partnerRepository := database.NewPartnerRepository(sqlDB)
	customerRepository := database.NewCustomerRepository(sqlDB)
	contractorRepository := database.NewContractorRepository(sqlDB)
	caseRepository := database.NewCaseRepository(sqlDB)
	productRepository := database.NewProductRepository(sqlDB)
	commentRepository := database.NewCommentRepository(sqlDB)
	transactionRepository := database.NewTransactionRepository(sqlDB)
	attachmentRepository := database.NewAttachmentRepository(sqlDB)

	// services
	userService := application.NewUserService(userRepository)
	partnerService := application.NewPartnerService(partnerRepository)
	customerService := application.NewCustomerService(customerRepository)
	contractorService := application.NewContractorService(contractorRepository)
	authService := application.NewAuthService(userRepository, appConfig.SecretKey())
	productService := application.NewProductService(productRepository)
	batchCaseService := application.NewBatchCaseService(customerService, productService, contractorService, caseRepository)
	commentService := application.NewCommentService(commentRepository, attachmentRepository, attachmentBucket)
	transactionService := application.NewTransactionService(transactionRepository, caseRepository)
	caseService := application.NewCaseService(
		customerService,
		caseRepository,
		productService,
		userService,
		commentService,
		transactionService,
		partnerService,
		contractorService,
	)
	reportService := application.NewReportService(
		appConfig.ReportFolder,
		caseService,
		productService,
		customerService,
		commentService,
		partnerService,
		contractorService,
		attachmentBucket,
	)
	attachmentService := application.NewAttachmentService(attachmentRepository, attachmentBucket)
	caseActionService := application.NewCaseActionService(caseRepository, commentService, reportService, attachmentService, transactionService)

	// controllers
	pingController := rest.NewPingController()
	userController := rest.NewUserController(userService)
	partnerController := rest.NewPartnerController(partnerService)
	customerController := rest.NewCustomerController(customerService)
	contractorController := rest.NewContractorController(contractorService)
	webMessageController := rest.NewWebMessageController()
	authController := rest.NewAuthController(authService)
	caseController := rest.NewCaseController(caseService, batchCaseService)
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
