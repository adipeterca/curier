package main

import "time"

type FileMeta struct {
	OriginalFilename string    `json:"original_filename"`
	UploadedAt       time.Time `json:"uploaded_at"`
	UploaderIP       string    `json:"uploader_ip"`
}

type ShareData struct {
	FileMeta
	ID string
}

type Config struct {
	MaxFileSize           int64    `json:"max_file_size"`
	AllowedFileExtensions []string `json:"allowed_extensions"`
}
