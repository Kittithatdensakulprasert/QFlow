package main

import (
	"qflow/config"
	"qflow/db"
	"qflow/internal/repository"
	"qflow/internal/router"
	"qflow/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DSN)
	db.Migrate(database)

	notificationRepo := repository.NewNotificationRepository(database)
	notificationSvc := service.NewNotificationService(notificationRepo)

	authRepo := repository.NewAuthRepository(database)
	authSvc := service.NewAuthService(authRepo)

	r := gin.Default()
	router.Setup(r, notificationSvc, authSvc)

	r.Run(":" + cfg.Port)
}
