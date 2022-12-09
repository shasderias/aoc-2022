package grid

import (
	"fmt"
	"math"
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

type Point struct {
	X, Y int
}

// DirOf b, when at a.
// Returns DirNil if a and b overlap.
func (a Point) DirOf(b Point) Dir {
	switch {
	case b.X == a.X && b.Y > a.Y:
		return DirN
	case b.X > a.X && b.Y > a.Y:
		return DirNE
	case b.X > a.X && b.Y == a.Y:
		return DirE
	case b.X > a.X && b.Y < a.Y:
		return DirSE
	case b.X == a.X && b.Y < a.Y:
		return DirS
	case b.X < a.X && b.Y < a.Y:
		return DirSW
	case b.X < a.X && b.Y == a.Y:
		return DirW
	case b.X < a.X && b.Y > a.Y:
		return DirNW
	}
	return DirNil
}

func (a Point) Add(b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y}
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
