package internal

import (
	"os"
	"strings"
)

func ReadDirectory(path string, opt *Options) ([]string, error) {

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []string
	// The standard ls -a always starts with these two literal entries
	if opt.All {
		files = append(files, ".", "..")
	}

	for _, entry := range entries {
		if !opt.All && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		files = append(files, entry.Name())
	}

	return files, nil
}
