package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"aoc-2022-22/grid"
)

func main() {
	star1("input.txt")
	star2("input.txt")
}

func star1(inputPath string) {
	maze, insts, err := parseInput(inputPath)
	if err != nil {
		panic(err)
	}

	curCoord := maze.FindStart()
	curDir := grid.DirE

	for _, in := range insts {
		switch in.Type() {
		case instTypeMove:
			for i := 0; i < in.move; i++ {
				nextCoord, nextTile := maze.Next(curCoord, curDir)
				if nextTile != grid.TileFloor {
					break
				}
				curCoord = nextCoord
			}
		case instTypeTurn:
			switch in.turn {
			case instDirR:
				curDir = curDir.Rotate(grid.RotateCW)
			case instDirL:
				curDir = curDir.Rotate(grid.RotateCCW)
			}
		}
	}

	fmt.Println((curCoord.X+1)*4 + (curCoord.Y+1)*1000 + curDir.PasswordValue())
}

func star2(inputPath string) {
	maze, insts, err := parseInput(inputPath)
	if err != nil {
		panic(err)
	}

	edgePairs := []grid.EdgePair{
		grid.NewEdgePair( // 1N - 6W
			"1N - 6W",
			grid.Vec{50, -1}, grid.Vec{99, -1},
			grid.Vec{-1, 150}, grid.Vec{-1, 199},
			grid.RotateCW),
		grid.NewEdgePair( // 2N to 6S
			"2N - 6S",
			grid.Vec{100, -1}, grid.Vec{149, -1},
			grid.Vec{0, 200}, grid.Vec{49, 200},
			grid.RotateNil),
		grid.NewEdgePair( //1W to 4W
			"1W - 4W",
			grid.Vec{49, 0}, grid.Vec{49, 49},
			grid.Vec{-1, 149}, grid.Vec{-1, 100},
			grid.Rotate180),
		grid.NewEdgePair( // 2E to 5E
			"2E - 5E",
			grid.Vec{150, 0}, grid.Vec{150, 49},
			grid.Vec{100, 149}, grid.Vec{100, 100},
			grid.Rotate180),
		grid.NewEdgePair( // 3W to 4N
			"3W - 4N",
			grid.Vec{49, 50}, grid.Vec{49, 99},
			grid.Vec{0, 99}, grid.Vec{49, 99},
			grid.RotateCCW),
		grid.NewEdgePair( // 3E to 2S
			"3E - 2S",
			grid.Vec{100, 50}, grid.Vec{100, 99},
			grid.Vec{100, 50}, grid.Vec{149, 50},
			grid.RotateCCW),
		grid.NewEdgePair( // 5S to 6E
			"5S - 6E",
			grid.Vec{50, 150}, grid.Vec{99, 150},
			grid.Vec{50, 150}, grid.Vec{50, 199},
			grid.RotateCW),
	}

	err = maze.LoadAndValidateEdgePairs(edgePairs)
	if err != nil {
		panic(err)
	}

	curCoord := maze.FindStart()
	curDir := grid.DirE

	for _, in := range insts {
		switch in.Type() {
		case instTypeMove:
			for i := 0; i < in.move; i++ {
				nextCoord, nextDir, nextTile := maze.NextCube(curCoord, curDir)
				if nextTile != grid.TileFloor {
					break
				}
				curCoord = nextCoord
				curDir = nextDir
			}
		case instTypeTurn:
			switch in.turn {
			case instDirR:
				curDir = curDir.Rotate(grid.RotateCW)
			case instDirL:
				curDir = curDir.Rotate(grid.RotateCCW)
			}
		}
	}

	fmt.Println((curCoord.X+1)*4 + (curCoord.Y+1)*1000 + curDir.PasswordValue())
}

type instDir string

const (
	instDirR = "R"
	instDirL = "L"
)

type instType byte

const (
	instTypeMove instType = iota
	instTypeTurn
)

type inst struct {
	move int
	turn instDir
}

func (i inst) Type() instType {
	if i.turn == "" {
		return instTypeMove
	}
	return instTypeTurn
}

func parseInput(inputPath string) (*grid.Grid, []inst, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	maxX := 0
	mazeLines := []string{}

	// read maze
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}

		mazeLines = append(mazeLines, line)
		if len(line) > maxX {
			maxX = len(line)
		}
	}

	maxY := len(mazeLines)

	scanner.Scan()

	instLine := scanner.Text()

	maze := grid.NewGrid(maxX, maxY)

	for y, line := range mazeLines {
		for x, c := range line {
			switch c {
			case '#':
				maze.Set(grid.Vec{x, y}, grid.TileWall)
			case '.':
				maze.Set(grid.Vec{x, y}, grid.TileFloor)
			case ' ':
				// do nothing
			default:
				panic(fmt.Sprintf("invalid tile %c", c))
			}
		}
	}

	insts, err := parseInst(instLine)
	if err != nil {
		return nil, nil, err
	}

	return maze, insts, nil
}

func parseInst(line string) ([]inst, error) {
	insts := []inst{}

	buf := ""

	for i := 0; i < len(line); i++ {
		c := line[i]
		switch c {
		case 'R', 'L':
			if len(buf) > 0 {
				steps, err := strconv.Atoi(buf)
				if err != nil {
					return nil, err
				}
				insts = append(insts, inst{
					move: steps,
				})
				insts = append(insts, inst{
					turn: instDir(c),
				})
				buf = ""
			}
		default:
			buf += string(c)
		}
	}

	if buf != "" {
		steps, err := strconv.Atoi(buf)
		if err != nil {
			return nil, err
		}
		insts = append(insts, inst{
			move: steps,
		})
	}

	return insts, nil
}
