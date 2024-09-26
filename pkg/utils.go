package pkg

import (
	"os"
	"path/filepath"
)

// isValidPath checks if the path is valid or not
func isValidPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// isDirectory checks if the path is a directory or not
func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

// extractName extracts the name from the path
func extractName(path string) string {
	return filepath.Base(path)
}
