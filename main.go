package main

import (
	"qflow/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api := r.Group("/api/queues")
	{
		api.POST("/book", handlers.BookQueue)
		api.GET("/history", handlers.GetHistory) // ต้องอยู่ก่อน /:queueNumber
		api.GET("/:queueNumber", handlers.GetQueue)
		api.PATCH("/:id/cancel", handlers.CancelQueue)
	}

	r.Run(":8080")
}
