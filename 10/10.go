package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const inputFilePath = "input.txt"

func main() {
	star12()
}

type Opcode int

const (
	OpcodeNoop Opcode = iota
	OpcodeAddX
)

type Instruction struct {
	opcode  Opcode
	operand int
}

func star12() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	instructions := []Instruction{}

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case line == "noop":
			instructions = append(instructions, Instruction{OpcodeNoop, 0})
		case strings.HasPrefix(line, "addx"):
			operand, err := strconv.Atoi(strings.Split(line, " ")[1])
			if err != nil {
				return err
			}
			instructions = append(instructions, Instruction{OpcodeNoop, 0})
			instructions = append(instructions, Instruction{OpcodeAddX, operand})
		}
	}

	fmt.Println(run(instructions, 20))
	draw(instructions, 40)

	return nil
}

func run(instructions []Instruction, nextSignalCycle int) int {
	var (
		x   = 1
		acc = 0
	)

	for i, inst := range instructions {
		cycle := i + 1
		if cycle == nextSignalCycle {
			acc += cycle * x
			nextSignalCycle += 40
		}

		switch inst.opcode {
		case OpcodeNoop:
		case OpcodeAddX:
			x += inst.operand
		}
	}

	return acc
}

func draw(instructions []Instruction, screenWidth int) {
	var (
		x = 1
	)

	for i, inst := range instructions {
		var (
			cycle = i + 1
			pos   = i % screenWidth
		)

		if pos == x || pos == x+1 || pos == x-1 {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}

		if cycle%screenWidth == 0 {
			fmt.Println()
		}

		switch inst.opcode {
		case OpcodeNoop:
		case OpcodeAddX:
			x += inst.operand
		}
	}
}
