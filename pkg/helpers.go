package sourcecollector

import (
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrInvalidRepositoryUrl = fmt.Errorf("invalid repository url format")
)

func validateRepositoryURL(rawUrl string) error {
	if !strings.HasPrefix(rawUrl, "https://") {
		return fmt.Errorf("invalid repository url format")
	}

	if _, err := url.Parse(rawUrl); err != nil {
		return err
	}

	return nil
}

func extractOwnerAndRepo(url string) (string, string, error) {
	parts := strings.Split(url, "/")

	partLen := len(parts)
	if partLen < 4 {
		return "", "", ErrInvalidRepositoryUrl
	}

	owner := parts[partLen-2]

	// Remove ".git" from the repository name if present
	repository := strings.TrimSuffix(parts[partLen-1], ".git")

	return owner, repository, nil
}
