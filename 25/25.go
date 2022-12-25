package main

import (
	"bufio"
	"fmt"
	"os"

	"aoc-2022-25/snafu"
)

func main() {
	star1("input.txt")
}

func star1(inputPath string) {
	nums := readInput(inputPath)
	acc := 0
	for _, num := range nums {
		acc += num.Int()
	}
	fmt.Println(snafu.FromInt(acc).Snafu())
}

func readInput(inputPath string) []snafu.Num {
	f, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	num := []snafu.Num{}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		num = append(num, snafu.FromString(line))
	}

	return num
}
