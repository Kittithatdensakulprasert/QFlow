package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	blocklist := []string{"", "secret", "your-secret-key-here", "change-me-to-a-long-random-jwt-secret-for-local-dev"} //nolint:gosec // G101: blocklist of weak values, not actual credentials
	isWeak := false
	for _, v := range blocklist {
		if cfg.JWTSecret == v {
			isWeak = true
			break
		}
	}
	if isWeak {
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
	categoryRepo := repository.NewCategoryGormRepository(database)
	categorySvc := service.NewCategoryService(categoryRepo)
	authRepo := repository.NewAuthRepository(database)
	seedBootstrapUser(authRepo, cfg.BootstrapAdminPhone, cfg.BootstrapAdminName, "admin")
	seedBootstrapUser(authRepo, cfg.BootstrapProviderPhone, cfg.BootstrapProviderName, "provider")

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)
	authSvc := service.NewAuthService(authRepo, jwtManager)
	otpCleanupJob := service.NewOTPCleanupJob(authRepo, cfg.ParsedOTPCleanupInterval(), nil)
	otpCleanupJob.Start(appCtx)

	r := gin.Default()
	router.Setup(r, providerSvc, queueSvc, notificationSvc, authSvc, categorySvc, jwtManager, cfg.ExposeOTPInResponse())
	swagger.Register(r)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		<-appCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
}

func seedBootstrapUser(authRepo domain.AuthRepository, phone, name, role string) {
	if phone == "" {
		return
	}

	user, err := authRepo.FindUserByPhone(context.Background(), phone)
	if err != nil {
		user = &domain.User{
			Phone: phone,
			Name:  name,
			Role:  role,
		}
		if err := authRepo.CreateUser(context.Background(), user); err != nil {
			log.Fatalf("failed to create bootstrap %s user: %v", role, err)
		}
		return
	}

	user.Name = name
	user.Role = role
	if err := authRepo.UpdateUser(context.Background(), user); err != nil {
		log.Fatalf("failed to update bootstrap %s user: %v", role, err)
	}
}
