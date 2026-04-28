package main

import (
	"qflow/config"
	"qflow/db"
	"qflow/internal/jwt"
	"qflow/internal/repository"
	"qflow/internal/router"
	"qflow/internal/service"

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

	authRepo := repository.NewAuthRepository(database)
	jwtManager := jwt.NewJWTManager("your-secret-key-here") // Use environment variable in production
	authSvc := service.NewAuthService(authRepo, jwtManager)

	r := gin.Default()
	router.Setup(r, providerSvc, queueSvc, notificationSvc, authSvc, jwtManager)

	r.Run(":" + cfg.Port)
}
