package main

import (
	"log"

	"qflow/config"
	"qflow/db"
	"qflow/internal/domain"
	"qflow/internal/jwt"
	"qflow/internal/repository"
	"qflow/internal/router"
	"qflow/internal/service"
	"qflow/internal/swagger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if cfg.JWTSecret == "" ||
		cfg.JWTSecret == "secret" ||
		cfg.JWTSecret == "your-secret-key-here" ||
		cfg.JWTSecret == "change-me-to-a-long-random-jwt-secret-for-local-dev" {
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
	seedBootstrapUser(authRepo, cfg.BootstrapAdminPhone, cfg.BootstrapAdminName, "admin")
	seedBootstrapUser(authRepo, cfg.BootstrapProviderPhone, cfg.BootstrapProviderName, "provider")

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)
	authSvc := service.NewAuthService(authRepo, jwtManager)

	r := gin.Default()
	router.Setup(r, providerSvc, queueSvc, notificationSvc, authSvc, jwtManager, cfg.ExposeOTPInResponse())
	swagger.Register(r)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func seedBootstrapUser(authRepo domain.AuthRepository, phone, name, role string) {
	if phone == "" {
		return
	}

	user, err := authRepo.FindUserByPhone(phone)
	if err != nil {
		user = &domain.User{
			Phone: phone,
			Name:  name,
			Role:  role,
		}
		if err := authRepo.CreateUser(user); err != nil {
			log.Fatalf("failed to create bootstrap %s user: %v", role, err)
		}
		return
	}

	user.Name = name
	user.Role = role
	if err := authRepo.UpdateUser(user); err != nil {
		log.Fatalf("failed to update bootstrap %s user: %v", role, err)
	}
}
