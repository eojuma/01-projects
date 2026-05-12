package internal

import (
	"os"
	"testing"
)

func TestPrintShortNoPanic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-print-")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fileName := "test.txt"
	f, _ := os.Create(tmpDir + "/" + fileName)
	f.Close()

	// Should not panic
	PrintShort(tmpDir, []string{fileName})
}

func TestPrintLongNoPanic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-print-long-")
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fileName := "test.txt"
	f, _ := os.Create(tmpDir + "/" + fileName)
	f.Close()

	// Should not panic
	PrintLong(tmpDir, []string{fileName}, true)
}