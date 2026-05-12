// main.go
package main

import (
	"fmt"
	"my-ls-1/internal"
	"os"
	"strings"
)

func main() {
	inputs := internal.GetInputs()
	var filesToPrint []string
	var dirsToPrint []string

	for _, path := range inputs.Dirs {
		// 1. Get raw info (don't follow links yet)
		info, err := os.Lstat(path)
		if err != nil {
			// Handle trailing slash error for missing directories
			if strings.HasSuffix(path, "/") {
				basePath := strings.TrimSuffix(path, "/")
				if baseInfo, baseErr := os.Lstat(basePath); baseErr == nil && !baseInfo.IsDir() {
					fmt.Printf("ls: cannot access '%s': Not a directory\n", path)
					continue
				}
			}
			fmt.Printf("ls: cannot access '%s': No such file or directory\n", path)
			continue
		}

		// 2. Explicitly check for trailing slash error for existing objects
		if strings.HasSuffix(path, "/") && !info.IsDir() {
			realInfo, err := os.Stat(path)
			if err != nil || !realInfo.IsDir() {
				fmt.Printf("ls: cannot access '%s': Not a directory\n", path)
				continue
			}
		}

		// 3. Determine if we should treat this as a directory or a file
		isLink := info.Mode()&os.ModeSymlink != 0
		shouldFollow := info.IsDir() || (isLink && (strings.HasSuffix(path, "/") || !inputs.Long))

		if shouldFollow {
			realInfo, err := os.Stat(path)
			if err == nil && realInfo.IsDir() {
				dirsToPrint = append(dirsToPrint, path)
				continue
			}
		}

		filesToPrint = append(filesToPrint, path)
	}

	// Apply sorting to arguments
	internal.SortAtoZ(filesToPrint)
	internal.SortAtoZ(dirsToPrint)

	// Process files first
	if len(filesToPrint) > 0 {
		internal.PrintLong(".", filesToPrint, false)
	}

	// Process directories
	for i, path := range dirsToPrint {
		// Print a newline between arguments if there are multiple
		if len(filesToPrint) > 0 || i > 0 {
			fmt.Println()
		}
		
		// Header for multiple paths or recursion
		// The system ls prints "path:" if there's more than one argument
		if inputs.Recursive || (len(filesToPrint)+len(dirsToPrint) > 1) {
			fmt.Printf("%s:\n", path)
		}
		
		listDir(path, inputs, false)
	}
}

func listDir(path string, inputs *internal.Options, isSubDir bool) {
	if isSubDir {
		fmt.Printf("\n%s:\n", path)
	}

	files, err := internal.ReadDirectory(path, inputs)
	if err != nil {
		fmt.Printf("ls: cannot open directory '%s': %v\n", path, err)
		return
	}

	// If directory is empty and not in long format, do nothing (matches system ls)
	if len(files) == 0 && !inputs.Long {
		return
	}

	// Apply sorting logic
	if inputs.Time {
		internal.SortByTime(files, path)
	} else {
		internal.SortAtoZ(files)
	}
	
	if inputs.Reverse {
		internal.ReverseSort(files)
	}

	// Printing
	if inputs.Long {
		internal.PrintLong(path, files, true)
	} else {
		internal.PrintShort(path, files)
	}

	// Recursion
	if inputs.Recursive {
		for _, name := range files {
			if name == "." || name == ".." {
				continue
			}

			fullPath := internal.JoinPaths(path, name)
			fInfo, err := os.Lstat(fullPath)
			if err == nil && fInfo.IsDir() {
				listDir(fullPath, inputs, true)
			}
		}
	}
}