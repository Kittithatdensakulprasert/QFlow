package main

import (
	"log"
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
	if cfg.JWTSecret == "" || cfg.JWTSecret == "secret" || cfg.JWTSecret == "your-secret-key-here" {
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
	authRepo := repository.NewAuthRepository(database)

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)
	authSvc := service.NewAuthService(authRepo, jwtManager)

	r := gin.Default()
	router.Setup(r, providerSvc, queueSvc, notificationSvc, authSvc, jwtManager, cfg.ExposeOTPInResponse())

	r.Run(":" + cfg.Port)
}
