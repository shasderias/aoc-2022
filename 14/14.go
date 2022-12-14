package main

import (
	"aoc-2022-14/grid"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const inputFilePath = "input.txt"

func main() {
	star1()
	star2()
}

type Path []grid.Point

func (p *Path) Draw(g *grid.Sparse[Entity]) {
	for i := 0; i < len(*p)-1; i++ {
		var (
			s      = (*p)[i]
			e      = (*p)[i+1]
			relDir = s.DirOf(e)
			dist   = s.Distance(e) + 1
		)
		for j := 0; j < dist; j++ {
			g.Set(s.Add(relDir.Unit().Mul(j)), EntityRock)
		}
	}
}

type Entity byte

const (
	EntityNil Entity = iota
	EntityRock
	EntitySand
)

var sandSpawnPoint = grid.Point{500, 0}

func star1() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	paths := []Path{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		pathPoints := strings.Split(line, " -> ")
		path := Path{}
		for _, pathPoint := range pathPoints {
			xy := strings.Split(pathPoint, ",")
			x, err := strconv.Atoi(xy[0])
			if err != nil {
				return err
			}
			y, err := strconv.Atoi(xy[1])
			if err != nil {
				return err
			}
			path = append(path, grid.Point{X: x, Y: y})
		}
		paths = append(paths, path)
	}

	world := grid.NewSparse[Entity]()

	for _, path := range paths {
		path.Draw(world)
	}

	killPlane := 0
	world.Walk(func(p grid.Point, e Entity) {
		if e != EntityRock {
			return
		}
		if p.Y > killPlane {
			killPlane = p.Y
		}
	})
	killPlane++

	rested := 0

	func() {
		for {
			pt, err := spawnSandAndCalcRestCoord(world, killPlane, sandSpawnPoint)
			switch {
			case errors.Is(err, ErrSpawnBlocked):
				panic(err)
			case errors.Is(err, ErrSandVoid):
				return
			case err != nil:
				panic(err)
			default:
				rested++
				world.Set(pt, EntitySand)
			}
		}
	}()

	fmt.Println(rested)

	return nil
}

func star2() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	paths := []Path{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		pathPoints := strings.Split(line, " -> ")
		path := Path{}
		for _, pathPoint := range pathPoints {
			xy := strings.Split(pathPoint, ",")
			x, err := strconv.Atoi(xy[0])
			if err != nil {
				return err
			}
			y, err := strconv.Atoi(xy[1])
			if err != nil {
				return err
			}
			path = append(path, grid.Point{X: x, Y: y})
		}
		paths = append(paths, path)
	}

	world := grid.NewSparse[Entity]()

	for _, path := range paths {
		path.Draw(world)
	}

	floorY := 0
	world.Walk(func(p grid.Point, e Entity) {
		if e != EntityRock {
			return
		}
		if p.Y > floorY {
			floorY = p.Y
		}
	})
	floorY += 2

	rested := 0

	func() {
		for {
			pt, err := spawnSandWithFloorAndCalcRestCoord(world, floorY, sandSpawnPoint)
			switch {
			case errors.Is(err, ErrSpawnBlocked):
				return
			case errors.Is(err, ErrSandVoid):
				return
			case err != nil:
				panic(err)
			default:
				rested++
				world.Set(pt, EntitySand)
			}
		}
	}()

	fmt.Println(rested)

	return nil
}

var (
	ErrSpawnBlocked = errors.New("spawn blocked")
	ErrSandVoid     = fmt.Errorf("sand has returned to the void")
)

func spawnSandAndCalcRestCoord(g *grid.Sparse[Entity], killPlane int, spawnPoint grid.Point) (grid.Point, error) {
	if g.Get(spawnPoint) != EntityNil {
		return grid.Point{}, ErrSpawnBlocked
	}

	pos := spawnPoint

	for {
		if pos.Y >= killPlane {
			return grid.Point{}, ErrSandVoid
		}
		switch {
		case g.Get(pos.Add(grid.DirS.Unit())) == EntityNil:
			pos = pos.Add(grid.DirS.Unit())
		case g.Get(pos.Add(grid.DirSW.Unit())) == EntityNil:
			pos = pos.Add(grid.DirSW.Unit())
		case g.Get(pos.Add(grid.DirSE.Unit())) == EntityNil:
			pos = pos.Add(grid.DirSE.Unit())
		default:
			return pos, nil
		}
	}
}

func spawnSandWithFloorAndCalcRestCoord(g *grid.Sparse[Entity], floorY int, spawnPoint grid.Point) (grid.Point, error) {
	if g.Get(spawnPoint) != EntityNil {
		return grid.Point{}, ErrSpawnBlocked
	}

	pos := spawnPoint

	for {
		if pos.Y+1 == floorY {
			return pos, nil
		}
		switch {
		case g.Get(pos.Add(grid.DirS.Unit())) == EntityNil:
			pos = pos.Add(grid.DirS.Unit())
		case g.Get(pos.Add(grid.DirSW.Unit())) == EntityNil:
			pos = pos.Add(grid.DirSW.Unit())
		case g.Get(pos.Add(grid.DirSE.Unit())) == EntityNil:
			pos = pos.Add(grid.DirSE.Unit())
		default:
			return pos, nil
		}
	}
}
