package main

import "testing"

func TestRange(t *testing.T) {
	testCases := []struct {
		r1, r2 Range
		want   bool
	}{
		{Range{1, 10}, Range{1, 10}, true},
		{Range{1, 10}, Range{2, 9}, true},
		{Range{1, 10}, Range{1, 9}, true},
		{Range{1, 10}, Range{2, 10}, true},
		{Range{1, 10}, Range{0, 10}, false},
		{Range{1, 10}, Range{1, 11}, false},
		{Range{1, 10}, Range{0, 11}, false},
		{Range{2, 8}, Range{3, 7}, true},
		{Range{4, 6}, Range{6, 6}, true},
	}
	for _, tt := range testCases {
		got := tt.r1.Contains(tt.r2)
		if got != tt.want {
			t.Errorf("Range(%v).Contains(%v) = %v, want %v", tt.r1, tt.r2, got, tt.want)
		}
	}
}
