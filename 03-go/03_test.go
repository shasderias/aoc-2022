package main

import "testing"

func TestIntersection(t *testing.T) {
	testCases := []struct{ a, b, expected string }{
		{"abc", "abc", "abc"},
		{"abc", "def", ""},
		{"abc", "abd", "ab"},
		{"abc", "ade", "a"},
	}
	for _, tt := range testCases {
		actual := intersect(tt.a, tt.b)
		if actual != tt.expected {
			t.Errorf("intersect(%q, %q) = %q, want %q", tt.a, tt.b, actual, tt.expected)
		}
	}
}
