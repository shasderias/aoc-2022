package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const inputFilePath = "input.txt"

func main() {
	if err := star1(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if err := star2(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
func star1() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	prioritySum := 0

	lineScanner := bufio.NewScanner(f)
	for lineScanner.Scan() {
		line := lineScanner.Text()
		bagSize := len(line)
		comp1, comp2 := line[:bagSize/2], line[bagSize/2:]
		common := intersect(comp1, comp2)
		prioritySum += itemPriority(common)
	}

	fmt.Println(prioritySum)

	return nil
}

func star2() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	prioritySum := 0

	lineScanner := bufio.NewScanner(f)
	for lineScanner.Scan() {
		bag1 := lineScanner.Text()
		if !lineScanner.Scan() {
			panic("line count not multiple of 3")
		}
		bag2 := lineScanner.Text()
		if !lineScanner.Scan() {
			panic("line count not multiple of 3")
		}
		bag3 := lineScanner.Text()

		prioritySum += itemPriority(intersect(bag3, intersect(bag1, bag2)))
	}

	fmt.Println(prioritySum)

	return nil
}

func intersect(a, b string) string {
	intersection := strings.Builder{}
	for _, c := range a {
		if strings.ContainsRune(b, c) && !strings.ContainsRune(intersection.String(), c) {
			intersection.WriteRune(c)
		}
	}
	return intersection.String()
}

func itemPriority(item string) int {
	if len(item) != 1 {
		panic(fmt.Sprintf("itemPriority: item must be a single character, got %s", item))
	}
	cp := int(item[0])
	switch {
	case cp >= 'a' && cp <= 'z':
		return cp - 'a' + 1
	case cp >= 'A' && cp <= 'Z':
		return cp - 'A' + 1 + 26
	default:
		return 0
	}
}
