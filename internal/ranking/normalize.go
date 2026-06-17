package ranking

func Normalize(
	value float64,
	min float64,
	max float64,
) float64 {

	if max == min {
		return 1.0
	}

	return (value - min) / (max - min)
}
