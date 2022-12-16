package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

const inputFilePath = "input.txt"
const StartingRoom = "AA"

func main() {
	if err := star1(); err != nil {
		panic(err)
	}
}

func star1() error {
	n, err := parseInput(inputFilePath)
	if err != nil {
		return err
	}
	n.MaxMinutes = 30

	n.ConstructDistanceMatrix()
	n.SolveStar1()
	n.SolveStar2()

	return nil
}

func parseInput(filename string) (*Network, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lineBytes := bytes.Split(fileBytes, []byte("\n"))

	if len(lineBytes[len(lineBytes)-1]) == 0 {
		lineBytes = lineBytes[:len(lineBytes)-1]
	}

	n := NewNetwork()

	for i, line := range lineBytes {
		scanData := bytes.Split(line, []byte("; "))

		var roomName string
		var flowRate int
		_, err := fmt.Sscanf(string(scanData[0]), "Valve %s has flow rate=%d", &roomName, &flowRate)
		if err != nil {
			return nil, fmt.Errorf("err scanning line %d: %s: %w", i, scanData[0], err)
		}

		tunnelsStr := bytes.TrimPrefix(scanData[1], []byte("tunnels lead to valves "))
		tunnelsStr = bytes.TrimPrefix(tunnelsStr, []byte("tunnel leads to valve "))

		tunnelList := strings.Split(string(tunnelsStr), ", ")

		n.AddRoom(roomName, flowRate, tunnelList)
	}

	for _, room := range n.rooms {
		for _, tunnelStr := range room.TunnelsStr {
			if _, ok := n.rooms[tunnelStr]; !ok {
				return nil, fmt.Errorf("room %s does not exist", tunnelStr)
			}
			room.Connected = append(room.Connected, n.rooms[tunnelStr])
		}
	}

	return n, nil
}

type Room struct {
	Name       string
	FlowRate   int
	TunnelsStr []string
	Connected  []*Room
}

func (r *Room) String() string {
	return fmt.Sprintf("%s(%d)", r.Name, r.FlowRate)
}

type Network struct {
	rooms       map[string]*Room
	sortedRooms []*Room
	sigRooms    []*Room

	MaxMinutes int
	Distances  DistMatrix
}

func NewNetwork() *Network {
	return &Network{
		rooms: make(map[string]*Room),
	}
}

func (n *Network) AddRoom(roomName string, flowRate int, tunnelList []string) *Room {
	if _, ok := n.rooms[roomName]; ok {
		panic(fmt.Errorf("room %s already exists", roomName))
	}
	room := &Room{Name: roomName, FlowRate: flowRate, TunnelsStr: tunnelList}
	n.rooms[roomName] = room
	n.sortedRooms = nil
	n.sigRooms = nil
	return room
}

func (n *Network) Rooms() []*Room {
	if n.sortedRooms != nil {
		return n.sortedRooms
	}
	rooms := make([]*Room, 0, len(n.rooms))
	for _, room := range n.rooms {
		rooms = append(rooms, room)
	}
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Name < rooms[j].Name
	})
	n.sortedRooms = rooms
	return rooms
}

func (n *Network) SignificantRooms() []*Room {
	if n.sigRooms != nil {
		return n.sigRooms
	}
	sigRooms := []*Room{}
	for _, room := range n.Rooms() {
		if room.FlowRate > 0 {
			sigRooms = append(sigRooms, room)
		}
	}
	n.sigRooms = sigRooms
	return sigRooms
}

func (n *Network) ConstructDistanceMatrix() {
	sigRooms := n.SignificantRooms()
	sigRooms = append(sigRooms, n.rooms[StartingRoom])
	m := NewDistMatrix(sigRooms)
	for _, roomA := range sigRooms {
		for _, roomB := range sigRooms {
			if roomA == roomB {
				continue
			}
			dist := findShortestPath(roomA, roomB)
			m.Set(roomA, roomB, dist)
		}
	}
	n.Distances = *m
}

func (n *Network) SolveStar1() {
	startRoom := n.rooms[StartingRoom]

	vom := n.solveRecurseStar1(n.SignificantRooms(), startRoom, 0, valveOpenMap{})

	fmt.Println(vom.TotalFlowRate(n.MaxMinutes))
	fmt.Println(vom)
}

func (n *Network) solveRecurseStar1(significantRooms []*Room, r *Room, minutes int, vom valveOpenMap) valveOpenMap {
	if _, opened := vom[r]; !opened && r.FlowRate > 0 {
		return n.solveRecurseStar1(significantRooms, r, minutes+1, vom.With(r, minutes+1))
	}

	unvisitedRooms := vom.Diff(significantRooms)
	vomSet := []valveOpenMap{vom}
	for _, ur := range unvisitedRooms {
		minutesToRoom := n.Distances.Get(r, ur)

		if minutes+minutesToRoom > n.MaxMinutes {
			continue
		}
		vomCandidate := n.solveRecurseStar1(significantRooms, ur, minutes+minutesToRoom, vom.Clone())
		vomSet = append(vomSet, vomCandidate)
	}

	var bestVOM valveOpenMap
	var bestVOMScore int
	for _, vom := range vomSet {
		score := vom.TotalFlowRate(n.MaxMinutes)
		if score > bestVOMScore {
			bestVOMScore = score
			bestVOM = vom
		}
	}
	return bestVOM
}

func (n *Network) SolveStar2() {
	permutations := [][2][]*Room{}

	sigRooms := n.SignificantRooms()
	for i := 0; i < len(sigRooms)/2; i++ {
		roomSet := Comb(sigRooms, i)
		for _, roomSet := range roomSet {
			roomSet1 := roomSet
			roomSet2 := Diff(sigRooms, roomSet1)
			permutations = append(permutations, [2][]*Room{roomSet1, roomSet2})
		}
	}

	startRoom := n.rooms[StartingRoom]

	var bestVOM valveOpenMap
	var bestVOMScore int

	fmt.Println(len(permutations))

	wg := &sync.WaitGroup{}

	resultsChan := make(chan Result)
	stopCh := make(chan struct{})

	go func() {
		for {
			select {
			case result := <-resultsChan:
				if result.Solution.TotalFlowRate(n.MaxMinutes) > bestVOMScore {
					bestVOMScore = result.Solution.TotalFlowRate(n.MaxMinutes)
					bestVOM = result.Solution
					fmt.Printf("%d/%d\t:%d - %v\n", result.Index, len(permutations), bestVOMScore, bestVOM)
					fmt.Println(result.SolutionPart1)
					fmt.Println(result.SolutionPart2)
				}
			case <-stopCh:
				return
			}
		}
	}()

	for i, perm := range permutations {
		wg.Add(1)
		go n.solveStar2(wg, resultsChan, i, startRoom, perm[0], perm[1])
	}

	wg.Wait()
	close(stopCh)

	fmt.Println(bestVOMScore)
	fmt.Println(bestVOM)
}

type Result struct {
	Index                        int
	SolutionPart1, SolutionPart2 valveOpenMap
	Solution                     valveOpenMap
}

func (n *Network) solveStar2(wg *sync.WaitGroup, resultsCh chan<- Result, i int, startRoom *Room, roomSet1, roomSet2 []*Room) {
	solution1 := n.solveRecurseStar1(roomSet1, startRoom, 4, valveOpenMap{})
	solution2 := n.solveRecurseStar1(roomSet2, startRoom, 4, valveOpenMap{})
	solution := solution1.Merge(solution2)
	resultsCh <- Result{i, solution1, solution2, solution}
	wg.Done()
}

func findShortestPath(start, end *Room) int {
	visited := map[*Room]struct{}{}
	visiting := []*Room{start}
	visitNext := []*Room{}

	for cost := 0; len(visiting) > 0; cost++ {
		for _, room := range visiting {
			if room == end {
				return cost
			}

			for _, nextRoomCandidate := range room.Connected {
				if _, ok := visited[nextRoomCandidate]; !ok {
					visited[nextRoomCandidate] = struct{}{}
					visitNext = append(visitNext, nextRoomCandidate)
				}
			}
		}
		visiting = visitNext
		visitNext = []*Room{}
	}

	return -1
}

func NewDistMatrix(rooms []*Room) *DistMatrix {
	l := len(rooms)
	m := &DistMatrix{
		distances: make([]int, l*l),
		rooms:     rooms,
	}
	return m
}

type DistMatrix struct {
	distances []int
	rooms     []*Room
}

func (m *DistMatrix) roomIndex(r *Room) int {
	for i, room := range m.rooms {
		if room == r {
			return i
		}
	}
	panic(fmt.Errorf("room %s not found", r.Name))
}

func (m *DistMatrix) Set(r1, r2 *Room, dist int) {
	i1 := m.roomIndex(r1)
	i2 := m.roomIndex(r2)
	m.distances[i1*len(m.rooms)+i2] = dist
	m.distances[i2*len(m.rooms)+i1] = dist
}

func (m *DistMatrix) Get(r1, r2 *Room) int {
	i1 := m.roomIndex(r1)
	i2 := m.roomIndex(r2)
	return m.distances[i1*len(m.rooms)+i2]
}

type valveOpenMap map[*Room]int

func (vom valveOpenMap) With(room *Room, minutes int) valveOpenMap {
	clone := make(valveOpenMap, len(vom)+1)
	for k, v := range vom {
		clone[k] = v
	}
	clone[room] = minutes
	return clone
}

func (vom valveOpenMap) Diff(rooms []*Room) []*Room {
	d := []*Room{}
	for _, r := range rooms {
		if _, ok := vom[r]; !ok {
			d = append(d, r)
		}
	}
	return d
}

func (vom valveOpenMap) Merge(other valveOpenMap) valveOpenMap {
	clone := make(valveOpenMap, len(vom)+len(other))
	for k, v := range vom {
		clone[k] = v
	}
	for k, v := range other {
		clone[k] = v
	}
	return clone
}

func (vom valveOpenMap) Clone() valveOpenMap {
	clone := make(valveOpenMap, len(vom))
	for k, v := range vom {
		clone[k] = v
	}
	return clone
}

func (vom valveOpenMap) TotalFlowRate(maxMinutes int) (totalFlowRate int) {
	for room, minutes := range vom {
		totalFlowRate += room.FlowRate * (maxMinutes - minutes)
	}
	return
}

func (vom valveOpenMap) String() string {
	buf := strings.Builder{}
	buf.WriteString("valveOpenMap{")

	sorted := []struct {
		room    *Room
		minutes int
	}{}
	for room, minutes := range vom {
		sorted = append(sorted, struct {
			room    *Room
			minutes int
		}{room, minutes})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].minutes < sorted[j].minutes
	})

	for i, s := range sorted {
		if i < len(sorted)-1 {
			buf.WriteString(fmt.Sprintf("%s:%d, ", s.room.Name, s.minutes))
		} else {
			buf.WriteString(fmt.Sprintf("%s:%d", s.room.Name, s.minutes))
		}
	}

	buf.WriteString("}")

	return buf.String()
}

func Comb[T any](ary []T, m int) [][]T {
	var (
		output  = [][]T{}
		current = make([]T, m)
		last    = m - 1
		n       = len(ary)
	)

	if m > n {
		panic("m > len(ary)")
	}

	if m == 0 {
		return output
	}

	var recurse func(int, int)

	recurse = func(i, next int) {
		for j := next; j < n; j++ {
			current[i] = ary[j]
			if i == last {
				output = append(output, append([]T{}, current...))
			} else {
				recurse(i+1, j+1)
			}
		}
		return
	}

	recurse(0, 0)

	return output
}

func Diff[T comparable](a, b []T) []T {
	m := map[T]struct{}{}
	for _, v := range b {
		m[v] = struct{}{}
	}

	d := []T{}
	for _, v := range a {
		if _, ok := m[v]; !ok {
			d = append(d, v)
		}
	}

	return d
}
