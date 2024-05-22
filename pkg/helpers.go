package pkg

import (
	"os"
	"path/filepath"
)

// Helper function for GenerateSourceTree
func (sc *SourceCollector) generateSourceTree(path string) *SourceTree {
	// Check if the path is valid or not and if it is a supported file
	fileInfo, err := os.Stat(path)
	if err != nil || sc.Validator.IsIgnored(path) {
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
		sourceTree.Nodes = append(sourceTree.Nodes, sc.generateSourceTree(filepath.Join(path, file.Name())))
	}

	return &sourceTree
}

func (sc *SourceCollector) generateSourceTreeStructure(tree *SourceTree, level int) string {
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
		subTreeStructure := sc.generateSourceTreeStructure(node, level+1)
		treeStructure += subTreeStructure
	}

	return treeStructure
}
