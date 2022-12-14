package grid

import (
	"fmt"
	"math"
	"strings"
)

type Dir int

const (
	DirNil Dir = iota
	DirN
	DirNE
	DirE
	DirSE
	DirS
	DirSW
	DirW
	DirNW
)

var CardinalDirs = []Dir{DirN, DirE, DirS, DirW}

func (d Dir) String() string {
	switch d {
	case DirNil:
		return "nil"
	case DirN:
		return "N"
	case DirNE:
		return "NE"
	case DirE:
		return "E"
	case DirSE:
		return "SE"
	case DirS:
		return "S"
	case DirSW:
		return "SW"
	case DirW:
		return "W"
	case DirNW:
		return "NW"
	default:
		panic(fmt.Sprintf("unknown direction: %d", d))
	}
}

func (d Dir) Unit() Point {
	switch d {
	case DirN:
		return Point{0, -1}
	case DirNE:
		return Point{1, -1}
	case DirE:
		return Point{1, 0}
	case DirSE:
		return Point{1, 1}
	case DirS:
		return Point{0, 1}
	case DirSW:
		return Point{-1, 1}
	case DirW:
		return Point{-1, 0}
	case DirNW:
		return Point{-1, -1}
	default:
		panic(fmt.Sprintf("unknown direction: %d", d))
	}
}

type Point struct {
	X, Y int
}

// DirOf b, when at a.
// Returns DirNil if a and b overlap.
func (a Point) DirOf(b Point) Dir {
	switch {
	case b.X == a.X && b.Y > a.Y:
		return DirS
	case b.X > a.X && b.Y > a.Y:
		return DirSE
	case b.X > a.X && b.Y == a.Y:
		return DirE
	case b.X > a.X && b.Y < a.Y:
		return DirNE
	case b.X == a.X && b.Y < a.Y:
		return DirN
	case b.X < a.X && b.Y < a.Y:
		return DirNW
	case b.X < a.X && b.Y == a.Y:
		return DirW
	case b.X < a.X && b.Y > a.Y:
		return DirSW
	}
	return DirNil
}

func (a Point) Add(b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y}
}

func (p Point) Mul(n int) Point {
	return Point{p.X * n, p.Y * n}
}

func (a Point) Distance(b Point) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}

type Grid[T any] struct {
	data       []T
	maxX, maxY int
}

func New[T any](maxX, maxY int) *Grid[T] {
	return &Grid[T]{
		data: make([]T, maxX*maxY),
		maxX: maxX,
		maxY: maxY,
	}
}

func (g *Grid[T]) Set(p Point, v T) {
	g.data[p.Y*g.maxX+p.X] = v
}

func (g *Grid[T]) Get(p Point) T {
	return g.data[p.Y*g.maxX+p.X]
}

func (g *Grid[T]) IsInBounds(p Point) bool {
	return p.X >= 0 && p.X < g.maxX && p.Y >= 0 && p.Y < g.maxY
}

func (g *Grid[T]) Walk(f func(p Point, v T)) {
	for y := 0; y < g.maxY; y++ {
		for x := 0; x < g.maxX; x++ {
			p := Point{x, y}
			f(p, g.Get(p))
		}
	}
}

type PointSet map[Point]struct{}

func (ps PointSet) Add(p Point) { ps[p] = struct{}{} }

func (ps PointSet) String() string {
	l := len(ps)
	s := strings.Builder{}
	s.WriteString("[")
	i := 0
	for p := range ps {
		s.WriteString(fmt.Sprintf("%v", p))
		if i < l-1 {
			s.WriteString(", ")
		}
		i++
	}
	s.WriteString("]")
	return s.String()
}

func (ps1 PointSet) Merge(ps2 PointSet) {
	for p := range ps2 {
		ps1[p] = struct{}{}
	}
}

type Sparse[T any] struct {
	data map[Point]T
}

func NewSparse[T any]() *Sparse[T] {
	return &Sparse[T]{
		data: make(map[Point]T),
	}
}

func (s *Sparse[T]) Set(p Point, v T) {
	s.data[p] = v
}

func (s *Sparse[T]) Get(p Point) T {
	return s.data[p]
}

func (s *Sparse[T]) Walk(f func(p Point, v T)) {
	for p, v := range s.data {
		f(p, v)
	}
}
