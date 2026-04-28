package main

import (
	"qflow/config"
	"qflow/db"
	"qflow/internal/repository"
	"qflow/internal/router"
	"qflow/internal/service"
	"qflow/internal/swagger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DSN)
	db.Migrate(database)

	providerRepo := repository.NewProviderRepository(database)
	providerSvc := service.NewProviderService(providerRepo)
	queueRepo := repository.NewQueueRepository(database)
	queueSvc := service.NewQueueService(queueRepo)
	notificationRepo := repository.NewNotificationRepository(database)
	notificationSvc := service.NewNotificationService(notificationRepo)

	r := gin.Default()
	router.Setup(r, providerSvc, queueSvc, notificationSvc)
	swagger.Register(r)

	r.Run(":" + cfg.Port)
}
