package main

import (
	"log"
	"net/http"

	"github.com/hardikkum444/go-serverless/handlers"
)

func main() {
	http.HandleFunc("/api/submit", handlers.LoggingMiddleware(handlers.SubmitHandler))
	http.HandleFunc("/api/execute", handlers.LoggingMiddleware(handlers.ExecuteHandler))

	log.Println("serverless system running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}
}
