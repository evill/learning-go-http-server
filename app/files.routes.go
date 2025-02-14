package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func isFileExists(fileName string, filesDirectory string) bool {
	fullPath := path.Join(filesDirectory, fileName)
	_, err := os.Stat(fullPath)
	return !errors.Is(err, os.ErrNotExist)
}

func validateFileName(fileName string) error {
	// Check for path traversal attempts
	if strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		return fmt.Errorf("file name contains invalid path separators")
	}

	// Check if filename is actually a path
	cleanPath := filepath.Clean(fileName)
	if cleanPath != fileName {
		return fmt.Errorf("file name contains path elements")
	}

	// Additional security checks
	if fileName == "." || fileName == ".." || strings.HasPrefix(fileName, ".") {
		return fmt.Errorf("invalid file name")
	}

	return nil
}

func getFileRoute(request *HttpRequest, response *HttpResponse) {
	fileName, _ := strings.CutPrefix(request.path, "/files/")

	if fileName == "" {
		response.Status404().Text("Name of file is not passed in URL")
		return
	}

	if err := validateFileName(fileName); err != nil {
		response.Status400().Text(fmt.Sprintf("Invalid file name: %v", err))
		return
	}

	filesDirectory := *request.server.config.filesDirectory
	files, err := os.ReadDir(filesDirectory)
	if err != nil {
		log.Print(err)
		response.Status500().Text("File server feature is not available!")
		return
	}

	var targetFile fs.DirEntry
	for _, file := range files {
		if file.Name() == fileName && !file.IsDir() {
			targetFile = file
			break
		}
	}

	if targetFile == nil {
		log.Printf("Requested file '%s' doesn't exists in folder '$s'", fileName, filesDirectory)
		response.Status404().Text("Requested file not found")
		return
	}

	fullFilePath := path.Join(filesDirectory, fileName)

	response.SetHeader("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	response.Status200().LocalFile(fullFilePath)
}

func postFileRoute(request *HttpRequest, response *HttpResponse) {
	fileName, _ := strings.CutPrefix(request.path, "/files/")

	if err := validateFileName(fileName); err != nil {
		response.Status400().Text(fmt.Sprintf("Invalid file name: %v", err))
		return
	}

	filesDirectory := *request.server.config.filesDirectory
	fullPath := path.Join(filesDirectory, fileName)

	// Verify the joined path is still within target directory
	if !strings.HasPrefix(fullPath, filesDirectory) {
		response.Status400().Text("Invalid file path")
		return
	}

	if isFileExists(fileName, filesDirectory) {
		log.Printf("Conflict: file %s already exists.", fileName)
		response.Status409().Send()
		return
	}

	// Write content to file
	err := os.WriteFile(fullPath, []byte(request.body), 0644)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		response.Status500().Text("Failed to save file")
		return
	}

	response.Status(201, "Created").Send()
}
