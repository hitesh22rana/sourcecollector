package sourcecollector

type Repository interface {
	GetMetadata() (Metadata, error)
	SaveTextFile(string) error
}

func NewRepository(url string) Repository {
	return &RemoteRepository{
		URL: url,
	}
}
