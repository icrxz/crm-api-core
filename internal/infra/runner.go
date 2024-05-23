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

	// services
	userService := application.NewUserService(userRepository)
	partnerService := application.NewPartnerService(partnerRepository)
	customerService := application.NewCustomerService(customerRepository)
	contractorService := application.NewContractorService(contractorRepository)
	authService := application.NewAuthService(userRepository, appConfig.SecretJWTKey)

	// controllers
	pingController := rest.NewPingController()
	userController := rest.NewUserController(userService)
	partnerController := rest.NewPartnerController(partnerService)
	customerController := rest.NewCustomerController(customerService)
	contractorController := rest.NewContractorController(contractorService)
	webMessageController := rest.NewWebMessageController()
	authController := rest.NewAuthController(authService)

	// middlewares
	authMiddleware := middleware.NewAuthenticationMiddleware(authService)

	router := gin.Default()
	router.Use(entrypoint.CustomErrorEncoder())

	entrypoint.LoadRoutes(router, pingController, userController, webMessageController, partnerController, customerController, contractorController, authController, authMiddleware)

	return router.Run()
}
