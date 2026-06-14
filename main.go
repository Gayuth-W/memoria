package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	"memoria/internal/ranking"
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
	sessionService := &service.SessionService{Repo: sessionRepo}
	memoryService := &service.MemoryService{Repo: memoryRepo, Worker: backgroundWorker}
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
		r.Post("/search", searchHandler.Search)
	})

	cacheKey := "testing"

	cached, err := redisCache.Get(cacheKey)

	if err == nil {

		var results []ranking.SearchResult

		json.Unmarshal(
			[]byte(cached),
			&results,
		)

		bytes, _ := json.Marshal(results)

		redisCache.Set(cacheKey, string(bytes))
		fmt.Println(results)

	}

	// searchService := &search.Service{
	// 	Embedder: embedder,
	// 	Vector:   vectorStore,
	// 	Repo:     memoryRepo,
	// 	Cache:    redisCache,
	// }

	log.Println("Server running on : ", portString)
	http.ListenAndServe(":"+portString, r)
}
