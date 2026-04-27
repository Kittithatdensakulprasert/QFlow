package router

import (
	"qflow/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	auth := handler.NewAuthHandler()
	category := handler.NewCategoryHandler()
	provider := handler.NewProviderHandler()
	queue := handler.NewQueueHandler()
	notification := handler.NewNotificationHandler()

	api := r.Group("/api")

	// Auth
	api.POST("/auth/request-otp", auth.RequestOTP)
	api.POST("/auth/verify-otp", auth.VerifyOTP)
	api.POST("/auth/register", auth.Register)
	api.GET("/auth/me", auth.GetProfile)
	api.PUT("/auth/me", auth.UpdateProfile)

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

	// Notification
	api.GET("/notifications", notification.GetNotifications)
	api.POST("/notifications/send", notification.SendNotification)
	api.PATCH("/notifications/:id/read", notification.MarkNotificationRead)
	api.DELETE("/notifications/:id", notification.DeleteNotification)
}
