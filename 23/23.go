package main

import (
	"bufio"
	"fmt"
	"os"

	"aoc-2022-23/grid"
)

func main() {
	star1()
	star2()
}

func star1() {
	g := readInput("input.txt")

	round(g, 10)

	min, max := g.Bounds()

	emptyTiles := 0
	for y := min.Y; y <= max.Y; y++ {
		for x := min.X; x <= max.X; x++ {
			if !g.Get(grid.Vec{X: x, Y: y}) {
				emptyTiles++
			}
		}
	}
	fmt.Println(emptyTiles)
}

func star2() {
	g := readInput("input.txt")

	n := roundTerm(g)
	fmt.Println(n + 1)
}

func printGrid(g grid.Sparse) {
	min, max := g.Bounds()

	for y := min.Y; y <= max.Y; y++ {
		for x := min.X; x <= max.X; x++ {
			if g.Get(grid.Vec{X: x, Y: y}) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

type moveEvalFunc func(g grid.Sparse, v grid.Vec) *grid.Vec

func getMoveEvalFuncs() (*[]moveEvalFunc, func()) {
	moveEvalFuncs := []moveEvalFunc{
		func(g grid.Sparse, v grid.Vec) *grid.Vec {
			if !g.Get(v.Add(grid.N)) && !g.Get(v.Add(grid.NE)) && !g.Get(v.Add(grid.NW)) {
				nv := v.Add(grid.N)
				return &nv
			}
			return nil
		},
		func(g grid.Sparse, v grid.Vec) *grid.Vec {
			if !g.Get(v.Add(grid.S)) && !g.Get(v.Add(grid.SE)) && !g.Get(v.Add(grid.SW)) {
				nv := v.Add(grid.S)
				return &nv
			}
			return nil
		},
		func(g grid.Sparse, v grid.Vec) *grid.Vec {
			if !g.Get(v.Add(grid.W)) && !g.Get(v.Add(grid.NW)) && !g.Get(v.Add(grid.SW)) {
				nv := v.Add(grid.W)
				return &nv
			}
			return nil
		},
		func(g grid.Sparse, v grid.Vec) *grid.Vec {
			if !g.Get(v.Add(grid.E)) && !g.Get(v.Add(grid.NE)) && !g.Get(v.Add(grid.SE)) {
				nv := v.Add(grid.E)
				return &nv
			}
			return nil
		},
	}
	rotate := func() {
		moveEvalFuncs = append(moveEvalFuncs[1:], moveEvalFuncs[0])
	}

	return &moveEvalFuncs, rotate
}

func allAdjacentEmpty(g grid.Sparse, v grid.Vec) bool {
	if !g.Get(v.Add(grid.N)) && !g.Get(v.Add(grid.NE)) && !g.Get(v.Add(grid.NW)) &&
		!g.Get(v.Add(grid.S)) && !g.Get(v.Add(grid.SE)) && !g.Get(v.Add(grid.SW)) &&
		!g.Get(v.Add(grid.W)) && !g.Get(v.Add(grid.E)) {
		return true
	}
	return false
}

type proposal struct {
	from, to grid.Vec
	n        int
}

func propose(g grid.Sparse, moveEvalFuncs *[]moveEvalFunc) map[grid.Vec]*proposal {
	proposals := map[grid.Vec]*proposal{}

	g.Walk(func(v grid.Vec) {
		if allAdjacentEmpty(g, v) {
			return
		}
		for _, f := range *moveEvalFuncs {
			if nv := f(g, v); nv != nil {
				p := proposals[*nv]
				if p == nil {
					p = &proposal{from: v, to: *nv}
					proposals[*nv] = p
				}
				p.n++
				return
			}
		}
	})

	return proposals
}

func move(g grid.Sparse, proposals map[grid.Vec]*proposal) {
	for _, p := range proposals {
		if p.n == 1 {
			g.Clear(p.from)
			g.Set(p.to)
		}
	}
}

func round(g grid.Sparse, n int) {
	moveEvalFuncs, rotate := getMoveEvalFuncs()
	for i := 0; i < n; i++ {
		proposals := propose(g, moveEvalFuncs)
		move(g, proposals)
		rotate()
	}
}

func roundTerm(g grid.Sparse) int {
	moveEvalFuncs, rotate := getMoveEvalFuncs()
	for i := 0; ; i++ {
		proposals := propose(g, moveEvalFuncs)
		if len(proposals) == 0 {
			return i
		}
		move(g, proposals)
		rotate()
	}
}

func readInput(inputPath string) grid.Sparse {
	f, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	g := grid.NewSparse()

	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()
		if line == "" {
			break
		}
		for x, c := range line {
			switch c {
			case '#':
				g.Set(grid.Vec{X: x, Y: y})
			case '.':
			// do nothing
			default:
				panic(fmt.Sprintf("unexpected char %q", c))
			}
		}
	}

	return g
}
