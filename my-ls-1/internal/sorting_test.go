package internal

import (
	"os"
	"testing"
	"time"
)

func TestSortAtoZ(t *testing.T) {
	input := []string{"z", "a", "B"}
	expected := []string{"a", "B", "z"}

	SortAtoZ(input)

	for i := range input {
		if input[i] != expected[i] {
			t.Errorf("At index %d: got %s, want %s", i, input[i], expected[i])
		}
	}
}

func TestSortByTime(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-time-")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	f1 := tmpDir + "/old"
	f2 := tmpDir + "/new"

	os.WriteFile(f1, []byte("old"), 0644)
	os.WriteFile(f2, []byte("new"), 0644)

	// Set times manually: f2 is newer
	now := time.Now()
	os.Chtimes(f1, now.Add(-1*time.Hour), now.Add(-1*time.Hour))
	os.Chtimes(f2, now, now)

	files := []string{"old", "new"}
	SortByTime(files, tmpDir)

	if files[0] != "new" {
		t.Errorf("Expected 'new' to be first (newest), got %s", files[0])
	}
}