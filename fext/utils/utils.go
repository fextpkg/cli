package utils

func FindMinValue(values []int) int {
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}

	return min
}
