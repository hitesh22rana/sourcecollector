package validators

// DefaultValidator is the default validator
type DefaultValidator struct{}

// NewDefaultValidator creates a new DefaultValidator
func NewDefaultValidator() *DefaultValidator {
	return &DefaultValidator{}
}

// IsIgnored checks if the file is ignored or not
func (v *DefaultValidator) IsIgnored(path string) bool {
	// Check if the file is a sensitive file or a markdown file
	if isSensitiveFile(path) || isMarkdownFile(path) {
		return true
	}

	// Check if the file is not a directory and is not a programming file
	if !isDirectory(path) && !isProgrammingFile(path) {
		return true
	}

	// Lastly, check if the file is ignored by default
	return isUnwantedFilesAndFolders(path)
}
