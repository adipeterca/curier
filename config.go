package main

// All variables can be overwritten by using environment variables.
// All env vars need to start with `CURIER_` followed by the variable name in uppercase, each word separated with an underscore.
//
// Example:
// storagePath -> CURIER_STORAGE_PATH

// Where to save the uploaded files
var storagePath = "/var/lib/curier/uploads/"

// URL base path that will prefix all download links
var urlBasePath = "http://localhost"

// Network address to bind to - default 127.0.0.1
var host = "127.0.0.1"

// Port to listen on - default 8080
var port = "8080"
