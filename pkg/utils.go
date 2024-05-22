package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

// IsValidPath checks if the path is valid or not
func IsValidPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory checks if the path is a directory or not
func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

// ExtractName extracts the name from the path
func ExtractName(path string) string {
	return filepath.Base(path)
}

// GetFileContent returns the file content
func GetFileContent(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// SaveFileContent saves the file content to the output path, in append mode
func SaveFileContent(outputPath string, data []byte) error {
	// Write the file content to the output file in append mode
	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file")
	}
	defer file.Close()

	if _, err = file.Write(data); err != nil {
		return fmt.Errorf("failed to write output file")
	}

	return nil
}
