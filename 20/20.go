package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const decryptionKey = 811589153

func main() {
	if err := star1("input.txt"); err != nil {
		panic(err)
	}
	if err := star2("input.txt"); err != nil {
		panic(err)
	}
}

type ListElem interface {
	comparable
	Value() int
}

type List[T ListElem] []T

func NewListFromIntSlice(s ...int) List[Elem] {
	l := List[Elem]{}
	for i, v := range s {
		l = append(l, Elem{val: v, id: i})
	}
	return l
}

func (l *List[T]) IndexOf(value T) int {
	for i, v := range *l {
		if v == value {
			return i
		}
	}
	return -1
}

func (l *List[T]) FindIndexOfValue(v int) int {
	for i, elem := range *l {
		if elem.Value() == v {
			return i
		}
	}
	return -1
}

func (l *List[T]) Index(idx int) T {
	if idx >= len(*l) {
		idx = idx % len(*l)
	}
	return (*l)[idx]
}

func (l *List[T]) Shift(idx, offset int) {
	if offset == 0 {
		return
	}

	lenMinusOne := len(*l) - 1

	pos := ((((offset - 1 + lenMinusOne) % (lenMinusOne)) + idx) % lenMinusOne) + 1
	if pos < 0 {
		pos += lenMinusOne
	}
	effectiveOffset := pos - idx

	offset = effectiveOffset

	val := (*l)[idx]
	if offset > 0 {
		copy((*l)[idx:idx+offset], (*l)[idx+1:idx+offset+1])
		(*l)[pos] = val
	} else {
		copy((*l)[idx+offset+1:idx+1], (*l)[idx+offset:idx])
		(*l)[pos] = val
	}
}

type Elem struct {
	id  int
	val int
}

func (e Elem) Value() int {
	return e.val
}

func readInput(inputFilePath string) (List[Elem], error) {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	list := List[Elem]{}

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		if line == "" {
			break
		}

		num, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}

		list = append(list, Elem{i, num})
	}

	return list, nil
}

func star1(inputFilePath string) error {
	list, err := readInput(inputFilePath)
	if err != nil {
		return err
	}

	mix(list, 1)

	idxOf0 := list.FindIndexOfValue(0)

	coord1 := list.Index(idxOf0 + 1000)
	coord2 := list.Index(idxOf0 + 2000)
	coord3 := list.Index(idxOf0 + 3000)

	fmt.Println(coord1.val + coord2.val + coord3.val)

	return nil
}

func star2(inputFilePath string) error {
	list, err := readInput(inputFilePath)
	if err != nil {
		return err
	}

	applyDecrypt(list, decryptionKey)

	mix(list, 10)

	idxOf0 := list.FindIndexOfValue(0)

	coord1 := list.Index(idxOf0 + 1000)
	coord2 := list.Index(idxOf0 + 2000)
	coord3 := list.Index(idxOf0 + 3000)

	fmt.Println(coord1.val + coord2.val + coord3.val)

	return nil
}

func applyDecrypt(list List[Elem], key int) {
	for i := range list {
		list[i].val *= key
	}
}

func mix[T ListElem](l List[T], rounds int) List[T] {
	order := make(List[T], len(l))
	copy(order, l)

	for i := 0; i < rounds; i++ {
		for _, num := range order {
			l.Shift(l.IndexOf(num), num.Value())
		}
	}

	return l
}
