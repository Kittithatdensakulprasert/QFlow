package main

import (
	"log"
	"net/http"
	"qflow/handlers"
)

func main() {

	http.HandleFunc("/api/queues/book", handlers.BookQueue)
	http.HandleFunc("/api/queues/history", handlers.GetHistory)
	http.HandleFunc("/api/queues/", handlers.QueueHandler)

	http.HandleFunc("/api/manage/queues/", handlers.ManageHandler)

	http.HandleFunc("/api/notifications", handlers.GetNotifications)
	http.HandleFunc("/api/notifications/", handlers.NotificationHandler)

	//Category  routes
	http.HandleFunc("/api/categories", handlers.GetCategories)
	http.HandleFunc("/api/categories/", handlers.CategoryHandler)

	log.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))

}
