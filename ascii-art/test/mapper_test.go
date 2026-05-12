package ascii

import "testing"
import "ascii-art/ascii"

func TestMapToASCII(t *testing.T) {
	// Each char uses 9 lines, but only 8 are stored enough lines for all ASCII printable chars
	totalChars := int('~' - ' ' + 1)
	totalLines := 1 + totalChars*9
	lines := make([]string, totalLines)
	for i := 0; i < len(lines); i++ {
		lines[i] = "line"
	}

	charmap := ascii.MapToASCII(lines)

	//  Check number of mapped characters
	if len(charmap) != totalChars {
		t.Fatalf("expected %d characters, got %d", totalChars, len(charmap))
	}

	//  Check each character has exactly 8 lines
	for char := ' '; char <= '~'; char++ {
		asciiArt, ok := charmap[char]
		if !ok {
			t.Fatalf("missing character: %q", char)
		}

		if len(asciiArt) != 8 {
			t.Errorf("character %q has %d lines, expected 8", char, len(asciiArt))
		}
	}
}
