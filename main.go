package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"memoria/internal/cache"
	"memoria/internal/config"
	"memoria/internal/db"
	"memoria/internal/embedding"
	"memoria/internal/handler"
	"memoria/internal/middleware"
	"memoria/internal/observability"
	vector "memoria/internal/qdrant"
	"memoria/internal/repository"
	"memoria/internal/search"
	"memoria/internal/service"
	"memoria/internal/worker"
)

func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found")
	}

	logger := observability.NewLogger()
	metrics := observability.NewMetrics()
	database := db.NewDB()

	// embeddings
	embedder := embedding.NewOllamaEmbedder()

	// qdrant
	vectorStore := vector.NewVectorStore()

	if err := vectorStore.Init(); err != nil {
		log.Fatal(err)
	}

	// worker
	workerHandler := &worker.Handler{
		Embedder: embedder,
		Vector:   vectorStore,
	}

	backgroundWorker := worker.NewWorker(worker.Config{
		Buffer:     config.GetInt("WORKER_BUFFER", 100),
		Workers:    config.GetInt("WORKER_COUNT", 4),
		MaxRetries: config.GetInt("WORKER_MAX_RETRIES", 3),
		BaseDelay:  200 * time.Millisecond,
		Logger:     logger,
	}, workerHandler.Handle)

	// repositories
	userRepo := &repository.UserRepo{DB: database}
	sessionRepo := &repository.SessionRepo{DB: database}
	memoryRepo := &repository.MemoryRepo{DB: database}
	redisCache := cache.NewRedisCache()

	// services
	sessionService := &service.SessionService{Repo: sessionRepo, MemoryRepo: memoryRepo}
	memoryService := &service.MemoryService{Repo: memoryRepo, Worker: backgroundWorker, Cache: redisCache}
	userService := &service.UserService{Repo: userRepo}
	searchService := &search.Service{
		Embedder: embedder, Vector: vectorStore, Repo: memoryRepo,
		Cache: redisCache, Metrics: metrics, Logger: logger,
	}

	// handlers
	sessionHandler := &handler.SessionHandler{Service: sessionService}
	memoryHandler := &handler.MemoryHandler{Service: memoryService}
	userHandler := &handler.UserHandler{Service: userService}
	searchHandler := &handler.SearchHandler{Service: searchService}

	r := chi.NewRouter()

	// Add Request Logger middleware
	r.Use(middleware.RequestLogger(logger))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Expose metrics
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics.Snapshot())
	})

	// 1. Public Routes (No authentication needed)
	r.Post("/users", userHandler.Create)

	// 2. Protected Routes (Authenticated using API Key)
	r.Group(func(r chi.Router) {
		r.Use(middleware.APIKeyAuth(userRepo))
		// Rate limiting middleware: 100 requests per 1 minute window
		r.Use(middleware.RateLimit(redisCache, 100, 1*time.Minute))

		r.Get("/sessions", sessionHandler.ListByUser)
		r.Post("/sessions", sessionHandler.Create)
		r.Get("/sessions/{id}", sessionHandler.GetByID)
		r.Get("/sessions/{id}/memories", sessionHandler.GetMemories)

		r.Get("/memories", memoryHandler.ListByUser)
		r.Post("/memories", memoryHandler.Create)
		r.Delete("/memories/{id}", memoryHandler.Delete)

		r.Post("/search", searchHandler.Search)
	})

	srv := &http.Server{
		Addr:    ":" + portString,
		Handler: r,
	}

	go func() {
		log.Println("Server running on : ", portString)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Shutting down workers...")
	if err := backgroundWorker.Shutdown(ctx); err != nil {
		log.Fatal("Worker forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
