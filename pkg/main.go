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
func (sc *SourceCollector) GenerateSourceTreeStructure(sourceTree *SourceTree) (string, error) {
	// Check if the sourceTree is nil
	if sourceTree == nil {
		return "", fmt.Errorf("failed to generate source tree structure")
	}

	// Generate the tree structure
	sourceTreeStructure := sc.generateSourceTreeStructure(sourceTree, 0)
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

	// Open the output file in append mode
	file, err := os.OpenFile(sc.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file")
	}
	defer file.Close()

	// If sourceTreeStructure is not provided, then skip saving the source tree structure else, add the source code files tree structure to the output file and save it
	if sourceTreeStructure != "" {
		if _, err := file.Write([]byte(fmt.Sprintf("Source code files structure\n\n%s\n\n", sourceTreeStructure))); err != nil {
			return fmt.Errorf("failed to write output file")
		}
	}

	// Make a data channel to save the source code files
	dataChan := make(chan []byte)

	// Done channel to wait for the goroutine to finish
	done := make(chan bool)

	// Save the source code files, pick the data from the data channel and save it to the output file
	go func(dataChan chan []byte) {
		// defer wg.Done()
		for data := range dataChan {
			if _, err := file.Write(data); err != nil {
				fmt.Println("failed to write output file", err)
			}
		}

		// Signal the done channel
		done <- true
	}(dataChan)

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

			// Add the file content to the data channel
			dataChan <- data
		}
	}

	// Close the data channel after all data has been sent
	close(dataChan)

	// Wait for the goroutine to finish
	<-done

	// Close the done channel
	close(done)

	return nil
}
