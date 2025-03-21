package main

import (
	"log"
	"net/http"

	"github.com/hardikkum444/go-serverless/handlers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	reg := prometheus.NewRegistry()
	m := NewMatrics(reg)

	m.devices.Set(float64(len(Dvs)))

	// promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	http.HandleFunc("/api/submit", handlers.LoggingMiddleware(handlers.SubmitHandler))
	http.HandleFunc("/api/execute", handlers.LoggingMiddleware(handlers.ExecuteHandler))

	http.HandleFunc("/server/devices", GetDevice)
	http.Handle("/metrics", promhttp.Handler())

	log.Println("serverless system running on :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}
}
