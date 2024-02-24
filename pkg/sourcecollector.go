package sourcecollector

import (
	"net/url"
)

type Repository interface {
	GetMetadata() (Metadata, error)
	SaveTextFile(string) error
}

func NewRepository(rawUrl string) (Repository, error) {
	_, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	return &RemoteRepository{
		URL: rawUrl,
	}, nil
}
