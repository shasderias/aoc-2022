package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

var OrdinalDirections3D = []Point3D{
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
	{0, 0, 1},
	{0, 0, -1},
}

func main() {
	if err := star12("input.txt"); err != nil {
		panic(err)
	}
}

func star12(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	cubeCoords := []Point3D{}

	for scanner.Scan() {
		line := scanner.Text()

		var x, y, z int
		if _, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z); err != nil {
			return err
		}

		cubeCoords = append(cubeCoords, Point3D{x, y, z})
	}

	fmt.Println(len(calcExposedFace(cubeCoords)))

	var minX, minY, minZ = math.MaxInt, math.MaxInt, math.MaxInt
	var maxX, maxY, maxZ = 0, 0, 0

	for _, cube := range cubeCoords {
		minX = min(minX, cube.X)
		minY = min(minY, cube.Y)
		minZ = min(minZ, cube.Z)
		maxX = max(maxX, cube.X)
		maxY = max(maxY, cube.Y)
		maxZ = max(maxZ, cube.Z)
	}

	//fmt.Println(minX, minY, minZ, maxX, maxY, maxZ)

	minX--
	minY--
	minZ--
	maxX++
	maxY++
	maxZ++

	inBounds := func(p Point3D) bool {
		return p.X >= minX && p.X <= maxX && p.Y >= minY && p.Y <= maxY && p.Z >= minZ && p.Z <= maxZ
	}

	visitedSet := Point3DSet{}
	visitQueue := []Point3D{{minX, minY, minZ}}
	visitNext := Point3DSet{}

	occupiedPoints := Point3DSet{}
	for _, cube := range cubeCoords {
		occupiedPoints.Add(cube)
	}

	for len(visitQueue) > 0 {
		//fmt.Println("visiting", visitQueue)
		for _, p := range visitQueue {
			visitedSet.Add(p)
			for _, d := range OrdinalDirections3D {
				np := p.Add(d)
				if inBounds(np) && !visitedSet.Has(np) && !occupiedPoints.Has(np) {
					visitNext.Add(np)
				}
			}
		}
		visitQueue = visitNext.Slice()
		visitNext = Point3DSet{}

	}

	visitedSlice := []Point3D{}
	for p := range visitedSet {
		visitedSlice = append(visitedSlice, p)
	}

	externalFaces := calcExposedFace(visitedSlice)
	externalFaceCount := len(externalFaces)

	width := maxX - minX + 1
	height := maxY - minY + 1
	depth := maxZ - minZ + 1

	externalFaceCount -= ((width * height) + (width * depth) + (height * depth)) * 2

	fmt.Println(externalFaceCount)

	return nil
}

func calcExposedFace(unitCubeCoords []Point3D) FaceSet {
	exposedFaces := FaceSet{}
	coveredFaces := FaceSet{}

	for _, cube := range unitCubeCoords {
		faces := cubeFaces(cube)
		for _, face := range faces {
			switch {
			case coveredFaces.Has(face):
				continue
			case exposedFaces.Has(face):
				exposedFaces.Del(face)
				coveredFaces.Add(face)
			default:
				exposedFaces.Add(face)
			}
		}
	}
	return exposedFaces
}

type Point3D struct {
	X, Y, Z int
}

func (a Point3D) Add(b Point3D) Point3D {
	return Point3D{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

var Origin = Point3D{0, 0, 0}

func (a Point3D) Less(b Point3D) int {
	aDist, bDist := a.MDist(Origin), b.MDist(Origin)
	switch {
	case a.Eq(b):
		return 0
	case aDist < bDist:
		return -1
	case aDist > bDist:
		return 1
	default:
		switch {
		case a.X < b.X:
			return -1
		case a.X > b.X:
			return 1
		case a.Y < b.Y:
			return -1
		case a.Y > b.Y:
			return 1
		case a.Z < b.Z:
			return -1
		case a.Z > b.Z:
			return 1
		default:
			panic("unreachable")
		}
	}
}

func (a Point3D) Eq(b Point3D) bool {
	return a.X == b.X && a.Y == b.Y && a.Z == b.Z
}

func (a Point3D) MDist(b Point3D) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y) + abs(a.Z-b.Z)
}

type Face [4]Point3D

func (f Face) EqStr() string {
	var cp [4]Point3D
	copy(cp[:], f[:])
	sort.Slice(cp[:], func(i, j int) bool {
		return cp[i].Less(cp[j]) < 0
	})
	return fmt.Sprintf("%v", cp)
}

func cubeFaces(loc Point3D) []Face {
	return []Face{
		{ // bottom face
			{loc.X, loc.Y, loc.Z},
			{loc.X, loc.Y + 1, loc.Z},
			{loc.X + 1, loc.Y + 1, loc.Z},
			{loc.X + 1, loc.Y, loc.Z},
		},
		{ // front face
			{loc.X, loc.Y, loc.Z},
			{loc.X, loc.Y, loc.Z + 1},
			{loc.X + 1, loc.Y, loc.Z + 1},
			{loc.X + 1, loc.Y, loc.Z},
		},
		{ // left face
			{loc.X, loc.Y, loc.Z},
			{loc.X, loc.Y + 1, loc.Z},
			{loc.X, loc.Y + 1, loc.Z + 1},
			{loc.X, loc.Y, loc.Z + 1},
		},
		{ // top face
			{loc.X, loc.Y, loc.Z + 1},
			{loc.X, loc.Y + 1, loc.Z + 1},
			{loc.X + 1, loc.Y + 1, loc.Z + 1},
			{loc.X + 1, loc.Y, loc.Z + 1},
		},
		{ // right face
			{loc.X + 1, loc.Y, loc.Z},
			{loc.X + 1, loc.Y, loc.Z + 1},
			{loc.X + 1, loc.Y + 1, loc.Z + 1},
			{loc.X + 1, loc.Y + 1, loc.Z},
		},
		{ // back face
			{loc.X, loc.Y + 1, loc.Z},
			{loc.X, loc.Y + 1, loc.Z + 1},
			{loc.X + 1, loc.Y + 1, loc.Z + 1},
			{loc.X + 1, loc.Y + 1, loc.Z},
		},
	}
}

type FaceSet map[string]Face

func (fs FaceSet) Add(f Face) {
	fs[f.EqStr()] = f
}

func (fs FaceSet) Del(f Face) {
	delete(fs, f.EqStr())
}

func (fs FaceSet) Has(f Face) bool {
	_, ok := fs[f.EqStr()]
	return ok
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Point3DSet map[Point3D]struct{}

func (ps Point3DSet) Add(p Point3D) {
	ps[p] = struct{}{}
}

func (ps Point3DSet) Del(p Point3D) {
	delete(ps, p)
}

func (ps Point3DSet) Has(p Point3D) bool {
	_, ok := ps[p]
	return ok
}

func (ps Point3DSet) Slice() []Point3D {
	var out []Point3D
	for p := range ps {
		out = append(out, p)
	}
	return out
}
