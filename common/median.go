package common

import "sort"

func Median(data []int) float64 {
	dataCopy := make([]int, len(data))
	copy(dataCopy, data)

	sort.Ints(dataCopy)

	var median float64
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = float64(dataCopy[l/2-1]+dataCopy[l/2]) / 2
	} else {
		median = float64(dataCopy[l/2])
	}

	return median
}
