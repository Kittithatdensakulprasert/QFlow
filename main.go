package main

import (
	"log"
	"net/http"
	"qflow/handlers"
)

func main() {

	// Auth routes
	http.HandleFunc("/api/auth/request-otp", handlers.RequestOTP)
	http.HandleFunc("/api/auth/verify-otp", handlers.VerifyOTP)
	http.HandleFunc("/api/auth/register", handlers.Register)
	http.HandleFunc("/api/auth/me", handlers.ProfileHandler)

	// Queue routes
	http.HandleFunc("/api/queues/book", handlers.BookQueue)
	http.HandleFunc("/api/queues/history", handlers.GetHistory)
	http.HandleFunc("/api/queues/", handlers.QueueHandler)

	http.HandleFunc("/api/manage/queues/", handlers.ManageHandler)

	http.HandleFunc("/api/notifications", handlers.GetNotifications)
	http.HandleFunc("/api/notifications/", handlers.NotificationHandler)

	log.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
