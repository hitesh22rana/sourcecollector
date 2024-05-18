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

// GenerateSourceTree returns a the source tree
func GenerateSourceTree(path string) *SourceTree {
	// Check if the path is valid or not and if it is a supported file
	fileInfo, err := os.Stat(path)
	if err != nil || (!IsSupportedFile(path) && IsUnwantedFilesAndFolders(path)) {
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
		sourceTree.Nodes = append(sourceTree.Nodes, GenerateSourceTree(filepath.Join(path, file.Name())))
	}

	return &sourceTree
}

// GetSourceTreeStructure generates the source tree in a tree structure (string) for better understanding
func GetSourceTreeStructure(tree *SourceTree, level int) (string, error) {
	// Check if the tree is nil
	if tree == nil {
		return "", fmt.Errorf("failed to generate source tree")
	}

	// Generate the tree structure
	var treeStructure string
	for i := 0; i < level; i++ {
		treeStructure += "|  "
	}

	treeStructure += "|--" + tree.Root.Name + "\n"

	for _, node := range tree.Nodes {
		// Check if the node is nil
		if node == nil {
			continue
		}

		// Generate the source tree structure
		subTreeStructure, err := GetSourceTreeStructure(node, level+1)
		if err != nil {
			return "", err
		}

		treeStructure += subTreeStructure
	}

	return treeStructure, nil
}

// GetFileContent returns the file content
func GetFileContent(path string) ([]byte, error) {
	// Read the file content
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
