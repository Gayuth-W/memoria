package search

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"memoria/internal/cache"
	"memoria/internal/embedding"
	"memoria/internal/observability"
	vector "memoria/internal/qdrant"
	"memoria/internal/ranking"
	"memoria/internal/repository"
	"time"
)

type Service struct {
	Embedder embedding.Embedder
	Vector   *vector.VectorStore
	Repo     *repository.MemoryRepo
	Cache    *cache.RedisCache
	Metrics  *observability.Metrics
	Logger   *slog.Logger
}

type Trace struct {
	CacheHit       bool  `json:"cache_hit"`
	KeywordResults int   `json:"keyword_results"`
	VectorResults  int   `json:"vector_results"`
	MergedResults  int   `json:"merged_results"`
	FinalResults   int   `json:"final_results"`
	EmbedMs        int64 `json:"embed_ms"`
	VectorMs       int64 `json:"vector_ms"`
	KeywordMs      int64 `json:"keyword_ms"`
	TotalMs        int64 `json:"total_ms"`
}

func cacheKey(userID, session, query string) string {
	sum := sha256.Sum256([]byte(userID + "|" + session + "|" + query))
	return "search:" + hex.EncodeToString(sum[:8])
}

func (s *Service) Search(userID string, currentSession string, query string) ([]ranking.SearchResult, *Trace, error) {

	//caching layer
	start := time.Now()
	trace := &Trace{}
	key := cacheKey(userID, currentSession, query)

	if s.Cache != nil {
		if cached, err := s.Cache.Get(key); err == nil && cached != "" {
			var results []ranking.SearchResult
			if json.Unmarshal([]byte(cached), &results) == nil {
				trace.CacheHit = true
				trace.FinalResults = len(results)
				trace.TotalMs = time.Since(start).Milliseconds()
				if s.Metrics != nil {
					s.Metrics.CacheHit()
					s.Metrics.Search()
				}
				s.log(trace)
				return results, trace, nil
			}
		}
		if s.Metrics != nil {
			s.Metrics.CacheMiss()
		}
	}

	// keyword
	kwStart := time.Now()
	keywordIDs, _ := s.Repo.KeywordSearch(userID, query)
	trace.KeywordResults = len(keywordIDs)
	trace.KeywordMs = time.Since(kwStart).Milliseconds()

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

	var results []ranking.SearchResult

	for _, memory := range memories {

		results = append(
			results,
			ranking.SearchResult{

				MemoryID:  memory.ID,
				SessionID: memory.SessionID,
				Text:      memory.Text,

				Similarity: similarityMap[memory.ID],

				Recency: ranking.RecencyScore(
					memory.CreatedAt,
				),

				Importance: memory.ImportanceScore,

				SessionBoost: ranking.SessionBoost(
					currentSession,
					memory.SessionID,
				),

				CreatedAt: memory.CreatedAt,
			},
		)
	}
	ranked := ranking.Rank(results)
	return ranked, nil
}

func (s *Service) log(t *Trace) {
	if s.Logger == nil {
		return
	}
	s.Logger.Info("search",
		slog.Bool("cache_hit", t.CacheHit),
		slog.Int("keyword_results", t.KeywordResults),
		slog.Int("vector_results", t.VectorResults),
		slog.Int("final_results", t.FinalResults),
		slog.Int64("embed_ms", t.EmbedMs),
		slog.Int64("vector_ms", t.VectorMs),
		slog.Int64("keyword_ms", t.KeywordMs),
		slog.Int64("total_ms", t.TotalMs),
	)
}
