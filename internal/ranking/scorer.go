package ranking

func ComputeScore(similarity float64, recency float64, importance float64, session float64) float64 {

	return similarity*0.60 + recency*0.20 + importance*0.10 + session*0.10
}
