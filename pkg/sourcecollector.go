package sourcecollector

type Repository interface {
	GetMetadata() (Metadata, error)
	SaveTextFile(string) error
}

func NewRepository(url string) (Repository, error) {
	if err := validateRepositoryURL(url); err != nil {
		return nil, err
	}

	owner, repository, err := extractOwnerAndRepo(url)
	if err != nil {
		return nil, err
	}

	return &RemoteRepository{
		Owner:      owner,
		Repository: repository,
		Url:        url,
	}, nil
}
