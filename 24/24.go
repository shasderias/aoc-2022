package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"sort"

	"aoc-2022-24/grid"
	"aoc-2022-24/que"
)

func main() {
	star12("input.txt")
}

type blizzardState struct {
	m          map[grid.Vec]blizzardStateElem
	hashString string
}

func (bs blizzardState) Occupied(pos grid.Vec) bool {
	_, ok := bs.m[pos]
	return ok
}

func (bs blizzardState) HashString() string {
	return bs.hashString
}

type blizzardStateElem struct {
	t     grid.Tile
	count int
}

type BlizzardStateManager struct {
	maxX, maxY   int
	initialState []grid.Blizzard
	stateAt      map[int]blizzardState
}

func NewBlizzardStateManager(maxX, maxY int, initialState []grid.Blizzard) *BlizzardStateManager {
	return &BlizzardStateManager{
		maxX:         maxX,
		maxY:         maxY,
		initialState: initialState,
		stateAt:      map[int]blizzardState{},
	}
}

func (bsm *BlizzardStateManager) At(minute int) blizzardState {
	if s, ok := bsm.stateAt[minute]; ok {
		return s
	}

	if bsm.stateAt == nil {
		bsm.stateAt = make(map[int]blizzardState)
	}

	bsMap := make(map[grid.Vec]blizzardStateElem)
	for _, b := range bsm.initialState {
		v := b.At(bsm.maxX, bsm.maxY, minute)
		bsMap[v] = blizzardStateElem{
			t:     b.Type,
			count: bsMap[v].count + 1,
		}

	}

	bsSlice := make([]grid.Vec, 0)
	for k := range bsMap {
		bsSlice = append(bsSlice, k)
	}

	sort.Slice(bsSlice, func(i, j int) bool {
		return bsSlice[i].X < bsSlice[j].X || bsSlice[i].Y < bsSlice[j].Y
	})

	hashString := fmt.Sprintf("%v", bsSlice)

	bsm.stateAt[minute] = blizzardState{
		m:          bsMap,
		hashString: hashString,
	}

	return bsm.stateAt[minute]
}

func star12(ip string) {
	maze := readInput(ip)
	fmt.Println(maze)
	start, exit := findOpenings(maze)
	fmt.Println(start, exit, start.MDist(exit))
	blizzards := findBlizzards(maze)

	bsm := NewBlizzardStateManager(maze.MaxX, maze.MaxY, blizzards)

	startToExit := solve(maze, bsm, start, exit, 0)
	fmt.Println(startToExit)
	exitToStart := solve(maze, bsm, exit, start, startToExit)
	fmt.Println(exitToStart)
	startToExit2 := solve(maze, bsm, start, exit, exitToStart)
	fmt.Println(startToExit2)

}

func solve(maze *grid.Dense, bsm *BlizzardStateManager, start, goal grid.Vec, startMinute int) int {
	pq := que.NewPriority[State]()

	pq.Push(State{
		Pos:    start,
		Goal:   goal,
		Minute: startMinute,
	})

	minSteps := math.MaxInt

	seenStates := make(map[string]struct{})

	pushIfViable := func(s State) {
		if s.MDistToGoal()+s.Minute > minSteps {
			return
		}

		stateHash := s.SeenHash(bsm.At(s.Minute))
		if _, ok := seenStates[stateHash]; ok {
			return
		}

		pq.Push(s)
		seenStates[stateHash] = struct{}{}
	}

	for pq.Len() > 0 {
		var (
			eval = pq.Pop()
			bs   = bsm.At(eval.Minute + 1)
		)

		for _, c := range grid.Cardinal {
			nextPos := eval.Pos.Add(c)

			switch {
			case nextPos == goal:
				minSteps = eval.Minute + 1
				fmt.Println("found solution in", minSteps, "minutes")
				fmt.Println(eval.Route)
				goto foundExit
			case nextPos == start:
				continue
			case !maze.InInBounds(nextPos):
				continue
			case bs.Occupied(nextPos):
				continue
			}

			pushIfViable(State{
				Pos:     nextPos,
				Goal:    goal,
				Minute:  eval.Minute + 1,
				Penalty: eval.Penalty,
				Route:   eval.Route + c.CardinalString(),
			})
		}
		if !bs.Occupied(eval.Pos) {
			pushIfViable(State{
				Pos:     eval.Pos,
				Goal:    goal,
				Minute:  eval.Minute + 1,
				Penalty: eval.Penalty + 10,
				Route:   eval.Route + "Z",
			})
		}

	foundExit:
	}

	return minSteps
}

type State struct {
	Pos     grid.Vec
	Goal    grid.Vec
	Minute  int
	Penalty int
	Route   string
}

func (s State) MDistToGoal() int {
	return s.Pos.MDist(s.Goal)
}

func (s State) Priority() int {
	return s.MDistToGoal()
}

func (s State) SeenHash(bs blizzardState) string {
	return bs.HashString() + fmt.Sprintf("%v", s.Pos)
}

func findOpenings(maze *grid.Dense) (start, exit grid.Vec) {
	var startFound, exitFound bool
	for x := 0; x < maze.MaxX; x++ {
		if maze.Get(grid.Vec{x, 0}) == grid.TileNil {
			start = grid.Vec{x, 0}
			startFound = true
		}
		if maze.Get(grid.Vec{x, maze.MaxY - 1}) == grid.TileNil {
			exit = grid.Vec{x, maze.MaxY - 1}
			exitFound = true
		}
	}
	if !startFound {
		panic("no start found")
	}
	if !exitFound {
		panic("no exit found")
	}
	return start, exit
}

func findBlizzards(maze *grid.Dense) []grid.Blizzard {
	blizzards := make([]grid.Blizzard, 0)
	maze.Walk(func(v grid.Vec, t grid.Tile) {
		switch t {
		case grid.TileBlizzardE, grid.TileBlizzardW, grid.TileBlizzardN, grid.TileBlizzardS:
			blizzards = append(blizzards, grid.Blizzard{Type: t, Pos: v})
		}
	})
	return blizzards
}

func readInput(ip string) *grid.Dense {
	input, err := os.ReadFile(ip)
	if err != nil {
		panic(err)
	}

	lines := bytes.Split(input, []byte{'\n'})

	yMax := len(lines) - 1
	xMax := len(lines[0])

	maze := grid.New(xMax, yMax)

	scanner := bufio.NewScanner(bytes.NewBuffer(input))

	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()
		if line == "" {
			break
		}
		for i, t := range line {
			switch t {
			case '.':
				maze.Set(grid.Vec{i, y}, grid.TileNil)
			case '#':
				maze.Set(grid.Vec{i, y}, grid.TileWall)
			case '>':
				maze.Set(grid.Vec{i, y}, grid.TileBlizzardE)
			case '<':
				maze.Set(grid.Vec{i, y}, grid.TileBlizzardW)
			case '^':
				maze.Set(grid.Vec{i, y}, grid.TileBlizzardN)
			case 'v':
				maze.Set(grid.Vec{i, y}, grid.TileBlizzardS)
			default:
				panic(fmt.Sprintf("unknown token %q", t))
			}
		}
	}

	return maze
}
