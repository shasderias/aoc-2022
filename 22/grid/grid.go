package grid

import "fmt"

var (
	UnitN = Vec{0, -1}
	UnitE = Vec{1, 0}
	UnitS = Vec{0, 1}
	UnitW = Vec{-1, 0}
)

type Vec struct {
	X, Y int
}

func (p Vec) Add(q Vec) Vec {
	return Vec{
		X: p.X + q.X,
		Y: p.Y + q.Y,
	}
}

func (p Vec) AddWrap(q Vec, maxX, maxY int) Vec {
	return Vec{
		X: (p.X + q.X + maxX) % maxX,
		Y: (p.Y + q.Y + maxY) % maxY,
	}
}

func (p Vec) Sub(q Vec) Vec {
	return Vec{
		X: p.X - q.X,
		Y: p.Y - q.Y,
	}
}

func (p Vec) MulScalar(s int) Vec {
	return Vec{
		X: p.X * s,
		Y: p.Y * s,
	}
}

func (p Vec) Unit() Vec {
	var newX, newY int

	if p.X != 0 {
		newX = p.X / abs(p.X)
	}
	if p.Y != 0 {
		newY = p.Y / abs(p.Y)
	}
	return Vec{newX, newY}
}

func (p Vec) Magnitude() int {
	// restricted definition for this problem
	return abs(p.X) + abs(p.Y)
}

type Tile byte

const (
	TileNil Tile = iota
	TileFloor
	TileWall
)

type Dir byte

const (
	DirNil Dir = iota
	DirN
	DirE
	DirS
	DirW
)

func (d Dir) Rotate(r Rotate) Dir {
	switch r {
	case RotateNil:
		return d
	case RotateCW:
		return d%4 + 1
	case RotateCCW:
		return (d+2)%4 + 1
	case Rotate180:
		return (d+1)%4 + 1
	}
	panic(fmt.Sprintf("invalid rotation: %v", byte(r)))
}

func (d Dir) String() string {
	switch d {
	case DirN:
		return "N"
	case DirE:
		return "E"
	case DirS:
		return "S"
	case DirW:
		return "W"
	}
	panic(fmt.Sprintf("invalid direction: %v", byte(d)))
}

func (d Dir) Vec() Vec {
	switch d {
	case DirN:
		return UnitN
	case DirE:
		return UnitE
	case DirS:
		return UnitS
	case DirW:
		return UnitW
	}
	panic(fmt.Sprintf("invalid direction: %v", byte(d)))
}

func (d Dir) PasswordValue() int {
	switch d {
	case DirN:
		return 3
	case DirE:
		return 0
	case DirS:
		return 1
	case DirW:
		return 2
	}
	panic(fmt.Sprintf("invalid direction: %v", byte(d)))
}

type Rotate byte

const (
	RotateNil Rotate = iota
	RotateCW
	RotateCCW
	Rotate180
)

func (r Rotate) Invert() Rotate {
	switch r {
	case RotateCW:
		return RotateCCW
	case RotateCCW:
		return RotateCW
	case Rotate180:
		return Rotate180
	case RotateNil:
		return RotateNil
	}
	panic(fmt.Sprintf("invalid rotation: %v", byte(r)))
}

type Grid struct {
	MaxX, MaxY int
	Tiles      []Tile

	EdgePairs []EdgePair
}

func NewGrid(maxX, maxY int) *Grid {
	return &Grid{
		MaxX:  maxX,
		MaxY:  maxY,
		Tiles: make([]Tile, maxX*maxY),
	}
}

func (g *Grid) Get(p Vec) Tile {
	return g.Tiles[p.X+p.Y*g.MaxX]
}

func (g *Grid) Set(p Vec, t Tile) {
	g.Tiles[p.X+p.Y*g.MaxX] = t
}

func (g *Grid) Next(p Vec, d Dir) (Vec, Tile) {
	dirVec := d.Vec()

	cur := p.AddWrap(dirVec, g.MaxX, g.MaxY)
	for g.Get(cur) == TileNil {
		cur = cur.AddWrap(dirVec, g.MaxX, g.MaxY)
	}

	return cur, g.Get(cur)
}

func (g *Grid) InBounds(p Vec) bool {
	return 0 <= p.X && p.X < g.MaxX &&
		0 <= p.Y && p.Y < g.MaxY
}

func (g *Grid) NextCube(p Vec, d Dir) (Vec, Dir, Tile) {
	dirVec := d.Vec()

	destP := p.Add(dirVec)

	if g.InBounds(destP) && g.Get(destP) != TileNil {
		return destP, d, g.Get(destP)
	}

	ep := g.FindEdgePair(destP, d)

	warpedP, warpedDir := ep.Warp(destP, d)
	return warpedP, warpedDir, g.Get(warpedP)
}

func (g *Grid) FindEdgePair(p Vec, d Dir) EdgePair {
	orientation := OrientationNS
	if d != DirN && d != DirS {
		orientation = OrientationEW
	}

	foundEdgePairs := []EdgePair{}
	for _, ep := range g.EdgePairs {
		if ep.E1.Orientation() == orientation && ep.E1.Includes(p) {
			foundEdgePairs = append(foundEdgePairs, ep)
		}
		if ep.E2.Orientation() == orientation && ep.E2.Includes(p) {
			foundEdgePairs = append(foundEdgePairs, ep)
		}
	}

	if len(foundEdgePairs) != 1 {
		panic(fmt.Sprintf("expected 1 edge pair, found %v", len(foundEdgePairs)))
	} else {
		return foundEdgePairs[0]
	}

	panic("no edge pair found")
}

func (g *Grid) FindStart() Vec {
	for x := 0; x < g.MaxY; x++ {
		if g.Get(Vec{x, 0}) == TileFloor {
			return Vec{x, 0}
		}
	}
	panic("no start found")
}

func (g *Grid) LoadAndValidateEdgePairs(edgePairs []EdgePair) error {
	for _, ep := range edgePairs {
		if err := g.validateEdge(ep.E1); err != nil {
			return err
		}
		if err := g.validateEdge(ep.E2); err != nil {
			return err
		}
	}
	g.EdgePairs = edgePairs
	return nil
}

func (g *Grid) validateEdge(e Edge) error {
	diff := e.End.Sub(e.Start)

	if diff.X != 0 && diff.Y != 0 {
		return fmt.Errorf("edge is not straight: %v", e)
	}
	if diff.Magnitude() != 49 {
		return fmt.Errorf("edge is not 49 units long: %d, %v, %d", diff.Magnitude(), e)
	}

	for _, pt := range e.Points() {
		if !g.InBounds(pt) {
			continue
		}
		if g.Get(pt) != TileNil {
			return fmt.Errorf("edge pair has point within bounds: %v", pt)
		}
	}
	return nil
}

type Edge struct {
	Start, End Vec
}

func (e Edge) Includes(p Vec) bool {
	p1 := e.Start
	p2 := e.End
	if p1.X > p2.X || p1.Y > p2.Y {
		p1, p2 = p2, p1
	}

	return p1.X <= p.X && p.X <= p2.X &&
		p1.Y <= p.Y && p.Y <= p2.Y
}

func (e Edge) Points() []Vec {
	var points []Vec
	for x := e.Start.X; x <= e.End.X; x++ {
		for y := e.Start.Y; y <= e.End.Y; y++ {
			points = append(points, Vec{x, y})
		}
	}
	return points
}

type Orientation byte

const (
	OrientationNil Orientation = iota
	OrientationNS
	OrientationEW
)

func (o Orientation) String() string {
	switch o {
	case OrientationNS:
		return "NS"
	case OrientationEW:
		return "EW"
	}
	panic(fmt.Sprintf("invalid orientation: %v", byte(o)))
}

func (e Edge) Orientation() Orientation {
	if e.Start.X == e.End.X {
		return OrientationEW
	}
	return OrientationNS
}

type EdgePair struct {
	Name   string
	E1, E2 Edge
	E1ToE2 Rotate
}

func NewEdgePair(name string, e1s, e1e, e2s, e2e Vec, r Rotate) EdgePair {
	return EdgePair{
		Name: name,
		E1: Edge{
			Start: e1s,
			End:   e1e,
		},
		E2: Edge{
			Start: e2s,
			End:   e2e,
		},
		E1ToE2: r,
	}
}

func (ep *EdgePair) Warp(p Vec, d Dir) (Vec, Dir) {
	wantOrientation := OrientationNS
	if d != DirN && d != DirS {
		wantOrientation = OrientationEW
	}

	var from, to Edge
	var rotation Rotate

	switch {
	case ep.E1.Orientation() == wantOrientation && ep.E1.Includes(p):
		from, to = ep.E1, ep.E2
		rotation = ep.E1ToE2
	case ep.E2.Orientation() == wantOrientation && ep.E2.Includes(p):
		from, to = ep.E2, ep.E1
		rotation = ep.E1ToE2.Invert()
	default:
		panic("point not in edge pair")
	}

	mag := p.Sub(from.Start).Magnitude()

	vec := to.End.Sub(to.Start).Unit().MulScalar(mag)

	destPoint := to.Start.Add(vec)
	destDir := d.Rotate(rotation)

	return destPoint.Add(destDir.Vec()), destDir
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
