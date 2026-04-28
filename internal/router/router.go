package router

import (
	"qflow/internal/domain"
	"qflow/internal/handler"
	"qflow/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, queueSvc domain.QueueService, notificationSvc domain.NotificationService, jwtSecret string) {
	auth := handler.NewAuthHandler()
	category := handler.NewCategoryHandler()
	provider := handler.NewProviderHandler()
	queue := handler.NewQueueHandler(queueSvc)
	notification := handler.NewNotificationHandler(notificationSvc)

	api := r.Group("/api")

	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(jwtSecret))

	protected.GET("/auth/me", auth.GetProfile)
	protected.PUT("/auth/me", auth.UpdateProfile)

	protected.POST("/queues/book", queue.BookQueue)
	protected.GET("/queues/history", queue.GetHistory)
	protected.GET("/queues/:queueNumber", queue.GetQueue)
	protected.PATCH("/queues/:id/cancel", queue.CancelQueue)

	protected.GET("/notifications", notification.GetNotifications)
	protected.PATCH("/notifications/:id/read", notification.MarkNotificationRead)
	protected.DELETE("/notifications/:id", notification.DeleteNotification)

	// Auth
	api.POST("/auth/request-otp", auth.RequestOTP)
	api.POST("/auth/verify-otp", auth.VerifyOTP)
	api.POST("/auth/register", auth.Register)

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

	// Queue Management
	api.GET("/manage/queues/:zoneId", queue.GetQueuesByZone)
	api.PATCH("/manage/queues/:id/call", queue.CallQueue)
	api.PATCH("/manage/queues/:id/complete", queue.CompleteQueue)
	api.PATCH("/manage/queues/:id/skip", queue.SkipQueue)

	// Notification
	api.POST("/notifications/send", notification.SendNotification)
}
