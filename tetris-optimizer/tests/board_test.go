package tests

import (
	"bytes"
	"os"
	"strings"
	"tetris-optimizer/libs"
	"testing"
)

func TestNewBoard(t *testing.T) {
	size := 4
	board := libs.NewBoard(size)

	if len(board) != size {
		t.Errorf("NewBoard() returned board with %d rows, want %d", len(board), size)
	}

	for i, row := range board {
		if len(row) != size {
			t.Errorf("Row %d has length %d, want %d", i, len(row), size)
		}
		for j, cell := range row {
			if cell != '.' {
				t.Errorf("Cell (%d,%d) = %c, want '.'", i, j, cell)
			}
		}
	}
}

func TestMinBoardSize(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 4},
		{5, 5},
	}

	for _, tt := range tests {
		got := libs.MinBoardSize(tt.n)
		if got != tt.want {
			t.Errorf("MinBoardSize(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestPrintBoard(t *testing.T) {
	board := [][]rune{
		{'A', 'B'},
		{'.', 'C'},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	libs.PrintBoard(board)

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = old

	got := buf.String()
	want := "AB\n.C\n"

	// Normalize line endings
	got = strings.ReplaceAll(got, "\r\n", "\n")
	if got != want {
		t.Errorf("PrintBoard() output =\n%s\nwant:\n%s", got, want)
	}
}
