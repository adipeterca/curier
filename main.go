package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {

	log.Printf("Starting curier %s\n", version)

	parseEnvVars()
	parseFS()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", rootHandler)
	mux.HandleFunc("GET /config", configHandler)
	mux.HandleFunc("GET /download/{id}", downloadHandler)
	mux.HandleFunc("GET /share/{id}", shareHandler)
	mux.HandleFunc("GET /static/style.css", cssHandler)
	mux.HandleFunc("POST /upload", uploadHandler)

	var listenAddress = fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting and listening on http://%s ...\n", listenAddress)

	startCleanup()

	err := http.ListenAndServe(listenAddress, mux)
	if err != nil {
		log.Printf("ERROR: server failed at startup: %s\n", err)
	}
}

func parseEnvVars() {
	var err error

	if envVar := os.Getenv("CURIER_STORAGE_PATH"); envVar != "" {
		storagePath = envVar
	}

	if envVar := os.Getenv("CURIER_HOST"); envVar != "" {
		host = envVar
	}

	if envVar := os.Getenv("CURIER_PORT"); envVar != "" {
		port = envVar
	}

	if envVar := os.Getenv("CURIER_FILE_RETENTION_TIME"); envVar != "" {
		fileRetentionTime, err = strconv.ParseInt(envVar, 10, 64)
		if err != nil {
			log.Printf("CRITICAL: failed to parse fileRetentionTime, reason: %s\n", err)
			os.Exit(1)
		}
		if fileRetentionTime < 1 {
			log.Printf("WARNING: fileRetentionTime needs to be at least 1 hour - parsed value is %d\n", fileRetentionTime)
			log.Printf("WARNING: fileRetentionTime set to 1 hour\n")
			fileRetentionTime = 1
		}
	}

	if envVar := os.Getenv("CURIER_MAX_FILE_SIZE"); envVar != "" {
		maxFileSize, err = strconv.ParseInt(envVar, 10, 64)
		if err != nil {
			log.Printf("CRITICAL: failed to parse maxFileSize, reason: %s\n", err)
			os.Exit(1)
		}
	}

	if envVar := os.Getenv("CURIER_ALLOWED_FILE_EXTENSIONS"); envVar != "" {
		allowedFileExtensions = map[string]bool{}
		for _, ext := range strings.Split(envVar, ";") {
			ext = strings.TrimSpace(ext)
			if ext != "" {
				allowedFileExtensions["."+ext] = true
			}
		}

		if len(allowedFileExtensions) == 0 {
			log.Printf("CRITICAL: parsing allowedFileExtensions did not work. Exiting...\n")
			os.Exit(1)
		}
	}

	fullConfig := "\n\n  --- Environment variables ---\n"
	fullConfig += fmt.Sprintf("storagePath : %s\n", storagePath)
	fullConfig += fmt.Sprintf("host : %s\n", host)
	fullConfig += fmt.Sprintf("port : %s\n", port)
	fullConfig += fmt.Sprintf("fileRetentionTime : %d\n", fileRetentionTime)
	fullConfig += fmt.Sprintf("maxFileSize : %d bytes\n", maxFileSize)
	exts := ""
	for ext := range allowedFileExtensions {
		exts += fmt.Sprintf("\t\t\t%s\n", ext)
	}
	fullConfig += fmt.Sprintf("allowedFileExtensions:\n%s", exts)

	log.Println(fullConfig)
}

func parseFS() {
	var err error
	shareTemplate, err = template.ParseFS(templateFiles, "templates/share.html")
	if err != nil {
		log.Printf("CRITICAL: could not parse share template: %s\n", err)
		os.Exit(1)
	}
}
