package pkg

import (
	"github.com/hitesh22rana/sourcecollector/pkg/validators"
)

// SourceCollector is a struct that holds the input and output of the source code
type SourceCollector struct {
	// Input of the source code
	Input string

	// Output of the source code
	Output string

	// BasePath of the source code
	BasePath string

	// Validator of the source code
	Validator validators.Validator

	// Max Concurrency to be used for io operations
	MaxConcurrency int
}

// SourceTree is a struct that holds the source code tree structure
type SourceTree struct {
	// Root of the source code tree
	Root *SourceNode

	// Nodes of the source code tree
	Nodes []*SourceTree
}

// SourceNode is a struct that holds the source code node structure
type SourceNode struct {
	// Name of the source code node
	Name string

	// Path of the source code node
	Path string
}

// queueChanData is a struct that holds the name and path of the source code node
type queueChanData struct {
	name string
	path string
}
