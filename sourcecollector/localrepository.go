package sourcecollector

func (rr *RemoteRepository) GetMetadata() (Metadata, error) {
	return Metadata{}, nil
}

func (rr *RemoteRepository) SaveTextFile(path string) error {
	return nil
}
