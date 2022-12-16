package main

import "testing"

func TestComb(t *testing.T) {
	comb := Comb([]int{1, 2, 3, 4, 5}, 3)
	t.Log(comb)
}
