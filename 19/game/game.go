package game

import "fmt"

type ResourceType string

const (
	ResourceTypeNil ResourceType = ""
	Ore                          = "or"
	Clay                         = "cl"
	Obsidian                     = "ob"
	Geode                        = "ge"
)

var ResourceTypes = []ResourceType{Ore, Clay, Obsidian, Geode}

type ResourceVec map[ResourceType]int

var ResourceVecZero = ResourceVec{}

func (a ResourceVec) Add(b ResourceVec) ResourceVec {
	result := ResourceVec{}
	for _, rt := range ResourceTypes {
		result[rt] = a[rt] + b[rt]
	}
	return result
}

func (a ResourceVec) Sub(b ResourceVec) ResourceVec {
	result := ResourceVec{}
	for _, rt := range ResourceTypes {
		result[rt] = a[rt] - b[rt]
	}
	return result
}

func (a ResourceVec) Mul(b int) ResourceVec {
	result := ResourceVec{}
	for _, rt := range ResourceTypes {
		result[rt] = a[rt] * b
	}
	return result
}

func (a ResourceVec) Eq(b ResourceVec) bool {
	for _, rt := range ResourceTypes {
		if a[rt] != b[rt] {
			return false
		}
	}
	return true
}

func (v ResourceVec) Unit() ResourceVec {
	result := ResourceVec{}
	for _, rt := range ResourceTypes {
		switch {
		case v[rt] > 0:
			result[rt] = 1
		case v[rt] < 0:
			result[rt] = -1
		default:
			result[rt] = 0
		}
	}
	return result
}

func (a ResourceVec) CeilDiv(b ResourceVec) ResourceVec {
	result := ResourceVec{}
	for _, rt := range ResourceTypes {
		if a[rt] == 0 && b[rt] == 0 {
			result[rt] = 0
		} else {
			result[rt] = ceilDiv(a[rt], b[rt])
		}
	}
	return result
}

func (v ResourceVec) ClampZero() ResourceVec {
	result := ResourceVec{}
	for _, rt := range ResourceTypes {
		if v[rt] < 0 {
			result[rt] = 0
		} else {
			result[rt] = v[rt]
		}
	}
	return result
}

func (v ResourceVec) Max() int {
	max := 0
	for _, rt := range ResourceTypes {
		if v[rt] > max {
			max = v[rt]
		}
	}
	return max
}

func (v ResourceVec) Pos() bool {
	for _, rt := range ResourceTypes {
		if v[rt] < 0 {
			return false
		}
	}
	return true
}

func (v ResourceVec) Clone() ResourceVec {
	result := ResourceVec{}
	for _, rt := range ResourceTypes {
		result[rt] = v[rt]
	}
	return result
}

type Blueprint struct {
	Number    int
	RobotCost map[ResourceType]ResourceVec
}

type State struct {
	Blueprint

	ResourceCount ResourceVec
	RobotCount    ResourceVec

	Actions []Action
}

func (s *State) Advance(n int) {
	s.ResourceCount = s.ResourceCount.Add(s.RobotCount.Mul(n))
}

func (s *State) BuildRobot(robot ResourceType) {
	s.ResourceCount = s.ResourceCount.Sub(s.RobotCost[robot])
	s.RobotCount[robot]++
}

func (s *State) TurnsToBuild(robot ResourceType) int {
	robotCost := s.Blueprint.RobotCost[robot]
	robotCount := s.RobotCount

	if !robotCount.Unit().Sub(robotCost.Unit()).Pos() {
		return -1
	}

	turns := robotCost.Sub(s.ResourceCount).ClampZero().CeilDiv(robotCount).Max() + 1
	return turns
}

func (s *State) Clone() State {
	cActions := make([]Action, len(s.Actions))
	copy(cActions, s.Actions)
	return State{
		Blueprint:     s.Blueprint,
		ResourceCount: s.ResourceCount.Clone(),
		RobotCount:    s.RobotCount.Clone(),
		Actions:       cActions,
	}
}

func (s *State) Hash(minutes int) string {
	return fmt.Sprintf("%v %v %v %v", minutes, s.Number, s.ResourceCount, s.RobotCount)
}

func ceilDiv(a, b int) int {
	return (a + b - 1) / b
}

type Action struct {
	Minutes int
	Robot   ResourceType
}
