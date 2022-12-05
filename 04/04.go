package main

import (
	"bufio"
	"fmt"
	"os"
)

const inputFilePath = "input.txt"

type Range struct {
	s, e int
}

func (r1 Range) Contains(r2 Range) bool {
	return r1.s <= r2.s && r1.e >= r2.e
}

func (r1 Range) Overlap(r2 Range) bool {
	return r1.s <= r2.s && r1.e >= r2.s || r2.s <= r1.s && r2.e >= r1.s
}

func main() {
	star1()
	star2()
}

func star1() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0

	for scanner.Scan() {
		line := scanner.Text()

		var r1, r2 Range

		if _, err := fmt.Sscanf(line, "%d-%d,%d-%d", &r1.s, &r1.e, &r2.s, &r2.e); err != nil {
			return err
		}

		if r1.Contains(r2) || r2.Contains(r1) {
			count++
		}
	}

	fmt.Println(count)

	return nil
}

func star2() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0

	for scanner.Scan() {
		line := scanner.Text()

		var r1, r2 Range

		if _, err := fmt.Sscanf(line, "%d-%d,%d-%d", &r1.s, &r1.e, &r2.s, &r2.e); err != nil {
			return err
		}

		if r1.Overlap(r2) {
			count++
		}
	}

	fmt.Println(count)

	return nil
}
