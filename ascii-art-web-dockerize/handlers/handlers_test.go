package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"html/template"
)

// --- helper to setup template ---
func setupTemplate(t *testing.T) {
	var err error
	Tmpl, err = template.New("test").Parse("<html>{{.Result}}</html>")
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}
}

// --- Minimal GET test for "/" ---
func TestHomeHandler_Success(t *testing.T) {
	setupTemplate(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	HomeHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
}

// --- GET wrong method ---
func TestHomeHandler_WrongMethod(t *testing.T) {
	setupTemplate(t)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	HomeHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", resp.StatusCode)
	}
}

// --- Minimal POST test for "/ascii-art" ---
func TestAsciiArtHandler_Success(t *testing.T) {
	setupTemplate(t)

	body := strings.NewReader("text=Hello&banner=standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	AsciiArtHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
}

// --- POST empty text ---
func TestAsciiArtHandler_EmptyText(t *testing.T) {
	setupTemplate(t)

	body := strings.NewReader("text=&banner=standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	AsciiArtHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", resp.StatusCode)
	}
}

// --- POST wrong method ---
func TestAsciiArtHandler_WrongMethod(t *testing.T) {
	setupTemplate(t)

	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	w := httptest.NewRecorder()

	AsciiArtHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", resp.StatusCode)
	}
}

// --- Missing template (simulate internal server error) ---
func TestAsciiArtHandler_TemplateMissing(t *testing.T) {
	Tmpl = nil // simulate missing template

	body := strings.NewReader("text=Hello&banner=standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	AsciiArtHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found due to missing template, got %d", resp.StatusCode)
	}
}

// --- Simulate ascii.Generate failure (invalid banner) ---
func TestAsciiArtHandler_GenerateFail(t *testing.T) {
	setupTemplate(t)

	body := strings.NewReader("text=Hello&banner=invalidbanner")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	AsciiArtHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found due to ascii.Generate failure, got %d", resp.StatusCode)
	}
}
