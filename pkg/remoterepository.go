package sourcecollector

import (
	"context"

	"github.com/google/go-github/github"
)

func (rr *RemoteRepository) GetMetadata() (Metadata, error) {
	client := github.NewClient(nil)
	repo, _, err := client.Repositories.Get(context.Background(), rr.Owner, rr.Repository)
	if err != nil {
		return Metadata{}, err
	}

	return Metadata{
		Name:        *repo.Name,
		Description: *repo.Description,
		Topics:      repo.Topics,
		Stars:       *repo.StargazersCount,
		Size:        *repo.Size,
	}, nil
}

func (rr *RemoteRepository) SaveTextFile(path string) error {
	return nil
}
