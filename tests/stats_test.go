package tests

import (
	"testing"

	"weather/internal"
)

func TestComputeStats_Empty(t *testing.T) {
	s := internal.ComputeStats(nil)
	if s.Average != 0 || s.Median != 0 || s.Min != 0 || s.Max != 0 {
		t.Errorf("expected zero stats for empty input, got %+v", s)
	}
}

func TestComputeStats_Single(t *testing.T) {
	s := internal.ComputeStats([]float64{5.0})
	if s.Min != 5 || s.Max != 5 || s.Average != 5 || s.Median != 5 {
		t.Errorf("unexpected stats for single value: %+v", s)
	}
}

func TestComputeStats_OddCount(t *testing.T) {
	s := internal.ComputeStats([]float64{1, 3, 5})
	if s.Min != 1 {
		t.Errorf("expected Min=1, got %v", s.Min)
	}
	if s.Max != 5 {
		t.Errorf("expected Max=5, got %v", s.Max)
	}
	if s.Average != 3 {
		t.Errorf("expected Average=3, got %v", s.Average)
	}
	if s.Median != 3 {
		t.Errorf("expected Median=3, got %v", s.Median)
	}
}

func TestComputeStats_EvenCount(t *testing.T) {
	s := internal.ComputeStats([]float64{10, 20, 30, 40})
	if s.Median != 25 {
		t.Errorf("expected Median=25, got %v", s.Median)
	}
	if s.Average != 25 {
		t.Errorf("expected Average=25, got %v", s.Average)
	}
	if s.Min != 10 || s.Max != 40 {
		t.Errorf("unexpected min/max: %+v", s)
	}
}

func TestComputeStats_UnsortedInput(t *testing.T) {
	s := internal.ComputeStats([]float64{5, 1, 3})
	if s.Min != 1 || s.Max != 5 {
		t.Errorf("expected Min=1 Max=5 for unsorted input, got %+v", s)
	}
}
