package router

import (
	"qflow/internal/domain"
	"qflow/internal/handler"
	"qflow/internal/jwt"
	"qflow/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(
	r *gin.Engine,
	providerSvc domain.ProviderService,
	queueSvc domain.QueueService,
	notificationSvc domain.NotificationService,
	authSvc domain.AuthService,
	jwtManager *jwt.JWTManager,
	exposeOTPResponse bool,
) {
	auth := handler.NewAuthHandler(authSvc, exposeOTPResponse)
	category := handler.NewCategoryHandler()
	provider := handler.NewProviderHandler(providerSvc)
	queue := handler.NewQueueHandler(queueSvc)
	notification := handler.NewNotificationHandler(notificationSvc)

	api := r.Group("/api")

	// Public Auth routes (rate limited)
	api.POST("/auth/request-otp", middleware.OTPRequestLimiter.Middleware(), auth.RequestOTP)
	api.POST("/auth/verify-otp", middleware.OTPVerifyLimiter.Middleware(), auth.VerifyOTP)
	api.POST("/auth/register", auth.Register)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(jwtManager))
	protected.GET("/auth/me", auth.GetProfile)
	protected.PUT("/auth/me", auth.UpdateProfile)
	protected.GET("/notifications", notification.GetNotifications)
	protected.POST("/notifications/send", notification.SendNotification)
	protected.PATCH("/notifications/:id/read", notification.MarkNotificationRead)
	protected.DELETE("/notifications/:id", notification.DeleteNotification)
	protected.POST("/queues/book", queue.BookQueue)
	protected.GET("/queues/history", queue.GetHistory)
	protected.GET("/queues/:queueNumber", queue.GetQueue)
	protected.PATCH("/queues/:id/cancel", queue.CancelQueue)

	// Admin routes
	admin := protected.Group("/")
	admin.Use(middleware.RequireRole("admin"))
	admin.POST("/categories", category.CreateCategory)
	admin.PUT("/categories/:id", category.UpdateCategory)
	admin.DELETE("/categories/:id", category.DeleteCategory)
	admin.POST("/providers", provider.CreateProvider)

	// Provider routes
	providerRoutes := protected.Group("/")
	providerRoutes.Use(middleware.RequireRole("provider"))
	providerRoutes.POST("/providers/:id/zones", provider.CreateZone)
	providerRoutes.PATCH("/zones/:id/toggle", provider.ToggleZone)
	providerRoutes.GET("/manage/queues/:zoneId", queue.GetQueuesByZone)
	providerRoutes.PATCH("/manage/queues/:id/call", queue.CallQueue)
	providerRoutes.PATCH("/manage/queues/:id/complete", queue.CompleteQueue)
	providerRoutes.PATCH("/manage/queues/:id/skip", queue.SkipQueue)

	// Public Category & Provider
	api.GET("/categories", category.GetCategories)
	api.GET("/categories/:id", category.GetCategory)
	api.GET("/providers", provider.GetProviders)
	api.GET("/providers/:id/zones", provider.GetZones)
}
