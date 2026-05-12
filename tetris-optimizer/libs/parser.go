package libs

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func ParseFile(path string) ([]Tetromino, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var blocks [][]string
	var current []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			if len(current) > 0 {
				blocks = append(blocks, current)
				current = nil
			}
		} else {
			current = append(current, line)
		}
	}

	if len(current) > 0 {
		blocks = append(blocks, current)
	}
	var tetros []Tetromino
	letter := 'A'

	for _, b := range blocks {
		t, err := ParseBlock(b, letter)
		if err != nil {
			return nil, err
		}
		tetros = append(tetros, t)
		letter++
	}
	return tetros, nil
}

func ParseBlock(lines []string, letter rune) (Tetromino, error) {
	if !IsValidFormat(lines) {
		return Tetromino{}, errors.New("Invalid Format")
	}

	var blocks []Point
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if lines[y][x] == '#' {
				blocks = append(blocks, Point{X: x, Y: y})
			}
		}
	}

	if !IsConnected(blocks) {
		return Tetromino{}, errors.New("Not connected")
	}
	blocks = Normalize(blocks)
	return Tetromino{
		Blocks: blocks,
		Letter: letter,
	}, nil
}
