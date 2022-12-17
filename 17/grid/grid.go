package grid

import (
	"fmt"
	"reflect"
)

type State byte

const (
	StateNil  State = iota
	StateRock       // #
)

var (
	UnitLeft  = Point{-1, 0}
	UnitRight = Point{1, 0}
	UnitDown  = Point{0, -1}
)

type Point struct{ X, Y int }

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

type Points []Point

func (p Points) Add(q Point) Points {
	for i := range p {
		p[i] = p[i].Add(q)
	}
	return p
}

type Dense struct {
	maxX, maxY int
	data       []State
	highestY   int
}

func NewDense(maxX, maxY int) *Dense {
	return &Dense{
		maxX:     maxX,
		maxY:     maxY,
		data:     make([]State, maxX*maxY),
		highestY: 0,
	}
}

func (d *Dense) Get(p Point) State {
	return d.data[p.Y*d.maxX+p.X]
}

func (d *Dense) Set(p Point, state State) {
	if p.Y+1 > d.highestY {
		d.highestY = p.Y + 1
	}
	d.data[p.Y*d.maxX+p.X] = state
}

func (d *Dense) AddSpawned(b *SpawnedBlock) {
	for _, p := range b.points {
		if d.Get(p) != StateNil {
			panic(fmt.Sprintf("collision at %v", p))
		}
		d.Set(p, StateRock)
	}
}

func (d *Dense) HighestY() int { return d.highestY }

func (d *Dense) InBounds(p Point) bool {
	return p.X >= 0 && p.X < d.maxX && p.Y >= 0 && p.Y < d.maxY
}

func (d *Dense) Period() bool {
	hy := d.highestY
	if hy == 0 {
		return false
	}
	//fmt.Printf("comparing %d: data[:%d] == data[%d:%d]\n", hy, (hy/2)*d.maxX, (hy/2)*d.maxX, (hy)*d.maxX)
	return reflect.DeepEqual(d.data[:(hy/2)*d.maxX], d.data[(hy/2)*d.maxX:(hy)*d.maxX])
}

func (d *Dense) Period2() int {
	hy := d.highestY
	if hy < 3 {
		return -1
	}

	depth := 30
	if hy/2 < depth {
		depth = hy / 2
	}
	for i := 0; i < hy/2; i++ {
		if reflect.DeepEqual(
			d.data[(hy-(i*2+1))*d.maxX:(hy-i)*d.maxX],
			d.data[(hy-i)*d.maxX:(hy+1)*d.maxX],
		) {
			return i
		}
	}
	return -1
}

func (d *Dense) SafeHeight() int {
	bitMask := 0b0000000

	for y := d.highestY - 1; y >= 0; y-- {
		for x := 0; x < d.maxX; x++ {
			if d.Get(Point{x, y}) != StateNil {
				bitMask |= 1 << x
			}
			if bitMask == 0b01111111 {
				return y
			}
		}
	}

	return -1
}

func (d *Dense) Consolidate() int {
	safeHeight := d.SafeHeight()
	if safeHeight == -1 {
		return 0
	}
	copy(d.data[:safeHeight*d.maxX], d.data[safeHeight*d.maxX:])

	for i := (d.maxY - safeHeight) * d.maxX; i < d.maxY*d.maxX; i++ {
		d.data[i] = StateNil
	}

	d.highestY -= safeHeight
	return d.SafeHeight()
}

type Block struct {
	maxY, maxX int
	data       []bool
}

func NewBlock(data [][]byte) *Block {
	maxY := len(data)
	maxX := len(data[0])
	b := &Block{maxY, maxX, make([]bool, maxY*maxX)}
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			b.data[y*maxX+x] = data[maxY-y-1][x] > 0
		}
	}
	return b
}

func (b *Block) Get(p Point) bool {
	return b.data[p.Y*b.maxX+p.X]
}

func (b *Block) Points() Points {
	var points Points
	for y := 0; y < b.maxY; y++ {
		for x := 0; x < b.maxX; x++ {
			if b.data[y*b.maxX+x] {
				points = append(points, Point{x, y})
			}
		}
	}
	return points
}

func (b *Block) SpawnAt(g *Dense, p Point) *SpawnedBlock {
	blockPoints := b.Points()
	blockPoints = blockPoints.Add(p)

	return &SpawnedBlock{g, blockPoints}
}

type SpawnedBlock struct {
	g      *Dense
	points Points
}

func (s *SpawnedBlock) TryShift(d Point) bool {
	for _, p := range s.points {
		shiftedP := p.Add(d)
		if !s.g.InBounds(shiftedP) || s.g.Get(shiftedP) != StateNil {
			return false
		}
	}
	s.points = s.points.Add(d)
	return true
}
