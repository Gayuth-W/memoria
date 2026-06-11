package ranking

import (
	"math"
	"time"
)

func RecencyScore(createdAt time.Time) float64 {
	days := time.Since(createdAt).Hours() / 24
	return math.Exp(-days / 30)
}
