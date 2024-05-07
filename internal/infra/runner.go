package infra

import (
	"github.com/gin-gonic/gin"

	"github.com/icrxz/crm-api-core/internal/infra/entrypoint"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/rest"
)

func RunApp() error {
	// controllers
	pingController := rest.NewPingController()
	userController := rest.NewUserController()

	// routes
	router := gin.Default()
	entrypoint.LoadRoutes(router, pingController, userController)

	return router.Run()
}
