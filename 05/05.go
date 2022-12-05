package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const inputFilePath = "input.txt"

type ship struct {
	stacks []stack
}

type instruction struct {
	count, from, to int
}

func newShip(stackCount int) *ship {
	s := ship{
		stacks: make([]stack, stackCount),
	}
	return &s
}

func (s *ship) apply(inst instruction) {
	for i := 0; i < inst.count; i++ {
		c := s.stacks[inst.from-1].popTop()
		s.stacks[inst.to-1].pushTop(c)
	}
}

func (s *ship) apply9001(inst instruction) {
	c := s.stacks[inst.from-1].popNTop(inst.count)
	s.stacks[inst.to-1].pushTop(c...)
}

func (s *ship) topCrates() string {
	var b bytes.Buffer
	for _, stack := range s.stacks {
		b.WriteString(stack.peekTop())
	}
	return b.String()
}

type stack []string

func (s *stack) pushTop(c ...string) {
	*s = append(*s, c...)
}

func (s *stack) pushBot(c ...string) {
	*s = append(c, *s...)
}

func (s *stack) popTop() string {
	c := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return c
}

func (s *stack) popNTop(n int) []string {
	c := (*s)[len(*s)-n:]
	*s = (*s)[:len(*s)-n]
	return c
}

func (s *stack) peekTop() string {
	return (*s)[len(*s)-1]
}

func main() {
	if err := star1(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := star2(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func star1() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	stackLayoutStr, drawingBuf, err := readToStackLayout(scanner)

	stackLayout := strings.Split(strings.TrimSpace(stackLayoutStr), "   ")

	lastStack, err := strconv.Atoi(stackLayout[len(stackLayout)-1])
	if err != nil {
		return err
	}

	s, err := parseDrawing(drawingBuf, lastStack)
	if err != nil {
		return err
	}

	scanner.Scan()

	instructions, err := parseInstructions(scanner)
	if err != nil {
		return err
	}

	for _, inst := range instructions {
		s.apply(inst)
	}

	fmt.Println(s.topCrates())

	return nil
}

func star2() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	stackLayoutStr, drawingBuf, err := readToStackLayout(scanner)

	stackLayout := strings.Split(strings.TrimSpace(stackLayoutStr), "   ")

	lastStack, err := strconv.Atoi(stackLayout[len(stackLayout)-1])
	if err != nil {
		return err
	}

	s, err := parseDrawing(drawingBuf, lastStack)
	if err != nil {
		return err
	}

	scanner.Scan()

	instructions, err := parseInstructions(scanner)
	if err != nil {
		return err
	}

	for _, inst := range instructions {
		s.apply9001(inst)
	}

	fmt.Println(s.topCrates())

	return nil
}
func readToStackLayout(f *bufio.Scanner) (stackLayoutStr string, drawingBuf bytes.Buffer, err error) {
	for f.Scan() {
		line := f.Text()
		if strings.HasPrefix(line, " 1") {
			stackLayoutStr = line
			goto stackLayoutFound
		}
		drawingBuf.WriteString(line + "\n")
	}
	return "", bytes.Buffer{}, fmt.Errorf("EOF without stack layout")
stackLayoutFound:
	return
}

func parseDrawing(b bytes.Buffer, lastStack int) (*ship, error) {
	scanner := bufio.NewScanner(&b)

	s := newShip(lastStack)

	for scanner.Scan() {
		line := scanner.Text()
		for i := 0; i < (len(line)+1)/4; i++ {
			pos := i * 4
			crateStr := line[pos : pos+3]
			if crateStr[0] == '[' {
				s.stacks[i].pushBot(string(crateStr[1]))
			}
		}
	}

	return s, nil
}

func parseInstructions(scanner *bufio.Scanner) ([]instruction, error) {
	var instructions []instruction
	for scanner.Scan() {
		line := scanner.Text()
		var inst instruction
		if _, err := fmt.Sscanf(line, "move %d from %d to %d", &inst.count, &inst.from, &inst.to); err != nil {
			return nil, err
		}
		instructions = append(instructions, inst)
	}
	return instructions, nil
}
