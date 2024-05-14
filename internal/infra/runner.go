package infra

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/infra/config"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/rest"
	"github.com/icrxz/crm-api-core/internal/infra/repositories/database"
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

	// services
	userService := application.NewUserService(userRepository)

	// controllers
	pingController := rest.NewPingController()
	userController := rest.NewUserController(userService)
	webMessageController := rest.NewWebMessageController()

	router := gin.Default()
	router.Use(entrypoint.CustomErrorEncoder())

	entrypoint.LoadRoutes(router, pingController, userController, webMessageController)

	return router.Run()
}
