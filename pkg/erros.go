package pkg

import "errors"

var (
	ErrInvalidInputPath      = errors.New("input path is invalid")
	ErrInavlidInputDirectory = errors.New("input path is not a directory")
	ErrInvalidOutputPath     = errors.New("output path is invalid")
	ErrFailedToCreateFile    = errors.New("failed to create output file")
	ErrSourceTreeGeneration  = errors.New("failed to generate source tree")
	ErrSourceTreeStructure   = errors.New("failed to generate source tree structure")
	ErrSaveSourceTree        = errors.New("failed to save source tree to file")
	ErrOpenOutputFile        = errors.New("failed to open output file")
	ErrWriteOutputFile       = errors.New("failed to write to output file")
)
