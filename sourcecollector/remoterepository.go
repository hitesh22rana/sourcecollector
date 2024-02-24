package sourcecollector

func (rr *LocalRepository) GetMetadata() (Metadata, error) {
	return Metadata{}, nil
}

func (rr *LocalRepository) SaveTextFile(path string) error {
	return nil
}
