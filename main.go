package main

import (
	"log"
	"qflow/config"
	"qflow/db"
	"qflow/internal/repository"
	"qflow/internal/router"
	"qflow/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if cfg.JWTSecret == "" || cfg.JWTSecret == "secret" {
		log.Fatal("JWT_SECRET must be set to a strong non-default value")
	}

	database := db.Connect(cfg.DSN)
	db.Migrate(database)

	providerRepo := repository.NewProviderRepository(database)
	providerSvc := service.NewProviderService(providerRepo)
	queueRepo := repository.NewQueueRepository(database)
	queueSvc := service.NewQueueService(queueRepo)
	notificationRepo := repository.NewNotificationRepository(database)
	notificationSvc := service.NewNotificationService(notificationRepo)

	r := gin.Default()
	router.Setup(r, providerSvc, queueSvc, notificationSvc, cfg.JWTSecret)

	r.Run(":" + cfg.Port)
}
