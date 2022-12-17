package grid

import "testing"

func TestNewBlock(t *testing.T) {
	testCases := []struct {
		block              *Block
		wantRows, wantCols int
		wantPoints         map[Point]bool
	}{
		{
			NewBlock([][]byte{{1, 1, 0, 1}}),
			1, 4,
			map[Point]bool{
				{0, 0}: true,
				{2, 0}: false,
				{3, 0}: true,
			},
		},
		{
			NewBlock([][]byte{
				{0, 0, 0, 1},
				{0, 0, 0, 1},
				{1, 1, 1, 1},
			}),
			3, 4,
			map[Point]bool{
				{0, 0}: true,
				{1, 0}: true,
				{2, 0}: true,
				{3, 0}: true,
				{3, 0}: true,
				{3, 1}: true,
				{3, 2}: true,
			},
		},
	}
	for i, tt := range testCases {
		t.Log(tt.block)
		if tt.block.maxY != tt.wantRows {
			t.Errorf("got maxY %d, want %d", tt.block.maxY, tt.wantRows)
		}
		if tt.block.maxX != tt.wantCols {
			t.Errorf("got maxX %d, want %d", tt.block.maxX, tt.wantCols)
		}
		for p, want := range tt.wantPoints {
			if got := tt.block.Get(p); got != want {
				t.Errorf("%d: Get(%v): got %t, want %t", i, p, got, want)
			}
		}
	}

}
