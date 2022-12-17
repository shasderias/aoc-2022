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
	const gridHeight = 100000

	input, err := readInput("input-example.txt")
	if err != nil {
		return err
	}

	g := grid.NewDense(7, gridHeight)

	nextInput := makeIterator(input)
	nextBlock := makeIterator(Shapes)

	//multiple := len(input) * len(Shapes)
	dropHeightCount := make([]int, gridHeight)
	//dropHeightCountMap := make(map[int]int)
	//acc := 0

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

		//if gridHeight-g.HighestY() < 1000 {
		//	acc += g.Consolidate()
		//}

		//if blocksSpawned%multiple-1 == 0 {
		//	dropHeightCountMap[blocksSpawned] = g.HighestY() + acc
		//	if periodFound(dropHeightCountMap) {
		//		fmt.Println(blocksSpawned)
		//		break
		//	}
		//
		//}
		//
		dropHeightCount[blocksSpawned] = g.HighestY()

		maxPeriodGuess := min(40, blocksSpawned/2)
		if blocksSpawned > 3000 {
			for i := 3; i < maxPeriodGuess; i++ {
				periodDiff := dropHeightCount[blocksSpawned] - dropHeightCount[blocksSpawned-i]
				repeats := 0
				for ; blocksSpawned > repeats*i; repeats++ {
					if dropHeightCount[blocksSpawned-repeats*i] != dropHeightCount[blocksSpawned-(repeats+1)*i]+periodDiff {
						break
					}
					if repeats > 10 {
						fmt.Println("found", i, blocksSpawned, repeats)
						break
					}
				}
			}
		}

		//
		//if dropHeightCount[blocksSpawned] == dropHeightCount[blocksSpawned/2]/2 {
		//	fmt.Println("period at", blocksSpawned, " - ", g.HighestY())
		//}
		//
		//if g.Period() {
		//	fmt.Println("period at", blocksSpawned, " - ", g.HighestY())
		//}
		//if period := g.Period2(); period != -1 {
		//	fmt.Println("period found", period)
		//	return nil
		//}
	}
	//fmt.Println(acc)
	fmt.Println(g.HighestY())

	return nil
}

func periodFound(countMap map[int]int) bool {
	for blocksDropped, height := range countMap {
		if doubleHeight, ok := countMap[blocksDropped*2]; ok && doubleHeight == height*2 {
			return true
		}
	}
	return false

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
