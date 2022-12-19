package game_test

import (
	"testing"

	"aoc-2022-19/game"
)

func TestResourceVec_Add(t *testing.T) {
	var (
		v1   = game.ResourceVec{game.Ore: 1, game.Obsidian: 1}
		v2   = game.ResourceVec{game.Geode: 1, game.Clay: 1}
		want = game.ResourceVec{game.Ore: 1, game.Obsidian: 1, game.Geode: 1, game.Clay: 1}
	)

	if result := v1.Add(v2); !result.Eq(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestResourceVec_Sub(t *testing.T) {
	var (
		v1   = game.ResourceVec{game.Ore: 1, game.Obsidian: 1}
		v2   = game.ResourceVec{game.Geode: 1, game.Clay: 1}
		want = game.ResourceVec{game.Ore: 1, game.Obsidian: 1, game.Geode: -1, game.Clay: -1}
	)

	if result := v1.Sub(v2); !result.Eq(want) {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestResourceVec_NonNeg(t *testing.T) {
	testCases := []struct {
		vec  game.ResourceVec
		want bool
	}{
		{game.ResourceVec{game.Ore: 1, game.Obsidian: 1}, true},
		{game.ResourceVec{game.Ore: 1, game.Obsidian: -1}, false},
		{game.ResourceVec{}, true},
	}

	for _, tt := range testCases {
		if result := tt.vec.Pos(); result != tt.want {
			t.Errorf("got %v, want %v", result, tt.want)
		}
	}
}

func TestTurnsToBuild(t *testing.T) {
	testCases := []struct {
		s            game.State
		robotToBuild game.ResourceType
		want         int
	}{
		{game.State{
			Blueprint: game.Blueprint{
				RobotCost: map[game.ResourceType]game.ResourceVec{
					game.Ore: {game.Ore: 1},
				},
			},
			ResourceCount: game.ResourceVecZero,
			RobotCount:    game.ResourceVec{game.Ore: 1},
		}, game.Ore, 1},
		{game.State{
			Blueprint: game.Blueprint{
				RobotCost: map[game.ResourceType]game.ResourceVec{
					game.Ore: {game.Ore: 1, game.Obsidian: 2},
				},
			},
			ResourceCount: game.ResourceVecZero,
			RobotCount:    game.ResourceVec{game.Ore: 1},
		}, game.Ore, -1},
		{game.State{
			Blueprint: game.Blueprint{
				RobotCost: map[game.ResourceType]game.ResourceVec{
					game.Ore: {game.Ore: 1, game.Obsidian: 8},
				},
			},
			ResourceCount: game.ResourceVec{game.Ore: 100, game.Obsidian: 3},
			RobotCount:    game.ResourceVec{game.Ore: 1, game.Obsidian: 1},
		}, game.Ore, 5},
	}

	for _, tt := range testCases {
		if result := tt.s.TurnsToBuild(tt.robotToBuild); result != tt.want {
			t.Errorf("got %v, want %v", result, tt.want)
		}
	}
}

func TestState_BuildRobot(t *testing.T) {
	s := game.State{
		Blueprint: game.Blueprint{
			RobotCost: map[game.ResourceType]game.ResourceVec{
				game.Ore: {game.Ore: 1, game.Obsidian: 8},
			},
		},
		ResourceCount: game.ResourceVec{game.Ore: 5, game.Obsidian: 12},
		RobotCount:    game.ResourceVec{},
	}

	s.BuildRobot(game.Ore)

	if !s.ResourceCount.Eq(game.ResourceVec{game.Ore: 4, game.Obsidian: 4}) {
		t.Errorf("got %v, want %v", s.ResourceCount, game.ResourceVec{game.Ore: 4, game.Obsidian: 4})
	}
	if !s.RobotCount.Eq(game.ResourceVec{game.Ore: 1}) {
		t.Errorf("got %v, want %v", s.RobotCount, game.ResourceVec{game.Ore: 1})
	}
}

func TestState(t *testing.T) {
	s := game.State{
		Blueprint: game.Blueprint{
			RobotCost: map[game.ResourceType]game.ResourceVec{
				game.Ore:      {game.Ore: 4},
				game.Clay:     {game.Ore: 2},
				game.Obsidian: {game.Ore: 3, game.Clay: 14},
				game.Geode:    {game.Ore: 2, game.Obsidian: 7},
			},
		},
		ResourceCount: game.ResourceVec{game.Ore: 1},
		RobotCount:    game.ResourceVec{game.Ore: 1},
	}

	if s.TurnsToBuild(game.Clay) != 2 {
		t.Errorf("got %v, want %v", s.TurnsToBuild(game.Clay), 2)
	}
	s.Advance(2) // minute 3
	s.BuildRobot(game.Clay)
	checkState(t, s,
		game.ResourceVec{game.Ore: 1},
		game.ResourceVec{game.Ore: 1, game.Clay: 1},
	)

	s.Advance(2) // minute 5
	s.BuildRobot(game.Clay)
	checkState(t, s,
		game.ResourceVec{game.Ore: 1, game.Clay: 2},
		game.ResourceVec{game.Ore: 1, game.Clay: 2},
	)

	s.Advance(2) // minute 7
	s.BuildRobot(game.Clay)
	checkState(t, s,
		game.ResourceVec{game.Ore: 1, game.Clay: 6},
		game.ResourceVec{game.Ore: 1, game.Clay: 3},
	)

	s.Advance(4) // minute 11
	s.BuildRobot(game.Obsidian)
	checkState(t, s,
		game.ResourceVec{game.Ore: 2, game.Clay: 4},
		game.ResourceVec{game.Ore: 1, game.Clay: 3, game.Obsidian: 1},
	)

	s.Advance(1) // minute 12
	s.BuildRobot(game.Clay)
	checkState(t, s,
		game.ResourceVec{game.Ore: 1, game.Clay: 7, game.Obsidian: 1},
		game.ResourceVec{game.Ore: 1, game.Clay: 4, game.Obsidian: 1},
	)

}

func checkState(t *testing.T, s game.State, wantResourceCount, wantRobotCount game.ResourceVec) {
	if !s.ResourceCount.Eq(wantResourceCount) {
		t.Errorf("got %v, want %v", s.ResourceCount, wantResourceCount)
	}
	if !s.RobotCount.Eq(wantRobotCount) {
		t.Errorf("got %v, want %v", s.RobotCount, wantRobotCount)
	}
}
