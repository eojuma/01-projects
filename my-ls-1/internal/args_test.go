package internal

import (
	"os"
	"testing"
)

func TestGetInputs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *Options
	}{
		{
			name: "no arguments defaults to current directory",
			args: []string{},
			expected: &Options{
				Dirs: []string{"."},
				All:  false, Long: false, Reverse: false, Time: false,
			},
		},
		{
			name: "hyphen as directory",
			args: []string{"-"},
			expected: &Options{
				Dirs: []string{"-"},
				All:  false, Long: false, Reverse: false, Time: false,
			},
		},
		{
			name: "combined flags and path",
			args: []string{"-la", "myDir"},
			expected: &Options{
				Dirs: []string{"myDir"},
				All:  true,
				Long: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = append([]string{"cmd"}, tt.args...)
			result := GetInputs()

			if len(result.Dirs) != len(tt.expected.Dirs) {
				t.Errorf("Dirs count: got %v, want %v", result.Dirs, tt.expected.Dirs)
			}
			if result.All != tt.expected.All || result.Long != tt.expected.Long {
				t.Errorf("Flags mismatch for %s", tt.name)
			}
		})
	}
}