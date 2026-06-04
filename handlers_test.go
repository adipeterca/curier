package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	dir, _ := os.MkdirTemp("", "curier-test")
	storagePath = dir
	defer os.RemoveAll(dir)
	m.Run()
}

// --- /upload/ endpoint testing ---

func TestUploadValidFile(t *testing.T) {
	body, contentType := makeMultipartFile(t, "test.txt", "hello world")

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	uploadHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["path"] == "" {
		t.Error("expected path in response, got empty string")
	}
}

func TestUploadNoFile(t *testing.T) {
	req := httptest.NewRequest("POST", "/upload", nil)
	w := httptest.NewRecorder()

	uploadHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUploadInvalidFilename(t *testing.T) {
	// filepath.Base(".")  == "." which should be rejected
	body, contentType := makeMultipartFile(t, ".", "hello world")

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	uploadHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// Could also add a simple fuzzing function with testing.F
// just to see how it behaves

// --- /download/ endpoint testing ---

func TestDownloadValidFile(t *testing.T) {
	url := uploadTestFile(t, "hello.txt", "hello world")
	id := extractID(url)

	req := httptest.NewRequest("GET", "/download/"+id, nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()

	downloadHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "hello world") {
		t.Errorf("expected file content in response, got: %s", body)
	}
}

func TestDownloadNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/download/"+strings.Repeat("a", 32), nil)
	req.SetPathValue("id", strings.Repeat("a", 32))
	w := httptest.NewRecorder()

	downloadHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDownloadInvalidID(t *testing.T) {
	req := httptest.NewRequest("GET", "/download/tooshort", nil)
	req.SetPathValue("id", "tooshort")
	w := httptest.NewRecorder()

	downloadHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// --- Helpers ---

// Simplifies the creation of a file for uploads
func makeMultipartFile(t *testing.T, filename, content string) (*bytes.Buffer, string) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(part, content)
	writer.Close()
	return &body, writer.FormDataContentType()
}

// uploadTestFile uploads a file and returns the download URL
func uploadTestFile(t *testing.T, filename, content string) string {
	t.Helper()
	body, contentType := makeMultipartFile(t, filename, content)
	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	uploadHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("upload failed with status %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	return response["path"]
}

// extractID pulls the ID from a full download URL
func extractID(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}
