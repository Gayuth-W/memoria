package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"memoria/internal/cache"
	"memoria/internal/db"
	"memoria/internal/embedding"
	"memoria/internal/handler"
	"memoria/internal/middleware"
	vector "memoria/internal/qdrant"
	"memoria/internal/ranking"
	"memoria/internal/repository"
	"memoria/internal/service"
	"memoria/internal/worker"
)

func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found")
	}

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

	backgroundWorker := worker.NewWorker(
		100,
		workerHandler.Handle,
	)

	// repositories
	userRepo := &repository.UserRepo{DB: database}
	sessionRepo := &repository.SessionRepo{DB: database}
	memoryRepo := &repository.MemoryRepo{DB: database}

	// services
	sessionService := &service.SessionService{Repo: sessionRepo}
	memoryService := &service.MemoryService{Repo: memoryRepo, Worker: backgroundWorker}
	userService := &service.UserService{Repo: userRepo}

	// handlers
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

	redisCache := cache.NewRedisCache()

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
