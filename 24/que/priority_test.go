package que

import (
	"fmt"
	"testing"
)

type elem int

func (e elem) Priority() int { return int(e) }

func TestSanity(t *testing.T) {
	minQueue := NewPriority[elem]()
	elems := []elem{7, 1, 4, 3, 5, 5, 1}

	for _, e := range elems {
		minQueue.Push(e)
	}

	sortedElems := []elem{1, 1, 3, 4, 5, 5, 7}
	for _, e := range sortedElems {
		poppedElem := minQueue.Pop()
		fmt.Println(poppedElem)
		if poppedElem != e {
			t.Errorf("expected %d, got %d", e, poppedElem)
		}
	}
}
