package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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

func getRemoteAddress(h *http.Header) string {

	// Cloudflare proxying traffic
	ip := h.Get("Cf-Connecting-Ip")
	if ip != "" {
		return ip
	}

	// Traffic coming from on-host reverse proxies like Caddy
	ip = h.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	return ""
}

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
