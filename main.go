package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {

	parseEnvVars()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", rootHandler)
	mux.HandleFunc("GET /config", configHandler)
	mux.HandleFunc("GET /download/{id}", downloadHandler)
	mux.HandleFunc("GET /share/{id}", shareHandler)
	mux.HandleFunc("GET /static/style.css", cssHandler)
	mux.HandleFunc("POST /upload", uploadHandler)

	var listenAddress = fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Starting and listening on http://%s ...\n", listenAddress)

	err := http.ListenAndServe(listenAddress, mux)
	if err != nil {
		fmt.Printf("ERROR: server failed at startup: %s\n", err)
	}
}

func parseEnvVars() {
	if envVar := os.Getenv("CURIER_STORAGE_PATH"); envVar != "" {
		storagePath = envVar
	}

	if envVar := os.Getenv("CURIER_URL_BASE_PATH"); envVar != "" {
		urlBasePath = envVar
	}

	if envVar := os.Getenv("CURIER_HOST"); envVar != "" {
		host = envVar
	}

	if envVar := os.Getenv("CURIER_PORT"); envVar != "" {
		port = envVar
	}

	if envVar := os.Getenv("CURIER_MAX_FILE_SIZE"); envVar != "" {
		var err error
		maxFileSize, err = strconv.ParseInt(envVar, 10, 64)
		if err != nil {
			fmt.Printf("CRITICAL: failed to parse maxFileSize, reason: %s\n", err)
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
			fmt.Printf("CRITICAL: parsing allowedFileExtensions did not work. Exiting...\n")
			os.Exit(1)
		}
	}

	fmt.Printf("\n\n  --- Environment variables ---\n")
	fmt.Printf("storagePath : %s\n", storagePath)
	fmt.Printf("urlBasePath : %s\n", urlBasePath)
	fmt.Printf("host : %s\n", host)
	fmt.Printf("port : %s\n", port)
	fmt.Printf("maxFileSize : %d bytes\n", maxFileSize)
	fmt.Printf("allowedFileExtensions: ")
	for ext := range allowedFileExtensions {
		fmt.Printf("%s ", ext)
	}
	fmt.Printf("\n\n")
}
