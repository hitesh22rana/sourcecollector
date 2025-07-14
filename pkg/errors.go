package pkg

import "errors"

var (
	ErrInvalidInputPath      = errors.New("input path is invalid")
	ErrInvalidInputDirectory = errors.New("input path is not a valid directory")
	ErrInvalidOutputPath     = errors.New("output path is invalid")
	ErrFailedToCreateFile    = errors.New("failed to create output file")
	ErrSourceTreeGeneration  = errors.New("failed to generate source tree")
	ErrSourceTreeStructure   = errors.New("failed to generate source tree structure")
	ErrSaveSourceTree        = errors.New("failed to save source tree to file")
	ErrOpenOutputFile        = errors.New("failed to open output file")
	ErrWriteOutputFile       = errors.New("failed to write to output file")
)
