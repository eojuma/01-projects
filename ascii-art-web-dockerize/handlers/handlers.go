package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"learn.zone01kisumu.ke/git/aochuka/ascii-art/ascii"
)

var Tmpl *template.Template

func RunWeb() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/ascii-art", AsciiArtHandler)

	fmt.Println("Server starting at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	if Tmpl == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	err := Tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	if Tmpl == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	text := r.FormValue("text")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	banner := r.FormValue("banner")

	if text == "" {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	result, err := ascii.Generate(text, banner)
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	// --- FIXED: join lines with newline and execute template only once ---
	resultStr := strings.Join(result, "\n")
	w.WriteHeader(http.StatusOK)
	err = Tmpl.Execute(w, map[string]any{
		"Result": resultStr,
	})
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}
