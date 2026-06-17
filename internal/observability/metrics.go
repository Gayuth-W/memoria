package observability

import "sync/atomic"

type Metrics struct {
	cacheHits     atomic.Int64
	cacheMisses   atomic.Int64
	searches      atomic.Int64
	embeddings    atomic.Int64
	embedErrors   atomic.Int64
	totalSearchMs atomic.Int64
	totalEmbedMs  atomic.Int64
}

func NewMetrics() *Metrics { return &Metrics{} }

func (m *Metrics) CacheHit()               { m.cacheHits.Add(1) }
func (m *Metrics) CacheMiss()              { m.cacheMisses.Add(1) }
func (m *Metrics) Search()                 { m.searches.Add(1) }
func (m *Metrics) Embedding()              { m.embeddings.Add(1) }
func (m *Metrics) EmbedError()             { m.embedErrors.Add(1) }
func (m *Metrics) RecordSearchMs(ms int64) { m.totalSearchMs.Add(ms) }
func (m *Metrics) RecordEmbedMs(ms int64)  { m.totalEmbedMs.Add(ms) }

type Snapshot struct {
	CacheHits    int64   `json:"cache_hits"`
	CacheMisses  int64   `json:"cache_misses"`
	CacheHitRate float64 `json:"cache_hit_rate"`
	Searches     int64   `json:"searches"`
	Embeddings   int64   `json:"embeddings"`
	EmbedErrors  int64   `json:"embed_errors"`
	AvgSearchMs  float64 `json:"avg_search_ms"`
	AvgEmbedMs   float64 `json:"avg_embed_ms"`
}

func (m *Metrics) Snapshot() Snapshot {
	//Accessing them using load because they are atomic.Int64 values, so it reads the current values safely when multiple go routines access it
	hits, misses := m.cacheHits.Load(), m.cacheMisses.Load()
	rate := 0.0
	if total := hits + misses; total > 0 {
		rate = float64(hits) / float64(total)
	}

	searches := m.searches.Load()
	embeddings := m.embeddings.Load()

	avgSearch := 0.0
	if searches > 0 {
		avgSearch = float64(m.totalSearchMs.Load()) / float64(searches)
	}

	avgEmbed := 0.0
	if embeddings > 0 {
		avgEmbed = float64(m.totalEmbedMs.Load()) / float64(embeddings)
	}

	//Creating a snapshot since Metrics uses atmoic.Int64 values are internal synchronization primitives, without exposing them I am creating a clean struct containing plain numbers
	return Snapshot{
		CacheHits: hits, CacheMisses: misses, CacheHitRate: rate,
		Searches: searches, Embeddings: embeddings,
		EmbedErrors: m.embedErrors.Load(),
		AvgSearchMs: avgSearch, AvgEmbedMs: avgEmbed,
	}
}
