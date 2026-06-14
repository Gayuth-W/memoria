package observability

import "sync/atomic"

type Metrics struct {
	cacheHits   atomic.Int64
	cacheMisses atomic.Int64
	searches    atomic.Int64
	embeddings  atomic.Int64
	embedErrors atomic.Int64
}

func NewMetrics() *Metrics { return &Metrics{} }

func (m *Metrics) CacheHit()   { m.cacheHits.Add(1) }
func (m *Metrics) CacheMiss()  { m.cacheMisses.Add(1) }
func (m *Metrics) Search()     { m.searches.Add(1) }
func (m *Metrics) Embedding()  { m.embeddings.Add(1) }
func (m *Metrics) EmbedError() { m.embedErrors.Add(1) }

type Snapshot struct {
	CacheHits    int64   `json:"cache_hits"`
	CacheMisses  int64   `json:"cache_misses"`
	CacheHitRate float64 `json:"cache_hit_rate"`
	Searches     int64   `json:"searches"`
	Embeddings   int64   `json:"embeddings"`
	EmbedErrors  int64   `json:"embed_errors"`
}
