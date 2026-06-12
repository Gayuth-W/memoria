package ranking

import "time"

type SearchResult struct {
	MemoryID  string
	SessionID string
	Text      string

	Similarity   float64
	Recency      float64
	Importance   float64
	SessionBoost float64

	FinalScore float64

	CreatedAt time.Time
}
