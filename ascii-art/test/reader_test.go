package ascii

import "testing"
import "ascii-art/ascii"

func TestReader(t *testing.T) {
	lines, _ := ascii.ReadLinesFromFile("../standard.txt")
	got := len(lines)

	want := 855

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
