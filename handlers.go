package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Printf("WARNING: invalid file upload: %s\n", err)
		http.Error(w, "invalid file upload", http.StatusBadRequest)
		return
	}

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
		fmt.Printf("ERROR: could not generate ID: %s\n", err)
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
	defer file.Close()

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
	fmt.Fprintf(w, `{"url": "%s/download/%s"}`, urlBasePath, id)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	id = filepath.Base(id) // Simple check against path traversal
	if len(id) != 32 {
		fmt.Printf("WARNING: invalid ID %s\n", id)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(storagePath, id)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("INFO: did not find file %s\n", filePath)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	metaFilePath := filepath.Join(storagePath, id+".meta")
	metaBytes, err := os.ReadFile(metaFilePath)
	if err != nil {
		fmt.Printf("ERROR: could not read file %s: %s\n", metaFilePath, err)
		http.Error(w, "Sorry, I did my best :(", http.StatusInternalServerError)
		return
	}

	var meta FileMeta
	err = json.Unmarshal(metaBytes, &meta)
	if err != nil {
		fmt.Printf("ERROR: could not unmarshal file %s: %s\n", metaFilePath, err)
		http.Error(w, "Sorry, I did my best :(", http.StatusInternalServerError)
		return
	}

	fmt.Printf("INFO: Serving file %s uploaded by %s to %s\n", meta.OriginalFilename, meta.UploaderIP, r.RemoteAddr)

	w.Header().Set("Content-Disposition", "attachment; filename=\""+meta.OriginalFilename+"\"")
	http.ServeFile(w, r, filePath)
}
