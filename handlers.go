package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		log.Printf("ERROR: static/index.html not found in embedded files")
		http.Error(w, "Something did not work. Contact the administrator.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(allowedFileExtensions))
	for ext := range allowedFileExtensions {
		keys = append(keys, ext)
	}

	config := Config{
		MaxFileSize:           maxFileSize,
		AllowedFileExtensions: keys,
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		log.Printf("ERROR: failed to marshal config JSON")
		http.Error(w, "Something did not work. Contact the administrator.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(configBytes)
}

// uploadHandler saves the uploaded files to disk, if it passes the verifications.
// It will only return simple errors (JSON based)
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	// A memory buffer limit (at most 32MB) - anything over it will be written to disk, up to maxFileSize bytes.
	var maxRAMSize int64 = min(maxFileSize, 32*1024*1024)
	if err := r.ParseMultipartForm(maxRAMSize); err != nil {
		log.Printf("WARNING: could not parse multipart form: %s\n", err)
		http.Error(w, "file too big", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")

	if err != nil {
		log.Printf("WARNING: no file provided: %s\n", err)
		http.Error(w, "no file provided", http.StatusBadRequest)
		return
	}

	defer file.Close()

	fileName, err := validateFile(header)
	if err != nil {
		log.Printf("WARNING: problem with file upload, reason: %s\n", err)
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}

	// Create a unique filename from 128 random bits -> 32 char string
	id, err := generateId()
	if err != nil {
		log.Printf("ERROR: could not generate ID, reason: %s\n", err)
		http.Error(w, "Something did not work. Contact the administrator.", http.StatusInternalServerError)
		return
	}

	meta := FileMeta{
		OriginalFilename: fileName,
		UploadedAt:       time.Now(),
		UploaderIP:       r.RemoteAddr, // needs testing
	}

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		log.Printf("ERROR: failed to marshal JSON: %s\n", err)
		http.Error(w, "Something did not work. Contact the administrator.", http.StatusInternalServerError)
		return
	}

	metaFilePath := filepath.Join(storagePath, id+".meta")
	err = os.WriteFile(metaFilePath, metaBytes, 0644)
	if err != nil {
		log.Printf("ERROR: could not write %s to disk: %s\n", metaFilePath, err)
		http.Error(w, "Something did not work. Contact the administrator.", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(storagePath, id)
	dst, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)

	if err != nil {
		log.Printf("ERROR: could not create %s on disk: %s\n", filePath, err)
		http.Error(w, "Something did not work. Contact the administrator.", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(metaFilePath)

		log.Printf("ERROR: could not write %s to disk: %s\n", filePath, err)
		http.Error(w, "Something did not work. Contact the administrator.", http.StatusInternalServerError)
		return
	}

	// File was uploaded successfully
	log.Printf("SUCCESS: file %s was saved to disk %s\n", fileName, filePath)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"path": "/share/%s"}`, id)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	err := validateId(id)
	if err != nil {
		log.Printf("WARNING: invalid ID %s, reason: %s\n", id, err)
		serveError(w, http.StatusNotFound)
		return
	}

	meta, err := readMeta(id)
	if err != nil {
		log.Printf("ERROR: reading meta file for %s: %s\n", id, err)
		serveError(w, http.StatusNotFound)
		return
	}

	filePath := filepath.Join(storagePath, id)

	log.Printf("INFO: Serving file %s uploaded by %s to %s\n", meta.OriginalFilename, meta.UploaderIP, r.RemoteAddr)

	w.Header().Set("Content-Disposition", "attachment; filename=\""+meta.OriginalFilename+"\"")
	http.ServeFile(w, r, filePath)
}

func shareHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	err := validateId(id)
	if err != nil {
		log.Printf("WARNING: invalid ID %s, reason: %s\n", id, err)
		// We return 404-NotFound because this is a browser-client facing application, not a pure API.
		// Thus, it makes more sense for the client to get a 404 instead of a 400-BadRequest
		serveError(w, http.StatusNotFound)
		return
	}

	meta, err := readMeta(id)
	if err != nil {
		log.Printf("ERROR: reading meta file for %s, reason: %s\n", id, err)
		serveError(w, http.StatusNotFound)
		return
	}
	shareData := ShareData{
		FileMeta: *meta,
		ID:       id,
	}

	w.Header().Set("Content-Type", "text/html")

	if err = shareTemplate.Execute(w, shareData); err != nil {
		log.Printf("ERROR: could not execute template: %s\n", err)
		// We don't return an error to the client, as partial output may already
		// have been written.
	}
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/styles.css")
	if err != nil {
		log.Printf("ERROR: static/styles.css not found in embedded files")
		http.Error(w, "Something did not work. Contact the administrator.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(data)
}

// --- Helper functions ---

func validateId(id string) error {

	id = filepath.Base(id) // Simple check against path traversal
	if len(id) != 32 {
		return fmt.Errorf("ID length is not 32")
	}

	for _, ch := range id {
		if !(ch >= 'a' && ch <= 'f' || ch >= '0' && ch <= '9') {
			return fmt.Errorf("ID is not a HEX encoded string")
		}
	}

	return nil
}

func validateFile(header *multipart.FileHeader) (string, error) {
	fileName := filepath.Base(header.Filename)
	if fileName == "." || fileName == "" {
		return "", fmt.Errorf("Filename is either empty or '.'")
	}

	if header.Size > maxFileSize {
		return "", fmt.Errorf("File size is bigger than allowed")
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	if !allowedFileExtensions[ext] {
		return "", fmt.Errorf("file extension %s is not allowed", ext)
	}

	return fileName, nil
}

func generateId() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		// Docs say: "It never returns an error, and always fills b entirely."
		// If you get here, I dunno what to do
		return "", err
	}
	id := hex.EncodeToString(bytes)

	return id, nil
}

func readMeta(id string) (*FileMeta, error) {

	metaFilePath := filepath.Join(storagePath, id+".meta")

	metaBytes, err := os.ReadFile(metaFilePath)
	if err != nil {
		return nil, err
	}
	var meta FileMeta
	if err = json.Unmarshal(metaBytes, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func serveError(w http.ResponseWriter, code int) {
	var file string
	switch code {
	case http.StatusNotFound:
		file = "static/404.html"
	case http.StatusInternalServerError:
		file = "static/500.html"
	default:
		file = "static/404.html"
	}

	content, err := staticFiles.ReadFile(file)
	if err != nil {
		// fallback if even the error page is missing
		log.Printf("CRITICAL: Missing error page for %d (are you sure the file was embedded?), reason: %s\n", code, err)
		http.Error(w, "Something really broke me :(", code)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	w.Write(content)
}
