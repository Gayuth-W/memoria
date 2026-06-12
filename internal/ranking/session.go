package ranking

func SessionBoost(
	currentSession string,
	memorySession string,
) float64 {

	if currentSession == memorySession {
		return 1.0
	}

	return 0.0
}
