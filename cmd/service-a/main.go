package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"text-analysis-test-task/pkg/servicea"
)

func main() {
	serviceBURL := os.Getenv("SERVICE_B_URL")
	if serviceBURL == "" {
		serviceBURL = "http://localhost:8081"
	}

	storage := servicea.NewStorage()
	client := servicea.NewClient(serviceBURL)
	handlers := servicea.NewHandlers(storage, client)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/text", handlers.HandlePostText)
	mux.HandleFunc("GET /api/v1/status/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGetStatus(w, r, r.PathValue("id"))
	})
	mux.HandleFunc("GET /api/v1/health", handlers.HandleHealth)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Service A starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("Service A stopped")
}

