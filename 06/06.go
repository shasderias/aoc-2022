package main

import (
	"fmt"
	"io"
	"os"
)

const inputFilePath = "input.txt"

func main() {
	star1()
	star2()
}

type sopDetector struct {
	buf [4]byte
}

func (d *sopDetector) isStartOfPacket() bool {
	if d.buf[0] == d.buf[1] || d.buf[0] == d.buf[2] || d.buf[0] == d.buf[3] || d.buf[1] == d.buf[2] || d.buf[1] == d.buf[3] || d.buf[2] == d.buf[3] {
		return false
	}
	return true
}

func (d *sopDetector) Write(p []byte) (n int, err error) {
	orgLen := len(p)

	if len(p) > 4 {
		p = p[len(p)-4:]
	}

	copy(d.buf[:], d.buf[len(p):])
	copy(d.buf[4-len(p):], p)

	return orgLen, nil
}

func star1() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	d := &sopDetector{}

	if _, err := io.CopyN(d, f, 4); err != nil {
		return err
	}

	processed := 4
	for ; !d.isStartOfPacket(); processed++ {
		if _, err := io.CopyN(d, f, 1); err != nil {
			return err
		}
	}

	fmt.Println(string(d.buf[:]), processed)

	return nil
}

type somDetector struct {
	buf     []byte
	somSize int
}

func newSOMDetector(somSize int) *somDetector {
	return &somDetector{
		buf:     make([]byte, somSize),
		somSize: somSize,
	}
}

func (d *somDetector) Write(p []byte) (n int, err error) {
	orgLen := len(p)
	s := d.somSize

	if len(p) > s {
		p = p[len(p)-s:]
	}

	copy(d.buf[:], d.buf[len(p):])
	copy(d.buf[s-len(p):], p)

	return orgLen, nil
}

func (d *somDetector) isStartOfMessage() bool {
	seen := map[byte]struct{}{}

	for _, b := range d.buf {
		if _, ok := seen[b]; ok {
			return false
		}
		seen[b] = struct{}{}
	}
	return true
}

func star2() error {
	f, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	d := newSOMDetector(14)

	if _, err := io.CopyN(d, f, 14); err != nil {
		return err
	}

	processed := 14
	for ; !d.isStartOfMessage(); processed++ {
		if _, err := io.CopyN(d, f, 1); err != nil {
			return err
		}
	}

	fmt.Println(string(d.buf[:]), processed)

	return nil
}
