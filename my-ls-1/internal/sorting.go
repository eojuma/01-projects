package internal

import (
	"os"
	"strings"
	"time"
)

func SortAtoZ(names []string) {
	for i := 0; i < len(names)-1; i++ {
		for j := 0; j < len(names)-1-i; j++ {
			left := strings.TrimPrefix(names[j], ".")
			right := strings.TrimPrefix(names[j+1], ".")

			if strings.ToLower(left) > strings.ToLower(right) {
				names[j], names[j+1] = names[j+1], names[j]
			}
		}
	}
}

func SortByTime(files []string, path string) {
	type fileData struct {
		name string
		time time.Time
	}

	fileInfo := make([]fileData, 0, len(files))

	for _, file := range files {
		fullpath := JoinPaths(path, file)
		// Use Lstat to stay consistent with system ls behavior
		info, err := os.Lstat(fullpath)
		if err != nil {
			fileInfo = append(fileInfo, fileData{file, time.Unix(0, 0)})
			continue
		}
		fileInfo = append(fileInfo, fileData{file, info.ModTime()})
	}

	// Manual Bubble Sort to replace sort.SliceStable
	n := len(fileInfo)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			shouldSwap := false

			// 1. Primary Sort: Modification Time (Newest First)
			if fileInfo[j].time.Before(fileInfo[j+1].time) {
				shouldSwap = true
			} else if fileInfo[j].time.Equal(fileInfo[j+1].time) {
				// 2. Secondary Sort: Case-Insensitive Alphabetical Tie-breaker
				// This fixes the README.md vs go.mod issue
				nameJ := strings.ToLower(fileInfo[j].name)
				nameJ1 := strings.ToLower(fileInfo[j+1].name)
				
				if nameJ > nameJ1 {
					shouldSwap = true
				} else if nameJ == nameJ1 {
					// 3. Final Tie-breaker: Case-Sensitive (ASCII)
					if fileInfo[j].name > fileInfo[j+1].name {
						shouldSwap = true
					}
				}
			}

			if shouldSwap {
				fileInfo[j], fileInfo[j+1] = fileInfo[j+1], fileInfo[j]
			}
		}
	}

	for i := range files {
		files[i] = fileInfo[i].name
	}
}

func ReverseSort(files []string) {
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}
}

func JoinPaths(dir, file string) string {
    if dir == "." {
        return file
    }
    if len(dir) > 0 && dir[len(dir)-1] == '/' {
        return dir + file
    }
    return dir + "/" + file
}