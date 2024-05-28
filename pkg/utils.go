package pkg

import (
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
