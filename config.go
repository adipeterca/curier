package main

// All variables can be overwritten by using environment variables.
// All env vars need to start with `CURIER_` followed by the variable name in uppercase, each word separated with an underscore.
//
// Example:
// storagePath -> CURIER_STORAGE_PATH

// Where to save the uploaded files
var storagePath = "/var/lib/curier/uploads/"

// URL base path that will prefix all download links
var urlBasePath = "https://myapp.example.com"
