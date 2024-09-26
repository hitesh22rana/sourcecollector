package pkg

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hitesh22rana/sourcecollector/pkg/validators"
)

// NewSourceCollector creates a new SourceCollector
func NewSourceCollector(input string, output string, fast bool) (*SourceCollector, error) {
	// Validate the input and output paths
	if !isValidPath(input) {
		return nil, ErrInvalidInputPath
	}

	// Validate if input file is a directory or not
	if !isDirectory(input) {
		return nil, ErrInavlidInputDirectory
	}

	// Validate if output file is a directory or don't have .txt extension
	if !isValidPath(filepath.Dir(output)) || filepath.Ext(output) != ".txt" {
		return nil, ErrInvalidOutputPath
	}

	// Make the output file if it does not exist
	outputFile, err := os.Create(output)
	if err != nil {
		return nil, ErrFailedToCreateFile
	}
	defer outputFile.Close()

	// Make a new gitignore based validator
	var validator validators.Validator
	validator, err = validators.NewGitIgnoreBasedValidator(input)

	// If the gitIgnoreBasedValidator is nil, then make a new default validator
	if err != nil {
		validator = validators.NewDefaultValidator()
	}

	// If fast is true, then set maxConcurrency to max cpu cores available, else 1
	var maxConcurrency int = 1
	if fast {
		maxConcurrency = max(runtime.NumCPU(), maxConcurrency)
	}

	return &SourceCollector{
		Input:          input,
		Output:         output,
		BasePath:       filepath.Dir(input),
		Validator:      validator,
		MaxConcurrency: maxConcurrency,
	}, nil
}

// GenerateSourceTree generates the source tree
func (sc *SourceCollector) GenerateSourceTree() (*SourceTree, error) {
	// Generate the source tree
	sourceTree := sc.generateSourceTree(sc.Input)
	if sourceTree == nil {
		return nil, ErrSourceTreeGeneration
	}

	return sourceTree, nil
}

// GenerateSourceTreeStructure generates the source tree structure in string format
func (sc *SourceCollector) GenerateSourceTreeStructure(sourceTree *SourceTree) (string, error) {
	// Check if the sourceTree is nil
	if sourceTree == nil {
		return "", ErrSourceTreeStructure
	}

	// Generate the tree structure
	sourceTreeStructure := sc.generateSourceTreeStructure(sourceTree, 0)
	if sourceTreeStructure == "" {
		return "", ErrSourceTreeStructure
	}

	return sourceTreeStructure, nil
}

// SaveSourceCode saves the source tree to the output path
func (sc *SourceCollector) SaveSourceCode(sourceTree *SourceTree, sourceTreeStructure string) error {
	// Check if the source tree is nil
	if sourceTree == nil {
		return ErrSaveSourceTree
	}

	// Open the output file in append mode
	file, err := os.OpenFile(sc.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return ErrOpenOutputFile
	}
	defer file.Close()

	// If sourceTreeStructure is not provided, then skip saving the source tree structure else, add the source code files tree structure to the output file and save it
	if sourceTreeStructure != "" {
		if _, err := file.Write([]byte("Source code files structure\n\n" + sourceTreeStructure + "\n\n")); err != nil {
			return ErrWriteOutputFile
		}
	}

	// Make a data channel to save the source code files
	dataChan := make(chan string)

	// Done channel to wait for the goroutine to finish
	done := make(chan bool)

	// Save the source code files, pick the data from the data channel and save it to the output file
	go func(dataChan chan string) {
		// Wait channel for the goroutine to finish of size equal to the sc.MaxConcurrency if order is not required else fallback to 1
		waitChan := make(chan bool, sc.MaxConcurrency)

		// Make multiple goroutines to write the data concurrently
		for i := 0; i < sc.MaxConcurrency; i++ {
			go func(dataChan chan string) {
				for data := range dataChan {
					if _, err := file.WriteString(data); err != nil {
						log.Println("failed to write output file", err)
					}
				}

				// Signal the wait channel
				waitChan <- true
			}(dataChan)
		}

		// Wait for the goroutines to finish
		for i := 0; i < sc.MaxConcurrency; i++ {
			<-waitChan
		}

		// Close the waitChan channel
		close(waitChan)

		// Signal the done channel
		done <- true
	}(dataChan)

	// Make a queue channel which takes the SourceNode as input and send the source code data to the data channel
	queueChan := make(chan SourceNode)

	// Process the source code files from the queue channel and send the data to the data channel
	go func(queueChan chan SourceNode, dataChan chan string) {
		// Wait channel for the goroutine to finish of size equal to the sc.MaxConcurrency
		waitChan := make(chan bool, sc.MaxConcurrency)

		// Make multiple goroutines to read and process the source code files concurrently
		for i := 0; i < sc.MaxConcurrency; i++ {
			go func(queueChan chan SourceNode) {
				for queueData := range queueChan {
					name := queueData.Name
					path := queueData.Path

					// Use bufio for efficient file reading
					file, err := os.Open(path)
					if err != nil {
						log.Fatalln("failed to open file", err)
					}

					// Get the relative path of the file
					relPath, _ := filepath.Rel(sc.BasePath, path)

					var sb strings.Builder

					sb.WriteString("Name: " + name + "\nPath: " + relPath + "\n```\n")
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						sb.Write(scanner.Bytes())
						sb.WriteString("\n")
					}
					file.Close()

					if err := scanner.Err(); err != nil {
						log.Fatalln("failed to read file", err)
					}

					sb.WriteString("\n```\n\n")

					// Add the file content to the data channel
					dataChan <- sb.String()
				}

				// Signal the wait channel
				waitChan <- true
			}(queueChan)
		}

		// Wait for the goroutines to finish
		for i := 0; i < sc.MaxConcurrency; i++ {
			<-waitChan
		}

		// Close the waitChan channel
		close(waitChan)

		// Close the data channel
		close(dataChan)
	}(queueChan, dataChan)

	queue := []*SourceTree{sourceTree}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		// Check if the node is nil
		if node == nil {
			continue
		}

		for _, node := range node.Nodes {
			// Check if the node is nil or the node is the output path
			if node == nil || node.Root.Path == sc.Output {
				continue
			}

			// Check if the node is a directory or a file, if it is a directory, add it to the queue and continue
			if node.Nodes != nil {
				queue = append(queue, node)
				continue
			}

			// Add the file path to the queue channel
			queueChan <- SourceNode{
				Name: node.Root.Name,
				Path: node.Root.Path,
			}
		}
	}

	// Close the queue channel
	close(queueChan)

	// Wait for the goroutine to finish
	<-done

	// Close the done channel
	close(done)

	return nil
}
