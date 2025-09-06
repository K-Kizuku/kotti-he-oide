package main

import (
	"log"
	"net/http"
	"os"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/infrastructure/persistence"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	userRepo := persistence.NewMemoryUserRepository()
	userService := service.NewUserService(userRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, userService)

	healthHandler := handler.NewHealthHandler()
	userHandler := handler.NewUserHandler(userUseCase)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthHandler.HealthCheck)

	mux.HandleFunc("GET /api/users", userHandler.GetUsers)
	mux.HandleFunc("POST /api/users", userHandler.CreateUser)
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetUser)
	mux.HandleFunc("DELETE /api/users/{id}", userHandler.DeleteUser)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
