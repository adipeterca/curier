package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func startCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			cleanup()
		}
	}()
}

// Will it even panic? If so, it will crash the whole app - need to keep an eye on this
func cleanup() {

	log.Printf("INFO: Started a routine cleanup for expired files at %s\n", time.Now())
	entries, err := os.ReadDir(storagePath)
	if err != nil {
		log.Printf("ERROR: could not read storage directory: %s\n", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".meta" {
			continue
		}

		id := entry.Name()
		filePath := filepath.Join(storagePath, id)

		if isExpired(id) {
			err = os.Remove(filePath)
			if err != nil {
				log.Printf("ERROR: could not remove %s, reason: %s\n", filePath, err)
			}

			err = os.Remove(filePath + ".meta")
			if err != nil {
				log.Printf("ERROR: could not remove %s, reason: %s\n", filePath+".meta", err)
			}
		}
	}

	log.Printf("INFO: Finished a routine cleanup at %s\n", time.Now())
}

func isExpired(id string) bool {

	metaFile, err := readMeta(id)
	if err != nil {
		log.Printf("ERROR: could not delete meta file for %s during cleanup, reason: %s\n", id, err)
		// Return 'false' as a "skip this file"
		return false
	}

	retention := time.Duration(fileRetentionTime) * time.Hour
	return time.Since(metaFile.UploadedAt) > retention
}
