package validators

import (
	"path/filepath"

	ignore "github.com/sabhiram/go-gitignore"
)

// GitIgnoreBasedValidator is a struct that implements the Validator interface
type GitIgnoreBasedValidator struct {
	// path is the path of the directory
	path string

	// gitIgnore is a set of .gitignore files
	gitIgnore *ignore.GitIgnore
}

// NewGitIgnoreBasedValidator creates a new GitIgnoreBasedValidator
func NewGitIgnoreBasedValidator(path string) (*GitIgnoreBasedValidator, error) {
	gitIgnore, err := ignore.CompileIgnoreFile(filepath.Join(path, ".gitignore"))
	if err != nil {
		return nil, err
	}

	return &GitIgnoreBasedValidator{
		path:      path,
		gitIgnore: gitIgnore,
	}, nil
}

// IsIgnored checks if the file is ignored by .gitignore
func (v *GitIgnoreBasedValidator) IsIgnored(path string) bool {
	isIgnored := isSensitiveFile(path) || isMarkdownFile(path)
	if isIgnored {
		return true
	}

	// Check if the file name or extension is ignored by .gitignore
	if v.gitIgnore.MatchesPath(path) {
		return true
	}

	// Lastly, check if the file is ignored by default
	return isUnwantedFilesAndFolders(path)
}
