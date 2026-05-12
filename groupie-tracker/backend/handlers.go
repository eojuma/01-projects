package backend

import (
	"html/template" 
	"log"
	"net/http"
	"strconv"
)

var Templates *template.Template

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	artists, err := GetArtistsWithRelations()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	query := r.URL.Query().Get("query")
	filtered := SearchArtists(query, artists)

	data := PageData{
		Artists: filtered,
		Query:   query,
		Total:   len(filtered),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := Templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func DetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	artists, err := GetArtistsWithRelations()
	if err != nil {
		log.Printf("Error fetching data for details: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	var targetArtist *ArtistDetail 
	for i := range artists {
		if artists[i].ID == id {
			targetArtist = &artists[i]
			break
		}
	}

	if targetArtist == nil {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := Templates.ExecuteTemplate(w, "details.html", targetArtist); err != nil {
		log.Printf("Error rendering details template: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	messages := map[int]string{
		http.StatusNotFound:            "The page you're looking for doesn't exist.",
		http.StatusInternalServerError: "Something went wrong on our end.",
		http.StatusMethodNotAllowed:    "This request method is not allowed.",
		http.StatusBadRequest:          "Bad request.",
	}

	msg, ok := messages[status]
	if !ok {
		msg = "An unexpected error occurred."
	}

	data := ErrorData{
		StatusCode: status,
		Message:    msg,
	}

	if err := Templates.ExecuteTemplate(w, "error.html", data); err != nil {
		log.Printf("Error rendering error template: %v", err)
		http.Error(w, msg, status)
	}
}