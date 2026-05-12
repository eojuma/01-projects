package tests

import (
	"os"
	"testing"
	"tetris-optimizer/libs"
)

func TestParseBlock(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		letter  rune
		wantErr bool
	}{
		{
			"Valid vertical line",
			[]string{"#...", "#...", "#...", "#..."},
			'A',
			false,
		},
		{
			"Valid square",
			[]string{"##..", "##..", "....", "...."},
			'B',
			false,
		},
		{
			"Invalid format - wrong char",
			[]string{"#..X", "#...", "#...", "#..."},
			'C',
			true,
		},
		{
			"Disconnected blocks",
			[]string{"#...", ".#..", "..#.", "...#"},
			'D',
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := libs.ParseBlock(tt.lines, tt.letter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	// Create a temporary file with multiple tetrominoes
	content := `#...
#...
#...
#...

##..
##..
....
....

.###
...#
....
....
`

	tmpFile, err := os.CreateTemp("", "tetromino_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	tetros, err := libs.ParseFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseFile() returned error: %v", err)
	}

	if len(tetros) != 3 {
		t.Fatalf("ParseFile() returned %d tetrominoes, want 3", len(tetros))
	}

	// Check that letters are assigned correctly
	letters := []rune{'A', 'B', 'C'}
	for i, tet := range tetros {
		if tet.Letter != letters[i] {
			t.Errorf("Tetromino %d has letter %c, want %c", i, tet.Letter, letters[i])
		}
	}

	// Check that blocks are normalized (minX and minY = 0)
	for _, tet := range tetros {
		for _, b := range tet.Blocks {
			if b.X < 0 || b.Y < 0 {
				t.Errorf("Tetromino %c has block with negative coordinate: %v", tet.Letter, b)
			}
		}
	}
}
