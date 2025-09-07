package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveFile saves the uploaded file to the specified folder and returns the file path
func SaveFile(file *multipart.FileHeader, folder string) (string, error) {
	// Ensure the folder exists
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique file name (timestamp + original extension)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
	filePath := filepath.Join(folder, fileName)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file content
	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
