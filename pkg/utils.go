package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

// ExtractName extracts the name from the path
func ExtractName(path string) string {
	return filepath.Base(path)
}

// GetSourceTree returns a the source tree
func GetSourceTree(path string) *SourceTree {
	// Check if the path is a directory or not and if it is a supported file
	fileInfo, err := os.Stat(path)
	if err != nil || (!fileInfo.IsDir() && !IsSupportedFile(path) && IsUnwantedFilesAndFolders(path)) {
		return nil
	}

	// If the path is not a directory, return the source node
	if !fileInfo.IsDir() {
		return &SourceTree{
			Root: &SourceNode{
				Name: ExtractName(path),
				Path: path,
			},
			Nodes: nil,
		}
	}

	// If the path is a directory, create a source tree
	var sourceTree SourceTree = SourceTree{
		Root: &SourceNode{
			Name: ExtractName(path),
			Path: path,
		},
		Nodes: []*SourceTree{},
	}

	// Get all the directories and files in the path
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	for _, file := range files {
		// Get the source tree of the file
		sourceTree.Nodes = append(sourceTree.Nodes, GetSourceTree(filepath.Join(path, file.Name())))
	}

	return &sourceTree
}

// GetFileContent returns the file content
func GetFileContent(path string) ([]byte, error) {
	// Read the file content
	return os.ReadFile(path)
}

// SaveFileContent saves the file content to the output path, in append mode
func SaveFileContent(outputPath string, fileName string, relPath string, data []byte) error {
	// Write the file content to the output file in append mode
	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file")
	}
	defer file.Close()

	if _, err = file.WriteString(fmt.Sprintf("Name: %s\nPath: %s\n```\n", fileName, relPath)); err != nil {
		return fmt.Errorf("failed to write output file")
	}

	if _, err = file.Write(data); err != nil {
		return fmt.Errorf("failed to write output file")
	}

	if _, err = file.WriteString("\n```\n\n"); err != nil {
		return fmt.Errorf("failed to write output file")
	}

	return nil
}
