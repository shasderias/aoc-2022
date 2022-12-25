package grid

import (
	"fmt"
	"strings"
)

type Tile string

const (
	TileNil       Tile = "."
	TileWall           = "#"
	TileBlizzardN      = "^"
	TileBlizzardS      = "v"
	TileBlizzardE      = ">"
	TileBlizzardW      = "<"
)

type Vec struct{ X, Y int }

func (v Vec) Add(o Vec) Vec {
	return Vec{v.X + o.X, v.Y + o.Y}
}

func (v Vec) MDist(o Vec) int {
	return abs(v.X-o.X) + abs(v.Y-o.Y)
}

func (v Vec) CardinalString() string {
	switch v {
	case N:
		return "N"
	case S:
		return "S"
	case E:
		return "E"
	case W:
		return "W"
	}
	panic(fmt.Sprintf("%v is not a cardinal direction", v))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var (
	N = Vec{0, -1}
	S = Vec{0, 1}
	E = Vec{1, 0}
	W = Vec{-1, 0}
)

var Cardinal = []Vec{N, S, E, W}

type Dense struct {
	MaxX, MaxY int
	data       []Tile
}

func New(maxX, maxY int) *Dense {
	return &Dense{
		MaxX: maxX,
		MaxY: maxY,
		data: make([]Tile, maxX*maxY),
	}
}

func (d *Dense) Set(v Vec, t Tile) {
	d.data[v.Y*d.MaxX+v.X] = t
}

func (d *Dense) Get(v Vec) Tile {
	return d.data[v.Y*d.MaxX+v.X]
}

func (d *Dense) InBounds(v Vec) bool {
	return v.X >= 0 && v.X < d.MaxX && v.Y >= 0 && v.Y < d.MaxY
}

func (d *Dense) InInBounds(v Vec) bool {
	return v.X > 0 && v.X < d.MaxX-1 && v.Y > 0 && v.Y < d.MaxY-1
}

func (d *Dense) String() string {
	buf := strings.Builder{}
	buf.Grow(d.MaxX*d.MaxY + d.MaxY)

	for y := 0; y < d.MaxY; y++ {
		for x := 0; x < d.MaxX; x++ {
			buf.WriteString(string(d.Get(Vec{x, y})))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func (d *Dense) Walk(fn func(v Vec, t Tile)) {
	for y := 0; y < d.MaxY; y++ {
		for x := 0; x < d.MaxX; x++ {
			v := Vec{x, y}
			fn(v, d.Get(v))
		}
	}
}

type Blizzard struct {
	Pos  Vec
	Type Tile
}

func (b Blizzard) At(maxX, maxY int, minute int) Vec {
	maxX -= 2
	maxY -= 2
	x := b.Pos.X - 1
	y := b.Pos.Y - 1

	switch b.Type {
	case TileBlizzardN:
		y = (((y-minute)%maxY)+maxY)%maxY + 1
		return Vec{b.Pos.X, y}
	case TileBlizzardS:
		y = (((y+minute)%maxY)+maxY)%maxY + 1
		return Vec{b.Pos.X, y}
	case TileBlizzardE:
		x = (((x+minute)%maxX)+maxX)%maxX + 1
		return Vec{x, b.Pos.Y}
	case TileBlizzardW:
		x = (((x-minute)%maxX)+maxX)%maxX + 1
		return Vec{x, b.Pos.Y}
	default:
		panic(fmt.Sprintf("unknown blizzard type %v", b.Type))
	}
}
