package tests

import (
	"fmt"
	"testing"

	"memoria/internal/db"
	"memoria/internal/embedding"
	vector "memoria/internal/qdrant"
	"memoria/internal/repository"
	"memoria/internal/search"

	"github.com/joho/godotenv"
)

func TestSearch(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Log("Warning: Error loading .env file")
	}
	database := db.NewDB()

	memoryRepo := &repository.MemoryRepo{
		DB: database,
	}

	embedder := embedding.NewOllamaEmbedder()

	vectorStore := vector.NewVectorStore()

	searchService := &search.Service{
		Repo:     memoryRepo,
		Embedder: embedder,
		Vector:   vectorStore,
	}

	results, _, err := searchService.Search(
		"7a6dbabb-6560-4d0c-90b6-ae54eea5a9ac",
		"sd3903812-bf3d-4e1b-9b8b-426e2bc23ba7",
		"test",
	)

	if err != nil {
		t.Fatal(err)
	}

	for _, r := range results {
		fmt.Printf(
			"%s %.4f\n",
			r.MemoryID,
			r.FinalScore,
		)
	}
}
