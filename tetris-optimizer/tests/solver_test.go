package tests

import (
	"reflect"
	"testing"
	"tetris-optimizer/libs"
)

func TestCanPlacePlaceRemove(t *testing.T) {
	// Create a simple 4x4 board
	board := libs.NewBoard(4)

	// Define a simple tetromino (2x2 square)
	tetro := libs.Tetromino{
		Blocks: []libs.Point{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
		},
		Letter: 'A',
	}

	// CanPlace should succeed at 0,0
	if !libs.CanPlace(board, tetro, 0, 0) {
		t.Errorf("CanPlace() failed at valid position (0,0)")
	}

	// Place the tetromino
	libs.Place(board, tetro, 0, 0)

	expected := [][]rune{
		{'A', 'A', '.', '.'},
		{'A', 'A', '.', '.'},
		{'.', '.', '.', '.'},
		{'.', '.', '.', '.'},
	}

	if !reflect.DeepEqual(board, expected) {
		t.Errorf("Place() failed. Got %+v, want %+v", board, expected)
	}

	// CanPlace should now fail at overlapping position
	if libs.CanPlace(board, tetro, 0, 0) {
		t.Errorf("CanPlace() returned true for overlapping position")
	}

	// Remove the tetromino
	libs.Remove(board, tetro, 0, 0)
	boardExpectedEmpty := libs.NewBoard(4)
	if !reflect.DeepEqual(board, boardExpectedEmpty) {
		t.Errorf("Remove() failed. Got %+v, want %+v", board, boardExpectedEmpty)
	}
}

func TestSolveSimple(t *testing.T) {
	// Two simple tetrominoes: two 2x2 squares
	tetros := []libs.Tetromino{
		{
			Blocks: []libs.Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
			Letter: 'A',
		},
		{
			Blocks: []libs.Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
			Letter: 'B',
		},
	}

	board := libs.Solve(tetros)
	// Board should be at least 3x3 to fit two squares (minimal square is 3x3)
	if len(board) < 3 || len(board[0]) < 3 {
		t.Errorf("Solve() returned too small board: %d x %d", len(board), len(board[0]))
	}

	// Check that all letters are present
	foundA := false
	foundB := false
	for _, row := range board {
		for _, cell := range row {
			if cell == 'A' {
				foundA = true
			}
			if cell == 'B' {
				foundB = true
			}
		}
	}

	if !foundA || !foundB {
		t.Errorf("Solve() did not place all tetrominoes, foundA=%v, foundB=%v", foundA, foundB)
	}
}
