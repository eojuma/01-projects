package backend

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// initTestTemplates sets up minimal templates so the handlers can render
// without panicking. This must be called at the start of every handler test.
func initTestTemplates() {
	funcMap := template.FuncMap{
		"formatLoc": func(s string) string {
			s = strings.ReplaceAll(s, "_", " ")
			s = strings.ReplaceAll(s, "-", ", ")
			return s
		},
	}
	Templates = template.Must(template.New("").Funcs(funcMap).Parse(
		`{{define "index.html"}}OK{{end}}` +
			`{{define "details.html"}}DETAIL{{end}}` +
			`{{define "error.html"}}ERROR {{.StatusCode}}{{end}}`,
	))
}

//ErrorHandler tests

func TestErrorHandler_NotFound(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	rr := httptest.NewRecorder()

	ErrorHandler(rr, req, http.StatusNotFound)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestErrorHandler_InternalServerError(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	ErrorHandler(rr, req, http.StatusInternalServerError)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestErrorHandler_MethodNotAllowed(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()

	ErrorHandler(rr, req, http.StatusMethodNotAllowed)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestErrorHandler_BadRequest(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodGet, "/details?id=abc", nil)
	rr := httptest.NewRecorder()

	ErrorHandler(rr, req, http.StatusBadRequest)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestErrorHandler_UnknownStatus(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	ErrorHandler(rr, req, http.StatusTeapot)

	if rr.Code != http.StatusTeapot {
		t.Errorf("expected status %d, got %d", http.StatusTeapot, rr.Code)
	}
}

//HomeHandler tests

func TestHomeHandler_WrongPath(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rr := httptest.NewRecorder()

	HomeHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d for wrong path, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHomeHandler_PostMethod(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()

	HomeHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d for POST, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

//DetailsHandler tests

func TestDetailsHandler_PostMethod(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodPost, "/details?id=1", nil)
	rr := httptest.NewRecorder()

	DetailsHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d for POST, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestDetailsHandler_MissingID(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodGet, "/details", nil)
	rr := httptest.NewRecorder()

	DetailsHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for missing id, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDetailsHandler_InvalidID(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodGet, "/details?id=abc", nil)
	rr := httptest.NewRecorder()

	DetailsHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for non-numeric id, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDetailsHandler_NegativeID(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodGet, "/details?id=-5", nil)
	rr := httptest.NewRecorder()

	DetailsHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for negative id, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDetailsHandler_ZeroID(t *testing.T) {
	initTestTemplates()
	req := httptest.NewRequest(http.MethodGet, "/details?id=0", nil)
	rr := httptest.NewRecorder()

	DetailsHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for zero id, got %d", http.StatusBadRequest, rr.Code)
	}
}
