package tests

import (
	"testing"
	"tetris-optimizer/libs"
)

func TestIsValidFormat(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected bool
	}{
		{
			"Valid Tetromino",
			[]string{"#...", "#...", "#...", "#..."},
			true,
		},
		{
			"Invalid Tetromino - wrong size",
			[]string{"#...", "#...", "#..."},
			false,
		},
		{
			"Invalid Tetromino - extra char",
			[]string{"#..X", "#...", "#...", "#..."},
			false,
		},
		{
			"Invalid Tetromino - wrong number of #",
			[]string{"#...", "#...", "#...", "...."},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := libs.IsValidFormat(tt.lines)
			if result != tt.expected {
				t.Errorf("IsValidFormat(%v) = %v; want %v", tt.lines, result, tt.expected)
			}
		})
	}
}

func TestIsConnected(t *testing.T) {
	tests := []struct {
		name     string
		blocks   []libs.Point
		expected bool
	}{
		{
			"Valid vertical line",
			[]libs.Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
			true,
		},
		{
			"Valid square",
			[]libs.Point{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			true,
		},
		{
			"Disconnected blocks",
			[]libs.Point{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
			false,
		},
		{
			"Valid L shape",
			[]libs.Point{{0, 0}, {0, 1}, {0, 2}, {1, 2}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := libs.IsConnected(tt.blocks)
			if result != tt.expected {
				t.Errorf("IsConnected(%v) = %v; want %v", tt.blocks, result, tt.expected)
			}
		})
	}
}
