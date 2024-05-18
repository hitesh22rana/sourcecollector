package pkg

// SourceCollector is a struct that holds the input and output of the source code
type SourceCollector struct {
	// Input of the source code
	Input string

	// Output of the source code
	Output string
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
