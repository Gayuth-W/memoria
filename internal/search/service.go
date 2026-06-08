package search

import (
	"memoria/internal/embedding"
	"memoria/internal/repository"
	"memoria/internal/vector"
)

type Service struct {
	Embedder embedding.Embedder
	Vector   *vector.VectorStore
	Repo     *repository.MemoryRepo
}

func (s *Service) Search(userID, query string) ([]string, error) {

	keywordIDs, _ := s.Repo.KeywordSearch(userID, query)

	vec, _ := s.Embedder.Embed(query)
	vectorIDs, _ := s.Vector.Search(vec, 10)

	seen := map[string]bool{}
	var result []string

	for _, id := range append(keywordIDs, vectorIDs...) {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}

	return result, nil
}
