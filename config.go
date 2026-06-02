package main

import "embed"

//go:embed static
var staticFiles embed.FS

//go:embed templates
var templateFiles embed.FS

// All variables can be overwritten by using environment variables.
// All env vars need to start with `CURIER_` followed by the variable name in uppercase, each word separated with an underscore.
//
// Example:
// storagePath -> CURIER_STORAGE_PATH

// -- Private variables ---

// Where to save the uploaded files
var storagePath = "uploads/"

// URL base path that will prefix all download links
var urlBasePath = "http://localhost"

// Network address to bind to - default 0.0.0.0
var host = "0.0.0.0"

// Port to listen on - default 39800
var port = "39800"

// --- Public variables ---
//
// This information can be queried by a GET request to the `/config/` endpoint.

// Max accepted file size for upload - default 20 GB
var maxFileSize int64 = 20 * 1024 * 1024 * 1024

// What file types (based on extension) can be uploaded.
// Env var looks like CURIER_ALLOWED_FILE_EXTENSIONS=jpg;jpeg;md
//
// DO NOT add a '.' for each extension - it will be added automatically.
var allowedFileExtensions = map[string]bool{
	".jpg":    true,
	".jpeg":   true,
	".webm":   true,
	".mkv":    true,
	".mp4":    true,
	".mp3":    true,
	".avi":    true,
	".png":    true,
	".pdf":    true,
	".zip":    true,
	".rar":    true,
	".tar.gz": true,
	".txt":    true,
	".md":     true,
}
