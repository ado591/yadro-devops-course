package internal

import (
	"sort"

	"weather/internal/models"
)

func ComputeStats(temps []float64) models.TemperatureStats {
	n := len(temps)
	if n == 0 {
		return models.TemperatureStats{}
	}

	sorted := make([]float64, n)
	copy(sorted, temps)
	sort.Float64s(sorted)

	min := sorted[0]
	max := sorted[n-1]

	sum := 0.0
	for _, t := range temps {
		sum += t
	}

	var median float64
	if n%2 == 0 {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	} else {
		median = sorted[n/2]
	}

	return models.TemperatureStats{
		Average: sum / float64(n),
		Median:  median,
		Min:     min,
		Max:     max,
	}
}
