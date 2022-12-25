package snafu_test

import (
	"testing"

	"aoc-2022-25/snafu"
)

func TestFromString(t *testing.T) {
	testCases := []struct {
		sna string
		dec snafu.Num
	}{
		{"1", 1},
		{"2", 2},
		{"1=", 3},
		{"1-", 4},
		{"10", 5},
		{"11", 6},
		{"12", 7},
		{"2=", 8},
		{"2-", 9},
		{"20", 10},
		{"1=0", 15},
		{"1-0", 20},
		{"1=11-2", 2022},
		{"1-0---0", 12345},
		{"1121-1110-1=0", 314159265},
	}
	for _, tt := range testCases {
		snafuNum := snafu.FromString(tt.sna)
		if snafuNum != tt.dec {
			t.Fatalf("FromString(%q) = %d, dec %d", tt.dec, snafuNum, tt.dec)
		}
		snafuStr := snafu.FromInt(int(tt.dec)).Snafu()
		if snafuStr != tt.sna {
			t.Fatalf("FromInt(%d).Snafu() = %q, sna %q", tt.dec, snafuStr, tt.sna)
		}
	}
}

func TestScratch(t *testing.T) {
	for i := 0; i < 50; i++ {
		snafu.FromInt(i).Snafu()
	}
}
