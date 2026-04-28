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
) {
	auth := handler.NewAuthHandler(authSvc)
	category := handler.NewCategoryHandler()
	provider := handler.NewProviderHandler(providerSvc)
	queue := handler.NewQueueHandler(queueSvc)
	notification := handler.NewNotificationHandler(notificationSvc)

	api := r.Group("/api")

	// Public Auth routes
	api.POST("/auth/request-otp", auth.RequestOTP)
	api.POST("/auth/verify-otp", auth.VerifyOTP)
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

	// Category
	api.GET("/categories", category.GetCategories)
	api.GET("/categories/:id", category.GetCategory)
	api.POST("/categories", category.CreateCategory)
	api.PUT("/categories/:id", category.UpdateCategory)
	api.DELETE("/categories/:id", category.DeleteCategory)

	// Provider & Zone
	api.POST("/providers", provider.CreateProvider)
	api.GET("/providers", provider.GetProviders)
	api.POST("/providers/:id/zones", provider.CreateZone)
	api.GET("/providers/:id/zones", provider.GetZones)
	api.PATCH("/zones/:id/toggle", provider.ToggleZone)

	// Queue Booking
	api.POST("/queues/book", queue.BookQueue)
	api.GET("/queues/history", queue.GetHistory)
	api.GET("/queues/:queueNumber", queue.GetQueue)
	api.PATCH("/queues/:id/cancel", queue.CancelQueue)

	// Queue Management
	api.GET("/manage/queues/:zoneId", queue.GetQueuesByZone)
	api.PATCH("/manage/queues/:id/call", queue.CallQueue)
	api.PATCH("/manage/queues/:id/complete", queue.CompleteQueue)
	api.PATCH("/manage/queues/:id/skip", queue.SkipQueue)
}
