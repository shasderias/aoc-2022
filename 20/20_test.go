package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestList(t *testing.T) {
	//Initial arrangement:
	//1, 2, -3, 3, -2, 0, 4
	list := NewListFromIntSlice(1, 2, -3, 3, -2, 0, 4)

	//1 moves between 2 and -3:
	//2, 1, -3, 3, -2, 0, 4
	list.Shift(0, 1)
	checkList(t, list, NewListFromIntSlice(2, 1, -3, 3, -2, 0, 4))

	//2 moves between -3 and 3:
	//1, -3, 2, 3, -2, 0, 4
	list.Shift(0, 2)
	checkList(t, list, NewListFromIntSlice(1, -3, 2, 3, -2, 0, 4))

	//-3 moves between -2 and 0:
	//1, 2, 3, -2, -3, 0, 4
	list.Shift(1, -3)
	checkList(t, list, NewListFromIntSlice(1, 2, 3, -2, -3, 0, 4))

	//3 moves between 0 and 4:
	//1, 2, -2, -3, 0, 3, 4
	list.Shift(2, 3)
	checkList(t, list, NewListFromIntSlice(1, 2, -2, -3, 0, 3, 4))

	//-2 moves between 4 and 1:
	//1, 2, -3, 0, 3, 4, -2
	list.Shift(2, -2)
	checkList(t, list, NewListFromIntSlice(1, 2, -3, 0, 3, 4, -2))

	//0 does not move:
	//1, 2, -3, 0, 3, 4, -2
	list.Shift(3, 0)
	checkList(t, list, NewListFromIntSlice(1, 2, -3, 0, 3, 4, -2))

	//4 moves between -3 and 0:
	//1, 2, -3, 4, 0, 3, -2
	list.Shift(5, 4)
	checkList(t, list, NewListFromIntSlice(1, 2, -3, 4, 0, 3, -2))
}

func TestList2(t *testing.T) {
	list := NewListFromIntSlice(1, 2, 3)

	list.Shift(0, 1)
	checkList(t, list, NewListFromIntSlice(2, 1, 3))

	list.Shift(1, 1)
	checkList(t, list, NewListFromIntSlice(2, 3, 1))

	list.Shift(2, 1)
	checkList(t, list, NewListFromIntSlice(2, 1, 3))

	list.Shift(1, 2)
	checkList(t, list, NewListFromIntSlice(2, 1, 3))

	list.Shift(1, 1000)
	checkList(t, list, NewListFromIntSlice(2, 1, 3))

	list.Shift(1, -1)
	checkList(t, list, NewListFromIntSlice(2, 3, 1))

	list.Shift(2, -1)
	checkList(t, list, NewListFromIntSlice(2, 1, 3))

	list.Shift(1, -4)
	checkList(t, list, NewListFromIntSlice(2, 1, 3))

	list.Shift(1, -4000)
	checkList(t, list, NewListFromIntSlice(2, 1, 3))

	list = NewListFromIntSlice(1, 2, 3, 4, 5)
	list.Shift(0, -1)
	checkList(t, list, NewListFromIntSlice(2, 3, 4, 1, 5))

	list = NewListFromIntSlice(1, 2, 3, 4, 5)
	list.Shift(0, -2)
	checkList(t, list, NewListFromIntSlice(2, 3, 1, 4, 5))

	list = NewListFromIntSlice(1, 2, 3, 4, 5)
	list.Shift(0, -3)
	checkList(t, list, NewListFromIntSlice(2, 1, 3, 4, 5))

	list = NewListFromIntSlice(1, 2, 3, 4, 5)
	list.Shift(0, -4000)
	checkList(t, list, NewListFromIntSlice(1, 2, 3, 4, 5))

	list = NewListFromIntSlice(1, 2, 3, 4, 5)
	list.Shift(0, -4001)
	checkList(t, list, NewListFromIntSlice(2, 3, 4, 1, 5))
}

func checkList(t *testing.T, list List[Elem], expected List[Elem]) {
	t.Helper()
	if diff := cmp.Diff(list, expected, cmpopts.IgnoreUnexported(Elem{})); diff != "" {
		t.Errorf("list mismatch (-want +got):\n%s", diff)
	}
}
