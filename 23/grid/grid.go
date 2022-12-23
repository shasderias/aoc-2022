package grid

import "math"

var (
	N = Vec{0, -1}
	S = Vec{0, 1}
	E = Vec{1, 0}
	W = Vec{-1, 0}
)

var (
	NE = Vec{1, -1}
	SE = Vec{1, 1}
	NW = Vec{-1, -1}
	SW = Vec{-1, 1}
)

type Vec struct{ X, Y int }

func (v Vec) Add(o Vec) Vec { return Vec{v.X + o.X, v.Y + o.Y} }

type Sparse struct {
	g map[Vec]struct{}
}

func NewSparse() Sparse {
	return Sparse{g: make(map[Vec]struct{})}
}

func (s Sparse) Set(v Vec) {
	s.g[v] = struct{}{}
}
func (s Sparse) Get(v Vec) bool {
	_, ok := s.g[v]
	return ok
}
func (s Sparse) Clear(v Vec) {
	delete(s.g, v)
}
func (s Sparse) Walk(f func(v Vec)) {
	for v := range s.g {
		f(v)
	}
}
func (s Sparse) Bounds() (min, max Vec) {
	min.X, min.Y = math.MaxInt, math.MaxInt

	for v := range s.g {
		if v.X < min.X {
			min.X = v.X
		}
		if v.Y < min.Y {
			min.Y = v.Y
		}
		if v.X > max.X {
			max.X = v.X
		}
		if v.Y > max.Y {
			max.Y = v.Y
		}
	}
	return
}
