package internal

import (
	"os"
	"testing"
)

func TestReadDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files using string concatenation
	testFiles := []string{"file1.txt", ".hidden.txt"}
	for _, name := range testFiles {
		f, err := os.Create(tmpDir + "/" + name)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		f.Close()
	}

	tests := []struct {
		name      string
		opt       *Options
		expectLen int // file1.txt (+ . and .. if -a)
	}{
		{
			name:      "hide dot files",
			opt:       &Options{All: false},
			expectLen: 1,
		},
		{
			name:      "show all",
			opt:       &Options{All: true},
			expectLen: 4, // ., .., file1.txt, .hidden.txt
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := ReadDirectory(tmpDir, tt.opt)
			if err != nil {
				t.Fatalf("ReadDirectory failed: %v", err)
			}
			if len(files) != tt.expectLen {
				t.Errorf("Expected %d files, got %d: %v", tt.expectLen, len(files), files)
			}
		})
	}
}