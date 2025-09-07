package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile saves the uploaded file and returns its URL
func SaveUploadedFile(file multipart.File, handler *multipart.FileHeader, uploadDir string) (string, error) {
	// Ensure directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	// Unique filename
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	// Return accessible path
	return "/" + filePath, nil
}
