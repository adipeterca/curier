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

	if envVar := os.Getenv("CURIER_HOST"); envVar != "" {
		host = envVar
	}

	if envVar := os.Getenv("CURIER_PORT"); envVar != "" {
		port = envVar
	}

	fmt.Printf("storagePath : %s\n", storagePath)
	fmt.Printf("urlBasePath : %s\n", urlBasePath)
	fmt.Printf("host : %s\n", host)
	fmt.Printf("port : %s\n", port)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", rootHandler)
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
