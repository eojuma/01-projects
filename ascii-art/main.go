package main

import (
	"ascii-art/ascii"
	"fmt"
	"os"
)

func main() {
	// Validate arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./ascii-art 'text' banner")
		os.Exit(1)
	}

	input := os.Args[1]
	banner := os.Args[2]

	// Call the reusable Generate function
	art, err := ascii.Generate(input, banner)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Print each line of the generated ASCII art
	for _, line := range art {
		fmt.Println(line)
	}
}
