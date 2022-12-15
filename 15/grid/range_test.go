package grid

import (
	"reflect"
	"testing"
)

func TestRange_Overlap(t *testing.T) {
	tests := []struct {
		r1, r2 Range
		want   bool
	}{
		{Range{0, 0}, Range{0, 0}, true},
		{Range{0, 0}, Range{1, 1}, false},
		{Range{1, 1}, Range{0, 0}, false},
		{Range{0, 1}, Range{1, 2}, true},
		{Range{1, 2}, Range{0, 1}, true},
		{Range{0, 2}, Range{1, 2}, true},
		{Range{0, 4}, Range{1, 2}, true},
		{Range{1, 2}, Range{0, 4}, true},
	}
	for _, tt := range tests {
		if got := tt.r1.Overlap(tt.r2); got != tt.want {
			t.Errorf("Overlap(%v, %v) = %v, want %v", tt.r1, tt.r2, got, tt.want)
		}
	}
}

func TestRange(t *testing.T) {
	rs := RangeSet{}

	rs.Add(Range{0, 0})
	rs.Add(Range{1, 1})
	rs.Add(Range{0, 1})

	want := RangeSet{Range{0, 1}}
	if !reflect.DeepEqual(rs, want) {
		t.Errorf("got %v, want %v", rs, want)
	}

	rs.Add(Range{2, 5})
	want = RangeSet{Range{0, 1}, Range{2, 5}}
	if !reflect.DeepEqual(rs, RangeSet{Range{0, 1}, Range{2, 5}}) {
		t.Errorf("got %v, want %v", rs, want)
	}
}