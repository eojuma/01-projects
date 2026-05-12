package main

import (
	"fmt"
	"os"
	"html/template"
	"ascii-art-web/handlers"
    "learn.zone01kisumu.ke/git/aochuka/ascii-art/ascii"
)

func main() {
	var err error
	handlers.Tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	input := ""
	banner := ""

	// -- INPUT FROM THE CLI
	if len(os.Args) > 1 {
		if len(os.Args) != 3 {
			fmt.Println("Usage: ./ascii-art 'text' banner")
			os.Exit(1)
		}

		input = os.Args[1]
		banner = os.Args[2]

		art, err := ascii.Generate(input, banner)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		for _, line := range art {
			fmt.Println(line)
		}
		return
	}

	// -- INPUT FROM THE WEB
	handlers.RunWeb()
}
