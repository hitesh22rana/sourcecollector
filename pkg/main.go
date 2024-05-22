package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hitesh22rana/sourcecollector/pkg/validators"
)

// NewSourceCollector creates a new SourceCollector
func NewSourceCollector(input string, output string) (*SourceCollector, error) {
	// Validate the input and output paths
	if !IsValidPath(input) {
		return nil, fmt.Errorf("input path is invalid")
	}

	// Validate if input file is a directory or not
	if !IsDirectory(input) {
		return nil, fmt.Errorf("input path is not a directory")
	}

	// Validate if output file is a directory or don't have .txt extension
	if !IsValidPath(filepath.Dir(output)) || filepath.Ext(output) != ".txt" {
		return nil, fmt.Errorf("output path is invalid")
	}

	// Make the output file if it does not exist
	outputFile, err := os.Create(output)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file")
	}
	defer outputFile.Close()

	// Make a new gitignore based validator
	var validator validators.Validator
	validator, err = validators.NewGitIgnoreBasedValidator(input)

	// If the gitIgnoreBasedValidator is nil, then make a new default validator
	if err != nil {
		fmt.Println("No .gitignore file found, proceeding with default settings.")
		validator = validators.NewDefaultValidator()
	}

	return &SourceCollector{
		Input:     input,
		Output:    output,
		BasePath:  filepath.Dir(input),
		Validator: validator,
	}, nil
}

// GenerateSourceTree generates the source tree
func (sc *SourceCollector) GenerateSourceTree() (*SourceTree, error) {
	// Generate the source tree
	sourceTree := sc.generateSourceTree(sc.Input)
	if sourceTree == nil {
		return nil, fmt.Errorf("no files found")
	}

	return sourceTree, nil
}

// GenerateSourceTreeStructure generates the source tree structure in string format
func (sc *SourceCollector) GenerateSourceTreeStructure(tree *SourceTree) (string, error) {
	// Check if the tree is nil
	if tree == nil {
		return "", fmt.Errorf("failed to generate source tree structure")
	}

	// Generate the tree structure
	sourceTreeStructure := sc.generateSourceTreeStructure(tree, 0)
	if sourceTreeStructure == "" {
		return "", fmt.Errorf("failed to generate source tree structure")
	}

	return sourceTreeStructure, nil
}

// Save saves the source tree to the output path
func (sc *SourceCollector) Save(sourceTree *SourceTree, sourceTreeStructure string) error {
	// Check if the source tree is nil
	if sourceTree == nil {
		return fmt.Errorf("source tree is nil, failed to save the source tree to the output file")
	}

	// If sourceTreeStructure is not provided, then skip saving the source tree structure else, add the source code files tree structure to the output file and save it
	if sourceTreeStructure != "" {
		if err := SaveFileContent(sc.Output, []byte(fmt.Sprintf("Source code files structure\n\n%s\n\n", sourceTreeStructure))); err != nil {
			return err
		}
	}

	queue := []*SourceTree{sourceTree}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		// Check if the node is nil
		if node == nil {
			continue
		}

		for _, child := range node.Nodes {
			// Check if the child is nil
			if child == nil {
				continue
			}

			// Check if the child is a directory or a file, if it is a directory, add it to the queue and continue
			if child.Nodes != nil {
				queue = append(queue, child)
				continue
			}

			// Check if the child is the output path
			if child.Root.Path == sc.Output {
				continue
			}

			name := child.Root.Name
			data, err := GetFileContent(child.Root.Path)
			if err != nil {
				return err
			}

			// Check if the file content is empty
			if len(data) == 0 {
				continue
			}

			// Get the relative path of the file
			relPath, err := filepath.Rel(sc.BasePath, child.Root.Path)
			if err != nil {
				return err
			}

			data = append([]byte(fmt.Sprintf("Name: %s\nPath: %s\n```\n", name, relPath)), data...)
			data = append(data, []byte("\n```\n\n")...)

			// Save the file content
			if err = SaveFileContent(sc.Output, data); err != nil {
				return err
			}
		}
	}

	return nil
}
