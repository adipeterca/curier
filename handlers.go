package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		fmt.Printf("ERROR: static/index.html not found in embedded files")
		http.Error(w, "Sorry :(", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Printf("WARNING: invalid file upload: %s\n", err)
		http.Error(w, "invalid file upload", http.StatusBadRequest)
		return
	}

	defer file.Close()

	fileName := filepath.Base(header.Filename)
	if fileName == "." || fileName == "" {
		fmt.Printf("WARNING: invalid filename provided: %s\n", err)
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	// Could be made to only accept specific files, based on their magic bytes
	// some ideas: ZIP, RAR, PNG, JPG, JPEG, text files (how do I recognise these - maybe by parsing the first 4 bytes?)

	// Create a unique filename from 128 random bits -> 32 char string
	bytes := make([]byte, 16)
	_, err = rand.Read(bytes)
	if err != nil {
		fmt.Printf("ERROR: could not generate ID, reason: %s\n", err)
		http.Error(w, "Sorry, I did my best :(", http.StatusInternalServerError)
		return
	}
	id := hex.EncodeToString(bytes)

	meta := FileMeta{
		OriginalFilename: fileName,
		UploadedAt:       time.Now(),
		UploaderIP:       r.RemoteAddr, // needs testing
	}

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		fmt.Printf("ERROR: failed to marshal JSON: %s\n", err)
		http.Error(w, "Sorry, I did my best :(", http.StatusInternalServerError)
		return
	}

	metaFilePath := filepath.Join(storagePath, id+".meta")
	err = os.WriteFile(metaFilePath, metaBytes, 0644)
	if err != nil {
		fmt.Printf("ERROR: could not write %s to disk: %s\n", metaFilePath, err)
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(storagePath, id)
	dst, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)

	if err != nil {
		fmt.Printf("ERROR: could not create %s on disk: %s\n", filePath, err)
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(metaFilePath)

		fmt.Printf("ERROR: could not write %s to disk: %s\n", filePath, err)
		http.Error(w, "Could not write file", http.StatusInternalServerError)
		return
	}

	// File was uploaded successfully
	fmt.Printf("SUCCESS: file %s was saved to disk %s\n", fileName, filePath)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"url": "%s:%s/share/%s"}`, host, port, id)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	err := validateId(id)
	if err != nil {
		fmt.Printf("WARNING: invalid ID %s, reason: %s\n", id, err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	meta, err := readMeta(id)
	if err != nil {
		fmt.Printf("ERROR: reading meta file for %s: %s\n", id, err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	filePath := filepath.Join(storagePath, id)

	fmt.Printf("INFO: Serving file %s uploaded by %s to %s\n", meta.OriginalFilename, meta.UploaderIP, r.RemoteAddr)

	w.Header().Set("Content-Disposition", "attachment; filename=\""+meta.OriginalFilename+"\"")
	http.ServeFile(w, r, filePath)
}

func shareHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	err := validateId(id)
	if err != nil {
		fmt.Printf("WARNING: invalid ID %s, reason: %s\n", id, err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	meta, err := readMeta(id)
	if err != nil {
		fmt.Printf("ERROR: reading meta file for %s, reason: %s\n", id, err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	shareData := ShareData{
		FileMeta: *meta,
		ID:       id,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.ParseFS(templateFiles, "templates/share.html")

	if err != nil {
		fmt.Printf("ERROR: could not parse template: %s\n", err)
		http.Error(w, "Sorry, I did my best :(", http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, shareData); err != nil {
		fmt.Printf("ERROR: could not execute template: %s\n", err)
		// We don't return an error to the client, as partial output may already
		// have been written.
	}
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/styles.css")
	if err != nil {
		fmt.Printf("ERROR: static/styles.css not found in embedded files")
		http.Error(w, "Sorry :(", http.StatusInternalServerError)
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

	return nil
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
