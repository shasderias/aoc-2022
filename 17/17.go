package main

import (
	"aoc-2022-17/grid"
	"fmt"
	"os"
)

type Dir byte

const (
	DirNil Dir = iota
	DirLeft
	DirRight
)

const Star1MaxBlocks = 2022
const Star2MaxBlocks = 1000000000000

var Shapes = []*grid.Block{
	grid.NewBlock([][]byte{
		{1, 1, 1, 1},
	}),
	grid.NewBlock([][]byte{
		{0, 1, 0},
		{1, 1, 1},
		{0, 1, 0},
	}),
	grid.NewBlock([][]byte{
		{0, 0, 1},
		{0, 0, 1},
		{1, 1, 1},
	}),
	grid.NewBlock([][]byte{
		{1},
		{1},
		{1},
		{1},
	}),
	grid.NewBlock([][]byte{
		{1, 1},
		{1, 1},
	}),
}

func main() {
	if err := star1(); err != nil {
		panic(err)
	}
	if err := star2(); err != nil {
		panic(err)
	}
}

func makeIterator[T any](slice []T) func() T {
	var i int
	l := len(slice)
	return func() T {
		defer func() { i++ }()
		return slice[i%l]
	}
}

func star1() error {
	input, err := readInput("input-example.txt")
	if err != nil {
		return err
	}

	g := grid.NewDense(7, Star1MaxBlocks*3+5)

	nextInput := makeIterator(input)
	nextBlock := makeIterator(Shapes)

	for blocksSpawned := 0; blocksSpawned < Star1MaxBlocks; blocksSpawned++ {
		var (
			block       = nextBlock()
			spawnPoint  = grid.Point{2, g.HighestY() + 3}
			activeBlock = block.SpawnAt(g, spawnPoint)
		)

		for {
			activeBlock.TryShift(nextInput())
			if !activeBlock.TryShift(grid.UnitDown) {
				break
			}
		}

		g.AddSpawned(activeBlock)
	}

	fmt.Println(g.HighestY())

	return nil
}

func star2() error {
	const (
		gridHeight       = 1000000
		minPeriodGuess   = 10
		maxPeriodGuess   = 10000
		repeatConfidence = 10
	)

	input, err := readInput("input.txt")
	if err != nil {
		return err
	}

	g := grid.NewDense(7, gridHeight)

	nextInput := makeIterator(input)
	nextBlock := makeIterator(Shapes)

	dropHeightCount := make([]int, gridHeight)

	var found bool
	var adjust int

	for blocksSpawned := 0; blocksSpawned < Star2MaxBlocks; blocksSpawned++ {
		var (
			block       = nextBlock()
			spawnPoint  = grid.Point{2, g.HighestY() + 3}
			activeBlock = block.SpawnAt(g, spawnPoint)
		)

		for {
			activeBlock.TryShift(nextInput())
			if !activeBlock.TryShift(grid.UnitDown) {
				break
			}
		}

		g.AddSpawned(activeBlock)

		if !found {
			dropHeightCount[blocksSpawned] = g.HighestY()
		}

		for periodGuess := minPeriodGuess; periodGuess < maxPeriodGuess; periodGuess++ {
			if found || blocksSpawned-periodGuess < 0 {
				break
			}
			diff := dropHeightCount[blocksSpawned] - dropHeightCount[blocksSpawned-periodGuess]
			for repeats := 0; blocksSpawned > (repeats+1)*periodGuess; repeats++ {
				if dropHeightCount[blocksSpawned-repeats*periodGuess] != dropHeightCount[blocksSpawned-(repeats+1)*periodGuess]+diff {
					break
				}
				if repeats > repeatConfidence {
					found = true
					toGoal := Star2MaxBlocks - blocksSpawned
					periodsToGoal := toGoal / periodGuess
					blocksSpawned += periodsToGoal * periodGuess
					adjust = periodsToGoal * diff
					fmt.Println("found period", periodGuess, diff)
					goto periodFound
				}
			}
		}
	periodFound:
	}

	fmt.Println(g.HighestY() + adjust)

	return nil
}

func readInput(filename string) ([]grid.Point, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	dirs := make([]grid.Point, 0, len(fileBytes))
	for _, b := range fileBytes {
		switch b {
		case '<':
			dirs = append(dirs, grid.UnitLeft)
		case '>':
			dirs = append(dirs, grid.UnitRight)
		case '\r', '\n':
			break
		default:
			return nil, fmt.Errorf("unexpected character %q", b)
		}
	}

	return dirs, nil
}
