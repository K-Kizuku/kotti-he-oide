package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	domainService "github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/infrastructure/persistence"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// User dependencies
	userRepo := persistence.NewMemoryUserRepository()
	userService := domainService.NewUserService(userRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, userService)

	// Web Push dependencies (using memory repositories for now)
	// TODO: Replace with actual PostgreSQL implementations
	subscriptionRepo := persistence.NewMemoryPushSubscriptionRepository()
	jobRepo := persistence.NewMemoryPushJobRepository()
	logRepo := persistence.NewMemoryPushLogRepository()

	// Initialize VAPID service
	vapidService, err := domainService.NewVAPIDService()
	if err != nil {
		log.Fatal("Failed to initialize VAPID service:", err)
	}

	// Push services
	pushService := domainService.NewPushService(subscriptionRepo, jobRepo)
	pushSenderService := service.NewPushSenderService(subscriptionRepo, jobRepo, logRepo, vapidService)

	// Use cases
	pushSubscriptionUseCase := usecase.NewPushSubscriptionUseCase(subscriptionRepo, pushService)
	pushNotificationUseCase := usecase.NewPushNotificationUseCase(jobRepo, subscriptionRepo, pushService)
	vapidUseCase := usecase.NewVAPIDUseCase(vapidService)

	// Handlers
	healthHandler := handler.NewHealthHandler()
	userHandler := handler.NewUserHandler(userUseCase)
	pushSubscriptionHandler := handler.NewPushSubscriptionHandler(pushSubscriptionUseCase)
	pushNotificationHandler := handler.NewPushNotificationHandler(pushNotificationUseCase)
	vapidHandler := handler.NewVAPIDHandler(vapidUseCase)
	mlHandler := handler.NewMLHandler()

	// Background service for processing push jobs
	go func() {
		for {
			if err := pushSenderService.ProcessPendingJobs(context.Background(), 100); err != nil {
				log.Printf("Error processing pending jobs: %v", err)
			}
			if err := pushSenderService.ProcessRetries(context.Background(), 5, 50); err != nil {
				log.Printf("Error processing retries: %v", err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /api/healthz", healthHandler.HealthCheck)

	// User API
	mux.HandleFunc("GET /api/users", userHandler.GetUsers)
	mux.HandleFunc("POST /api/users", userHandler.CreateUser)
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetUser)
	mux.HandleFunc("DELETE /api/users/{id}", userHandler.DeleteUser)

	// Web Push API
	mux.HandleFunc("GET /api/push/vapid-public-key", vapidHandler.GetPublicKey)
	mux.HandleFunc("POST /api/push/subscribe", pushSubscriptionHandler.Subscribe)
	mux.HandleFunc("DELETE /api/push/subscriptions/{id}", pushSubscriptionHandler.Unsubscribe)
	mux.HandleFunc("POST /api/push/send", pushNotificationHandler.SendNotification)
	mux.HandleFunc("POST /api/push/send/batch", pushNotificationHandler.SendBatchNotification)

	// ML (gRPC 経由) API プロキシ
	mux.HandleFunc("GET /api/ml/hello", mlHandler.HelloProxy)
	mux.HandleFunc("POST /api/ml/recognize", mlHandler.RecognizeImageProxy)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
