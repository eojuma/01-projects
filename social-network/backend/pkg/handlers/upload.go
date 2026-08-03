package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

var allowedMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
}

// ErrUnsupportedMediaType signals the caller should return HTTP 415.
var ErrUnsupportedMediaType = errors.New("unsupported media type")

func SaveUploadedImage(r *http.Request, fieldName, uploadDir string) (string, error) {
	file, _, err := r.FormFile(fieldName)
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	defer file.Close()

	// Read first 512 bytes to detect MIME type from magic bytes
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	mimeType := http.DetectContentType(buf[:n])
	ext, allowed := allowedMIME[mimeType]
	if !allowed {
		return "", ErrUnsupportedMediaType
	}

	// Seek back to beginning before copying full contents

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to seek file: %w", err)
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload dir: %w", err)
	}

	filename := uuid.Must(uuid.NewV4()).String() + ext
	savePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return savePath, nil
}

// SaveUploadedImages saves multiple uploaded files from the given field name and
// returns the saved paths (relative to the server working directory).
func SaveUploadedImages(r *http.Request, fieldName, uploadDir string) ([]string, error) {
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, nil
	}
	files := r.MultipartForm.File[fieldName]
	if len(files) == 0 {
		return nil, nil
	}
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload dir: %w", err)
	}

	var saved []string
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open uploaded file: %w", err)
		}
		// read header bytes
		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		mimeType := http.DetectContentType(buf[:n])
		ext, allowed := allowedMIME[mimeType]
		if !allowed {
			f.Close()
			return nil, ErrUnsupportedMediaType
		}
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			f.Close()
			return nil, fmt.Errorf("failed to seek file: %w", err)
		}

		filename := uuid.Must(uuid.NewV4()).String() + ext
		savePath := filepath.Join(uploadDir, filename)
		dst, err := os.Create(savePath)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("failed to save file: %w", err)
		}
		if _, err := io.Copy(dst, f); err != nil {
			dst.Close()
			f.Close()
			return nil, fmt.Errorf("failed to write file: %w", err)
		}
		dst.Close()
		f.Close()
		// Return path relative to server root so /uploads/... proxy works
		saved = append(saved, savePath)
	}
	return saved, nil
}
