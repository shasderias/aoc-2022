package main

import (
	"aoc-2022-15/grid"
	"bufio"
	"errors"
	"fmt"
	"os"
)

func main() {
	star1("input-example.txt", 10)
	star1("input.txt", 2000000)

	star2("input-example.txt", 0, 20, 0, 20)
	star2("input.txt", 0, 4000000, 0, 4000000)
}

type SensorBeacon struct {
	Sensor grid.Point
	Beacon grid.Point
}

func (sb SensorBeacon) Radius() int {
	return sb.Sensor.Distance(sb.Beacon)
}

func star1(inputFilePath string, targetY int) error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	sensorBeacons := []SensorBeacon{}

	for scanner.Scan() {
		line := scanner.Text()
		var scannerX, scannerY, beaconX, beaconY int
		_, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &scannerX, &scannerY, &beaconX, &beaconY)
		if err != nil {
			return err
		}

		sensorBeacons = append(sensorBeacons, SensorBeacon{Sensor: grid.Point{X: scannerX, Y: scannerY}, Beacon: grid.Point{X: beaconX, Y: beaconY}})
	}

	pointSet := grid.PointSet{}

	for _, sb := range sensorBeacons {
		if abs(targetY-sb.Sensor.Y) > sb.Radius() {
			continue
		}

		remainingX := sb.Radius() - abs(targetY-sb.Sensor.Y)

		minX, maxX := sb.Sensor.X-remainingX, sb.Sensor.X+remainingX

		points := grid.NewPointSetFromRange(grid.Point{X: minX, Y: targetY}, grid.Point{X: maxX, Y: targetY})

		pointSet.Merge(points)
	}

	for _, sb := range sensorBeacons {
		if pointSet.Has(sb.Beacon) {
			pointSet.Remove(sb.Beacon)
		}
		if pointSet.Has(sb.Sensor) {
			pointSet.Remove(sb.Sensor)
		}
	}

	fmt.Println(len(pointSet))

	return nil
}

func star2(inputFilePath string, minX, maxX, minY, maxY int) error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	sensorBeacons := []SensorBeacon{}

	for scanner.Scan() {
		line := scanner.Text()
		var scannerX, scannerY, beaconX, beaconY int
		_, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &scannerX, &scannerY, &beaconX, &beaconY)
		if err != nil {
			return err
		}

		sensorBeacons = append(sensorBeacons, SensorBeacon{Sensor: grid.Point{X: scannerX, Y: scannerY}, Beacon: grid.Point{X: beaconX, Y: beaconY}})
	}

	for y := minY; y <= maxY; y++ {
		distress, err := findDistressInRow(sensorBeacons, minX, maxX, y)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				continue
			} else {
				return err
			}
		}
		fmt.Println(distress.X*4000000 + distress.Y)
		goto found

	}
	fmt.Println("not found")
	return nil
found:
	return nil
}

var ErrNotFound = errors.New("not found")

func findDistressInRow(sensorBeacons []SensorBeacon, minX, maxX int, targetY int) (grid.Point, error) {
	rs := grid.RangeSet{}

	for _, sb := range sensorBeacons {
		if abs(targetY-sb.Sensor.Y) > sb.Radius() {
			continue
		}

		remainingX := sb.Radius() - abs(targetY-sb.Sensor.Y)

		minX, maxX := sb.Sensor.X-remainingX, sb.Sensor.X+remainingX

		rs.Add(grid.Range{Min: minX, Max: maxX})
	}

	switch {
	case len(rs) == 1:
		switch {
		case rs[0].Min == minX && rs[0].Max == maxX-1:
			return grid.Point{X: rs[0].Max, Y: targetY}, nil
		case rs[0].Min == minX+1 && rs[0].Max == maxX:
			return grid.Point{X: rs[0].Min, Y: targetY}, nil
		default:
			return grid.Point{}, ErrNotFound
		}
	case len(rs) == 2:
		switch {
		case rs[0].Width()+rs[1].Width() < (maxX-minX)-1:
			return grid.Point{}, ErrNotFound
		case rs[0].Min > minX || rs[1].Max < maxX:
			return grid.Point{}, ErrNotFound
		default:
			return grid.Point{X: rs[1].Min - 1, Y: targetY}, nil
		}
	default:
		return grid.Point{}, ErrNotFound
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
