package ranking

import "sort"

func Rank(results []SearchResult) []SearchResult {
	if len(results) == 0 {
		return results
	}

	var minSim, maxSim float64 = results[0].Similarity, results[0].Similarity
	var minRec, maxRec float64 = results[0].Recency, results[0].Recency
	var minImp, maxImp float64 = results[0].Importance, results[0].Importance

	for _, r := range results {
		if r.Similarity < minSim {
			minSim = r.Similarity
		}
		if r.Similarity > maxSim {
			maxSim = r.Similarity
		}
		if r.Recency < minRec {
			minRec = r.Recency
		}
		if r.Recency > maxRec {
			maxRec = r.Recency
		}
		if r.Importance < minImp {
			minImp = r.Importance
		}
		if r.Importance > maxImp {
			maxImp = r.Importance
		}
	}

	for i := range results {
		normSim := Normalize(results[i].Similarity, minSim, maxSim)
		normRec := Normalize(results[i].Recency, minRec, maxRec)
		normImp := Normalize(results[i].Importance, minImp, maxImp)

		results[i].FinalScore = ComputeScore(
			normSim,
			normRec,
			normImp,
			results[i].SessionBoost,
		)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].FinalScore > results[j].FinalScore
	})

	return results
}
