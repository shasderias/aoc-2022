package main

import (
	"aoc-2022-07/fs"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const inputFilePath = "input.txt"

func main() {
	if err := star1(); err != nil {
		fmt.Println(err)
	}
	if err := star2(); err != nil {
		fmt.Println(err)
	}
}

func star1() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	c := fs.New()

	p := newParser(f, c)
	if err := p.Run(); err != nil {
		return err
	}

	c.CD("/")

	totalSize := 0
	c.Walk(func(n fs.Node, path string) error {
		if _, ok := n.(fs.DirNode); ok && n.Size() <= 100_000 {
			totalSize += n.Size()
		}
		return nil
	})
	fmt.Println(totalSize)

	return nil
}

func star2() error {
	const (
		totalSpace    = 70_000_000
		spaceRequired = 30_000_000
	)

	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	c := fs.New()

	p := newParser(f, c)
	if err := p.Run(); err != nil {
		return err
	}

	c.CD("/")

	root := c.CWD()
	spaceUsed := root.Size()
	spaceFree := totalSpace - spaceUsed
	spaceToFree := spaceRequired - spaceFree

	smallestDir := totalSpace
	c.Walk(func(n fs.Node, path string) error {
		size := n.Size()
		if _, ok := n.(fs.DirNode); ok && size >= spaceToFree && size < smallestDir {
			smallestDir = size
		}
		return nil
	})
	fmt.Println(smallestDir)

	return nil
}

func newParser(r io.Reader, c *fs.Context) *Parser {
	scanner := bufio.NewScanner(r)
	return &Parser{scanner, c}
}

type Parser struct {
	*bufio.Scanner
	*fs.Context
}

func (p *Parser) Run() error {
	const (
		cdPrefix  = "$ cd "
		lsPrefix  = "$ ls"
		dirPrefix = "dir "
	)
	for p.Scan() {
		line := p.Text()
		switch {
		case strings.HasPrefix(line, cdPrefix):
			dirName := strings.TrimPrefix(line, cdPrefix)
			p.CD(dirName)
		case strings.HasPrefix(line, lsPrefix):
		// do nothing
		case strings.HasPrefix(line, dirPrefix):
			dirName := strings.TrimPrefix(line, dirPrefix)
			p.MkDir(dirName)
		default:
			var name string
			var size int
			if _, err := fmt.Sscanf(line, "%d %s", &size, &name); err != nil {
				return err
			}
			p.MkFile(name, size)
		}
	}
	return nil
}
