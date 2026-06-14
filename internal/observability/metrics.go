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
