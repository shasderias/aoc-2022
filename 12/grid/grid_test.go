package grid_test

import (
	"aoc-2022-09/grid"
	"testing"
)

func TestPoint_DirOf(t *testing.T) {
	testCases := []struct {
		a, b grid.Point
		want grid.Dir
	}{
		{grid.Point{0, 0}, grid.Point{0, 0}, grid.DirNil},
		{grid.Point{0, 0}, grid.Point{0, 1}, grid.DirN},
		{grid.Point{0, 0}, grid.Point{1, 1}, grid.DirNE},
		{grid.Point{0, 0}, grid.Point{1, 0}, grid.DirE},
		{grid.Point{0, 0}, grid.Point{1, -1}, grid.DirSE},
		{grid.Point{0, 0}, grid.Point{0, -1}, grid.DirS},
		{grid.Point{0, 0}, grid.Point{-1, -1}, grid.DirSW},
		{grid.Point{0, 0}, grid.Point{-1, 0}, grid.DirW},
		{grid.Point{0, 0}, grid.Point{-1, 1}, grid.DirNW},
		{grid.Point{0, 0}, grid.Point{0, 4}, grid.DirN},
		{grid.Point{4, 4}, grid.Point{0, 0}, grid.DirSW},
	}

	for _, tt := range testCases {
		got := tt.a.DirOf(tt.b)
		if got != tt.want {
			t.Errorf("DirOf(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestPoint_Distance(t *testing.T) {
	testCases := []struct {
		a, b grid.Point
		want int
	}{
		{grid.Point{0, 0}, grid.Point{0, 0}, 0},
		{grid.Point{0, 0}, grid.Point{0, 1}, 1},
		{grid.Point{0, 0}, grid.Point{1, 1}, 2},
	}

	for _, tt := range testCases {
		got := tt.a.Distance(tt.b)
		if got != tt.want {
			t.Errorf("Distance(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
