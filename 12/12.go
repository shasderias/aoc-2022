package main

import (
	"aoc-2022-12/grid"
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

const inputFilePath = "input.txt"

func main() {
	star1()
	star2()
}

type CoordType int

const (
	CoordTypeNil = iota
	CoordTypeStart
	CoordTypeEnd
	CoordTypeCell
)

type Coord struct {
	Elevation int
	Type      CoordType
	Step      int
}

type dbgFmt struct {
	f io.Writer
}

func (d dbgFmt) Println(a ...any) (int, error) {
	return fmt.Fprintln(d.f, a...)
}

func (d dbgFmt) Printf(format string, a ...any) (int, error) {
	return fmt.Fprintf(d.f, format, a...)
}

var dbg = dbgFmt{f: io.Discard}

func star1() error {
	g, pointStart, pointEnd, err := parseGrid()
	if err != nil {
		return err
	}

	fillGrid(g, pointStart, 0)

	fmt.Println(g.Get(pointEnd).Step)

	return nil
}

func star2() error {
	g, _, pointEnd, err := parseGrid()
	if err != nil {
		return err
	}

	lowestPoints := grid.PointSet{}

	g.Walk(func(p grid.Point, c Coord) {
		if c.Elevation == 0 {
			lowestPoints.Add(p)
		}
	})

	fewestSteps := math.MaxInt

	for p := range lowestPoints {
		fillGrid(g, p, 0)
		stepsTaken := g.Get(pointEnd).Step
		if stepsTaken != 00 && stepsTaken < fewestSteps {
			fewestSteps = stepsTaken
		}
	}

	fmt.Println(fewestSteps)

	return nil
}

func parseGrid() (g *grid.Grid[Coord], pointStart, pointEnd grid.Point, err error) {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return g, pointStart, pointEnd, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	maxY := len(lines)
	maxX := len(lines[0])

	dbg.Println("max:", maxX, maxY)
	g = grid.New[Coord](maxX, maxY)

	for y := range lines {
		for x := range lines[y] {
			point := grid.Point{x, y}
			var (
				elevation int
				coordType CoordType
			)

			switch {
			case lines[y][x] == 'S':
				elevation = 0
				coordType = CoordTypeStart
				pointStart = point
			case lines[y][x] == 'E':
				elevation = 'z' - 'a'
				coordType = CoordTypeEnd
				pointEnd = point
			case lines[y][x] >= 'a' && lines[y][x] <= 'z':
				elevation = int(lines[y][x] - 'a')
				coordType = CoordTypeCell
			default:
				panic(fmt.Sprintf("unknown char: %c", lines[y][x]))
			}

			g.Set(point, Coord{Elevation: elevation, Type: coordType})
		}
	}

	return g, pointStart, pointEnd, nil
}

func fillGrid(g *grid.Grid[Coord], t grid.Point, steps int) {
	var (
		currentTargets = grid.PointSet{t: {}}
		nextTargets    = grid.PointSet{}
	)

	for {
		dbg.Printf("cur[%d]: %s\n", steps, currentTargets)
		for t := range currentTargets {
			targets := findTargets(g, t)
			dbg.Printf("targets for %v are %v\n", t, targets)
			nextTargets.Merge(targets)
		}
		dbg.Println("nex", nextTargets)
		if len(nextTargets) == 0 {
			break
		}
		for t := range nextTargets {
			mark(g, t, steps+1)
		}
		currentTargets = nextTargets
		nextTargets = map[grid.Point]struct{}{}
		steps++
		dbg.Println("")
	}
}

func findTargets(g *grid.Grid[Coord], t grid.Point) grid.PointSet {
	var (
		targets  = grid.PointSet{}
		curCoord = g.Get(t)
	)

	for _, dir := range grid.CardinalDirs {
		nt := t.Add(dir.Unit())
		if !g.IsInBounds(nt) {
			dbg.Println("skip bounds", nt, nt)
			continue
		}

		nextCoord := g.Get(nt)
		if nextCoord.Type == CoordTypeStart {
			continue
		}
		if nextCoord.Step != 0 && nextCoord.Step <= curCoord.Step {
			dbg.Println("skip step", nt, nextCoord)
			continue
		}
		if nextCoord.Elevation-curCoord.Elevation > 1 {
			dbg.Println("skip elevation", nt, curCoord, nextCoord)
			continue
		}

		targets.Add(nt)
	}

	return targets
}

func mark(g *grid.Grid[Coord], t grid.Point, steps int) {
	curCoord := g.Get(t)
	curCoord.Step = steps
	g.Set(t, curCoord)
}
