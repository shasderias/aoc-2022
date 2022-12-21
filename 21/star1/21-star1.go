package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := star1(); err != nil {
		panic(err)
	}
}

func star1() error {
	monkeys, err := parseInput("input.txt")
	if err != nil {
		return err
	}

	for _, monkey := range monkeys {
		monkey.ResolveRefs(monkeys)
	}

	fmt.Println(monkeys["root"].Value())

	return nil
}

type Operation string

const (
	OperationNil Operation = ""
	OperationAdd           = "+"
	OperationSub           = "-"
	OperationMul           = "*"
	OperationDiv           = "/"
)

type Monkey struct {
	Name             string
	Ref1Str, Ref2Str string
	Ref1, Ref2       *Monkey
	Constant         int
	Operation        Operation
}

func (m *Monkey) Value() int {
	switch m.Operation {
	case OperationNil:
		return m.Constant
	case OperationAdd:
		return m.Ref1.Value() + m.Ref2.Value()
	case OperationSub:
		return m.Ref1.Value() - m.Ref2.Value()
	case OperationMul:
		return m.Ref1.Value() * m.Ref2.Value()
	case OperationDiv:
		return m.Ref1.Value() / m.Ref2.Value()
	default:
		panic(fmt.Sprintf("unknown operation %s", m.Operation))
	}
}

func (m *Monkey) ResolveRefs(monkeys map[string]*Monkey) {
	var ok bool
	if m.Ref1Str != "" {
		m.Ref1, ok = monkeys[m.Ref1Str]
		if !ok {
			panic(fmt.Sprintf("no monkey %s", m.Ref1Str))
		}
	}
	if m.Ref2Str != "" {
		m.Ref2, ok = monkeys[m.Ref2Str]
		if !ok {
			panic(fmt.Sprintf("no monkey %s", m.Ref2Str))
		}
	}
}

func parseInput(inputPath string) (map[string]*Monkey, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	monkeys := make(map[string]*Monkey)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		parts := strings.Split(line, ": ")

		name := parts[0]

		monkey := Monkey{
			Name: name,
		}

		if strings.ContainsAny(parts[1], "+-*/") {
			_, err := fmt.Sscanf(parts[1], "%s %s %s", &monkey.Ref1Str, &monkey.Operation, &monkey.Ref2Str)
			if err != nil {
				return nil, fmt.Errorf("error parsing %s: '%s'", line, parts[1])
			}
		} else {
			_, err := fmt.Sscanf(parts[1], "%d", &monkey.Constant)
			if err != nil {
				return nil, fmt.Errorf("error parsing %s", line)
			}
		}

		monkeys[name] = &monkey
	}

	return monkeys, nil
}
