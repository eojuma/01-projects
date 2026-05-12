package main

import (
	"fmt"
	"groupie-tracker/backend"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func main() {
	// 1. Define custom template functions for formatting data
	funcMap := template.FuncMap{
		"formatLoc": func(s string) string {
			s = strings.ReplaceAll(s, "_", " ")
			s = strings.ReplaceAll(s, "-", ", ")
			return strings.Title(s)
		},
	}

	// 2. Parse all templates in the frontend folder with the custom functions
	var err error
	backend.Templates, err = template.New("").Funcs(funcMap).ParseGlob("frontend/*.html")
	if err != nil {
		log.Fatalf("Critical Error: Failed to parse templates: %v", err)
	}

	// 3. Serve static files (CSS, Images) 
	// This maps the physical folder "frontend/static" to the URL path "/static/"
	fs := http.FileServer(http.Dir("frontend/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 4. Register Routes
	// Root route handles the search and the main grid
	http.HandleFunc("/", backend.HomeHandler)
	
	// Details route handles the "clickable card" event
	// This will look for the DetailsHandler function in your backend package
	http.HandleFunc("/details", backend.DetailsHandler)

	// 5. Start the Server
	port := ":8080"
	fmt.Printf("🎵 Groupie Tracker is live at http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}