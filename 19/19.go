package main

import (
	"bufio"
	"fmt"
	"os"

	"aoc-2022-19/game"
)

func main() {
	if err := star1(); err != nil {
		panic(err)
	}
	cache = make(map[string]int)
	if err := star2(); err != nil {
		panic(err)
	}
}

func searchParamsHeruistic(bp game.Blueprint, timeLimit int) searchParams {
	return searchParams{
		timeLimit: timeLimit,
		overproductionThresholds: game.ResourceVec{
			game.Ore:      12,
			game.Clay:     bp.RobotCost[game.Obsidian][game.Clay] * 4,
			game.Obsidian: 999,
			game.Geode:    999,
		},
		robotsCutoff: game.ResourceVec{
			game.Ore:      timeLimit / 3,
			game.Clay:     timeLimit/2 + 6,
			game.Obsidian: 99,
			game.Geode:    99,
		},
	}
}

func star1() error {
	blueprints, err := parseInput("input.txt")
	if err != nil {
		return err
	}

	acc := 0

	for _, bp := range blueprints {
		largestGeode := findLargest(bp, searchParamsHeruistic(bp, 24))
		quality := largestGeode * bp.Number
		acc += quality
		fmt.Println(bp.Number, quality)
	}

	fmt.Println(acc)

	return nil
}

func star2() error {
	blueprints, err := parseInput("input.txt")
	if err != nil {
		return err
	}

	if len(blueprints) > 3 {
		blueprints = blueprints[:3]
	}

	acc := 1

	for _, bp := range blueprints {
		largestGeode := findLargest(bp, searchParamsHeruistic(bp, 32))
		acc *= largestGeode
		fmt.Println(bp.Number, largestGeode)
	}

	fmt.Println(acc)

	return nil
}

type searchParams struct {
	timeLimit                int
	overproductionThresholds game.ResourceVec
	robotsCutoff             game.ResourceVec
}

func findLargest(bp game.Blueprint, sp searchParams) int {
	return recurse(sp, game.State{
		Blueprint: bp,
		RobotCount: map[game.ResourceType]int{
			game.Ore: 1,
		},
		ResourceCount: map[game.ResourceType]int{
			game.Ore: 1,
		},
	}, 1)
}

func recurse(sp searchParams, s game.State, minutes int) int {
	maxGeodeCount := 0

	for _, rt := range game.ResourceTypes {
		turns := s.TurnsToBuild(rt)
		if turns == -1 || minutes+turns > sp.timeLimit {
			continue
		}
		if s.ResourceCount[rt] > sp.overproductionThresholds[rt] {
			continue
		}
		if minutes > sp.robotsCutoff[rt] {
			continue
		}

		s := s.Clone()
		s.Advance(turns)
		s.BuildRobot(rt)
		s.Actions = append(s.Actions, game.Action{minutes + turns, rt})

		maxGeodeCount = max(maxGeodeCount, memoize(recurse)(sp, s, minutes+turns))
	}

	if maxGeodeCount == 0 {
		// no moves left, advance to end

		s.Advance(sp.timeLimit - minutes)

		return s.ResourceCount[game.Geode]
	}

	return maxGeodeCount
}

var cache = make(map[string]int)

func memoize(f func(searchParams, game.State, int) int) func(searchParams, game.State, int) int {
	return func(sp searchParams, s game.State, minutes int) int {
		if v, ok := cache[s.Hash(minutes)]; ok {
			return v
		}

		v := f(sp, s, minutes)
		cache[s.Hash(minutes)] = v
		return v
	}
}

func parseInput(inputPath string) ([]game.Blueprint, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	blueprints := []game.Blueprint{}

	for scanner.Scan() {
		line := scanner.Text()

		var (
			bpNumber               int
			oreRobotOreCost        int
			clayRobotOreCost       int
			obsidianRobotOreCost   int
			obsidianRobotClayCost  int
			geodeRobotOreCost      int
			geodeRobotObsidianCost int
		)
		_, err := fmt.Sscanf(line,
			"Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
			&bpNumber, &oreRobotOreCost, &clayRobotOreCost, &obsidianRobotOreCost, &obsidianRobotClayCost, &geodeRobotOreCost, &geodeRobotObsidianCost,
		)
		if err != nil {
			return nil, err
		}

		blueprints = append(blueprints, game.Blueprint{
			Number: bpNumber,
			RobotCost: map[game.ResourceType]game.ResourceVec{
				game.Ore:      {game.Ore: oreRobotOreCost},
				game.Clay:     {game.Ore: clayRobotOreCost},
				game.Obsidian: {game.Ore: obsidianRobotOreCost, game.Clay: obsidianRobotClayCost},
				game.Geode:    {game.Ore: geodeRobotOreCost, game.Obsidian: geodeRobotObsidianCost},
			},
		})
	}

	return blueprints, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
