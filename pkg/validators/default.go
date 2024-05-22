package validators

// DefaultValidator is the default validator
type DefaultValidator struct{}

// NewDefaultValidator creates a new DefaultValidator
func NewDefaultValidator() *DefaultValidator {
	return &DefaultValidator{}
}

// IsIgnored checks if the file is ignored or not
func (v *DefaultValidator) IsIgnored(path string) bool {
	return isSensitiveFile(path) || isMarkdownFile(path) || !isProgrammingFile(path) || isUnwantedFilesAndFolders(path)
}
