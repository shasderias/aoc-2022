package main

import (
	"aoc-2022-09/grid"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const inputFilePath = "input.txt"

var (
	MoveNil = grid.Point{0, 0}
	MoveN   = grid.Point{0, 1}
	MoveNE  = grid.Point{1, 1}
	MoveE   = grid.Point{1, 0}
	MoveSE  = grid.Point{1, -1}
	MoveS   = grid.Point{0, -1}
	MoveSW  = grid.Point{-1, -1}
	MoveW   = grid.Point{-1, 0}
	MoveNW  = grid.Point{-1, 1}
)

var (
	moveMap = map[grid.Dir]grid.Point{
		grid.DirN:   MoveN,
		grid.DirNE:  MoveNE,
		grid.DirS:   MoveS,
		grid.DirSE:  MoveSE,
		grid.DirE:   MoveE,
		grid.DirSW:  MoveSW,
		grid.DirW:   MoveW,
		grid.DirNW:  MoveNW,
		grid.DirNil: MoveNil,
	}
)

type cmd struct {
	dir grid.Dir
	n   int
}

func parseLine(line string) cmd {
	lineParts := strings.Split(line, " ")
	n, err := strconv.Atoi(lineParts[1])
	if err != nil {
		panic(err)
	}

	switch lineParts[0] {
	case "U":
		return cmd{grid.DirN, n}
	case "D":
		return cmd{grid.DirS, n}
	case "L":
		return cmd{grid.DirW, n}
	case "R":
		return cmd{grid.DirE, n}
	default:
		panic(fmt.Sprintf("unknown direction: %s", lineParts[0]))
	}
}

func main() {
	star1()
	star2()
}

func star1() error {
	head := grid.Point{0, 0}
	tail := grid.Point{0, 0}

	visited := map[grid.Point]struct{}{}

	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		cmd := parseLine(scanner.Text())
		for i := 0; i < cmd.n; i++ {
			visited[tail] = struct{}{}

			switch cmd.dir {
			case grid.DirN:
				head = head.Add(MoveN)
			case grid.DirS:
				head = head.Add(MoveS)
			case grid.DirE:
				head = head.Add(MoveE)
			case grid.DirW:
				head = head.Add(MoveW)
			default:
				panic(fmt.Sprintf("unsupported command direction: %s", cmd.dir))
			}

			switch tail.DirOf(head) {
			case grid.DirN, grid.DirS, grid.DirE, grid.DirW:
				if tail.Distance(head) <= 1 {
					continue
				}
			case grid.DirNE, grid.DirSE, grid.DirSW, grid.DirNW:
				if tail.Distance(head) <= 2 {
					continue
				}
			}
			tail = tail.Add(moveMap[tail.DirOf(head)])

		}
	}

	visited[tail] = struct{}{}

	visitedPos := 0
	for range visited {
		visitedPos++
	}

	fmt.Println(visitedPos)
	return nil
}

func star2() error {
	rope := [10]grid.Point{}

	visited := map[grid.Point]struct{}{}

	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	visited[rope[9]] = struct{}{}

	for scanner.Scan() {
		line := scanner.Text()
		cmd := parseLine(line)
		for i := 0; i < cmd.n; i++ {
			// move head
			switch cmd.dir {
			case grid.DirN:
				rope[0] = rope[0].Add(MoveN)
			case grid.DirS:
				rope[0] = rope[0].Add(MoveS)
			case grid.DirE:
				rope[0] = rope[0].Add(MoveE)
			case grid.DirW:
				rope[0] = rope[0].Add(MoveW)
			default:
				panic(fmt.Sprintf("unsupported command direction: %s", cmd.dir))
			}

			// compute tail movements
			for j := 0; j < len(rope)-1; j++ {
				rope[j+1] = calcPos(rope[j], rope[j+1])
			}

			visited[rope[9]] = struct{}{}

			//visualize(rope)
		}
	}

	visitedPos := 0
	for range visited {
		visitedPos++
	}

	fmt.Println(visitedPos)
	return nil
}

func calcPos(head, tail grid.Point) grid.Point {
	switch tail.DirOf(head) {
	case grid.DirN, grid.DirS, grid.DirE, grid.DirW:
		if tail.Distance(head) <= 1 {
			return tail
		}
	case grid.DirNE, grid.DirSE, grid.DirSW, grid.DirNW:
		if tail.Distance(head) <= 2 {
			return tail
		}
	}
	return tail.Add(moveMap[tail.DirOf(head)])
}

func visualize(rope [10]grid.Point) {
	const (
		xMin = -25
		xMax = 25
		yMin = -25
		yMax = 25
	)
	for y := yMax; y >= yMin; y-- {
		for x := xMin; x < xMax; x++ {
			p := grid.Point{x, y}
			for i := range rope {
				if rope[i] == p {
					if i == 0 {
						fmt.Print("H")
					} else {
						fmt.Printf("%d", i)
					}
					goto next
				}
			}
			fmt.Print(".")
		next:
		}
		fmt.Println()
	}
	fmt.Println("")
}
