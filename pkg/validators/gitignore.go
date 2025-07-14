package validators

import (
	"path/filepath"

	ignore "github.com/sabhiram/go-gitignore"
)

// GitIgnoreBasedValidator is a struct that implements the Validator interface
type GitIgnoreBasedValidator struct {
	// GitIgnore is used to check if the file is ignored by .gitignore
	GitIgnore *ignore.GitIgnore
}

// NewGitIgnoreBasedValidator creates a new GitIgnoreBasedValidator
func NewGitIgnoreBasedValidator(path string) (*GitIgnoreBasedValidator, error) {
	gitIgnore, err := ignore.CompileIgnoreFile(filepath.Join(path, ".gitignore"))
	if err != nil {
		return nil, err
	}

	return &GitIgnoreBasedValidator{
		GitIgnore: gitIgnore,
	}, nil
}

// IsIgnored checks if the file is ignored by .gitignore
func (v *GitIgnoreBasedValidator) IsIgnored(path string) bool {
	// Check if the file is a sensitive file
	if isSensitiveFile(path) {
		return true
	}

	// Check if the file name or extension is ignored by .gitignore
	if v.GitIgnore.MatchesPath(path) {
		return true
	}

	// Check if the file is not a directory and is not a programming file or informative file
	if !isDirectory(path) && !isProgrammingFile(path) && !isInformativeFile(path) {
		return true
	}

	// Lastly, check if the file is ignored by default
	return isUnwantedFilesAndFolders(path)
}
