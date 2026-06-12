package tests

import (
	"testing"
	"time"

	"memoria/internal/ranking"
)

func TestRanking(t *testing.T) {

	results := []ranking.SearchResult{
		{
			MemoryID: "1",

			Similarity:   0.9,
			Importance:   1.0,
			SessionBoost: 1.0,

			CreatedAt: time.Now(),
		},
		{
			MemoryID: "2",

			Similarity:   0.4,
			Importance:   0.1,
			SessionBoost: 0,

			CreatedAt: time.Now().Add(
				-30 * 24 * time.Hour,
			),
		},
	}

	ranked := ranking.Rank(results)

	if ranked[0].MemoryID != "1" {
		t.Fatal(
			"ranking failed",
		)
	}
}
