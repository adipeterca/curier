package main

import "time"

type FileMeta struct {
	OriginalFilename string    `json:"original_filename"`
	UploadedAt       time.Time `json:"uploaded_at"`
	UploaderIP       string    `json:"uploader_ip"`
	FileSize         int64     `json:"file_size"`
}

type ShareData struct {
	FileMeta
	ID        string
	ExpiresAt time.Time
	FileSize  int64
}

type Config struct {
	MaxFileSize           int64    `json:"max_file_size"`
	AllowedFileExtensions []string `json:"allowed_extensions"`
}
