package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	if envVar := os.Getenv("CURIER_STORAGE_PATH"); envVar != "" {
		storagePath = envVar
	}

	if envVar := os.Getenv("CURIER_URL_BASE_PATH"); envVar != "" {
		urlBasePath = envVar
	}

	fmt.Printf("storagePath : %s\n", storagePath)
	fmt.Printf("urlBasePath : %s\n", urlBasePath)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /upload", uploadHandler)
	mux.HandleFunc("GET /download/{id}", downloadHandler)

	fmt.Printf("Starting and listening on http://127.0.0.1:8080 ...\n")

	err := http.ListenAndServe("127.0.0.1:8080", mux)
	if err != nil {
		fmt.Printf("ERROR: server failed at startup: %s\n", err)
	}
}
