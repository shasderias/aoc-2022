package main

import "testing"

func TestStartOfPacket(t *testing.T) {
	testCases := []struct {
		buf  [4]byte
		want bool
	}{
		{[4]byte{0, 0, 0, 0}, false},
		{[4]byte{0, 0, 0, 1}, false},
		{[4]byte{0, 0, 1, 1}, false},
		{[4]byte{0, 1, 2, 3}, true},
		{[4]byte{3, 1, 2, 3}, false},
		{[4]byte{1, 2, 3, 3}, false},
		{[4]byte{3, 3, 1, 2}, false},
		{[4]byte{3, 3, 1, 2}, false},
		{[4]byte{3, 1, 2, 3}, false},
	}

	for _, tt := range testCases {
		if got := (&sopDetector{tt.buf}).isStartOfPacket(); got != tt.want {
			t.Errorf("isStartOfPacket(%v) = %v, want %v", tt.buf, got, tt.want)
		}
	}
}

func TestSopDetector_Write(t *testing.T) {
	d := &sopDetector{}

	d.Write([]byte{1, 2, 3, 4, 5})
	if got := d.buf; got != [4]byte{2, 3, 4, 5} {
		t.Fatal(got)
	}

	d.Write([]byte{6, 7})
	if got := d.buf; got != [4]byte{4, 5, 6, 7} {
		t.Fatal(got)
	}

	d.Write([]byte{8, 9, 10, 11, 12, 13})
	if got := d.buf; got != [4]byte{10, 11, 12, 13} {
		t.Fatal(got)
	}
}
