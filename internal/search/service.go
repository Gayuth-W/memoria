package search

import (
	"memoria/internal/embedding"
	vector "memoria/internal/qdrant"
	"memoria/internal/ranking"
	"memoria/internal/repository"
)

type Service struct {
	Embedder embedding.Embedder
	Vector   *vector.VectorStore
	Repo     *repository.MemoryRepo
}

func (s *Service) Search(userID string, currentSession string, query string) ([]ranking.SearchResult, error) {

	keywordIDs, _ := s.Repo.KeywordSearch(userID, query)

	vec, err := s.Embedder.Embed(query)
	if err != nil {
		return nil, err
	}

	vectorResults, err := s.Vector.Search(
		vec,
		10,
	)

	if err != nil {
		vectorResults = []vector.VectorResult{}
	}

	similarityMap := make(
		map[string]float64,
	)

	var vectorIDs []string

	for _, id := range vectorResults {
		similarityMap[id.MemoryID] = id.Score
		vectorIDs = append(vectorIDs, id.MemoryID)
	}

	seen := map[string]bool{}

	var ids []string

	for _, id := range append(
		keywordIDs,
		vectorIDs...,
	) {

		if !seen[id] {

			seen[id] = true

			ids = append(
				ids,
				id,
			)
		}
	}
	memories, err := s.Repo.GetByIDs(ids)

	if err != nil {
		return nil, err
	}

	return result, nil
}
