package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"memoria/internal/db"
	"memoria/internal/handler"
	"memoria/internal/middleware"
	"memoria/internal/repository"
	"memoria/internal/service"
)

func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found")
	}

	database := db.NewDB()

	userRepo := &repository.UserRepo{DB: database}
	sessionRepo := &repository.SessionRepo{DB: database}
	memoryRepo := &repository.MemoryRepo{DB: database}

	sessionService := &service.SessionService{Repo: sessionRepo}
	memoryService := &service.MemoryService{Repo: memoryRepo}
	userService := &service.UserService{Repo: userRepo}

	sessionHandler := &handler.SessionHandler{Service: sessionService}
	memoryHandler := &handler.MemoryHandler{Service: memoryService}
	userHandler := &handler.UserHandler{Service: userService}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// 1. Public Routes (No authentication needed)
	r.Post("/users", userHandler.Create)
	// 2. Protected Routes (Authenticated using API Key)
	r.Group(func(r chi.Router) {
		r.Use(middleware.APIKeyAuth(userRepo))
		r.Post("/sessions", sessionHandler.Create)
		r.Post("/memories", memoryHandler.Create)
	})

	log.Println("Server running on : ", portString)
	http.ListenAndServe(":"+portString, r)
}
