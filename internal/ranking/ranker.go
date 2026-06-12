package ranking

import "sort"

func Rank(results []SearchResult) []SearchResult {
	for i := range results {

		results[i].FinalScore =
			ComputeScore(
				results[i].Similarity,
				results[i].Recency,
				results[i].Importance,
				results[i].SessionBoost,
			)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].FinalScore >
			results[j].FinalScore
	})

	return results
}
