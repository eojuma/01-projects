package ascii

import "testing"
import "ascii-art/ascii"

func TestBuildASCII(t *testing.T) {
	input := []string{"A"}

	charmap := map[rune][]string{
		'A': {
			"  A  ",
			" A A ",
			"A   A",
			"AAAAA",
			"A   A",
			"A   A",
			"A   A",
			"     ",
		},
	}

	expected := []string{
		"  A  ",
		" A A ",
		"A   A",
		"AAAAA",
		"A   A",
		"A   A",
		"A   A",
		"     ",
	}

	result := ascii.BuildASCII(input, charmap)

	if len(result) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("line %d: expected %q, got %q", i, expected[i], result[i])
		}
	}
}
