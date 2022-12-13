package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
)

const inputFilePath = "input.txt"

func main() {
	if err := star1(); err != nil {
		panic(err)
	}
	if err := star2(); err != nil {
		panic(err)
	}
}

func star1() error {
	pairs, err := parseFile()
	if err != nil {
		return err
	}

	indicesSum := 0

	for i, pair := range pairs {
		switch isInOrder(pair[0], pair[1]) {
		case inOrderTrue:
			indicesSum += i + 1
			//fmt.Printf("pair %d: true\n", i+1)
		case inOrderFalse:
			//fmt.Printf("pair %d: false\n", i+1)
		case inOrderIndeterminate:
			//fmt.Printf("pair %d: indeterminate\n", i+1)
		default:
			panic("unreachable")
		}
	}

	fmt.Println(indicesSum)

	return nil
}

func star2() error {
	pairs, err := parseFile()
	if err != nil {
		return err
	}

	packets := make([][]any, 0, len(pairs)*2)
	for _, pair := range pairs {
		packets = append(packets, pair[0], pair[1])
	}

	var (
		divider1 = []any{[]any{2.0}}
		divider2 = []any{[]any{6.0}}
	)

	packets = append(packets, divider1, divider2)

	sort.Slice(packets, func(i, j int) bool {
		return isInOrder(packets[i], packets[j]) == inOrderTrue
	})

	divider1Pos, divider2Pos := -1, -1

	for i := range packets {
		if reflect.DeepEqual(packets[i], divider1) {
			divider1Pos = i + 1

		}
		if reflect.DeepEqual(packets[i], divider2) {
			divider2Pos = i + 1
		}
	}

	if divider1Pos == -1 || divider2Pos == -1 {
		return fmt.Errorf("dividers not found")
	}

	fmt.Println(divider1Pos * divider2Pos)

	return nil
}

type inOrder int

const (
	inOrderNil inOrder = iota
	inOrderTrue
	inOrderFalse
	inOrderIndeterminate
)

func isInOrder(l, r any) inOrder {
	var (
		lFloat, lFloatOk = l.(float64)
		rFloat, rFloatOk = r.(float64)
		lList, lListOk   = l.([]any)
		rList, rListOk   = r.([]any)
	)

	switch {
	case lFloatOk && rFloatOk:
		switch {
		case lFloat == rFloat:
			return inOrderIndeterminate
		case lFloat < rFloat:
			return inOrderTrue
		case lFloat > rFloat:
			return inOrderFalse
		default:
			panic("unreachable")
		}

	case lListOk && rListOk:
		commonLen := min(len(lList), len(rList))
		for i := 0; i < commonLen; i++ {
			switch isInOrder(lList[i], rList[i]) {
			case inOrderTrue:
				return inOrderTrue
			case inOrderFalse:
				return inOrderFalse
			case inOrderIndeterminate:
				continue
			default:
				panic("unreachable")
			}
		}
		switch {
		case len(lList) == len(rList):
			return inOrderIndeterminate
		case len(lList) < len(rList):
			return inOrderTrue
		case len(lList) > len(rList):
			return inOrderFalse
		default:
			panic("unreachable")
		}
	case lListOk && rFloatOk:
		return isInOrder(lList, []any{rFloat})
	case lFloatOk && rListOk:
		return isInOrder([]any{lFloat}, rList)
	default:
		fmt.Printf("unhandled types: %T, %T\n", l, r)
		panic(fmt.Sprintf("lListOk: %v, rListOk: %v, lFloatOk: %v, rFloatOk: %v", lListOk, rListOk, lFloatOk, rFloatOk))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parseFile() ([][2][]any, error) {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pairs := [][2][]any{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line1 := scanner.Bytes()
		var line1Data []any
		if err := json.Unmarshal(line1, &line1Data); err != nil {
			return nil, fmt.Errorf("error parsing line '%s', %w", line1, err)
		}

		if !scanner.Scan() {
			return nil, io.ErrUnexpectedEOF
		}

		line2 := scanner.Bytes()
		var line2Data []any
		if err := json.Unmarshal(line2, &line2Data); err != nil {
			return nil, fmt.Errorf("error parsing line '%s', %w", line2, err)
		}

		scanner.Scan()

		pairs = append(pairs, [2][]any{line1Data, line2Data})
	}
	return pairs, nil
}
