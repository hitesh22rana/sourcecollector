package sourcecollector

type RemoteRepository struct {
	URL string `json:"url"`
}

type LocalRepository struct {
	Path string `json:"path"`
}

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Data struct {
	Folders []Folder `json:"folders"`
	Files   []File   `json:"files"`
}

type Folder struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Files []File `json:"files"`
}

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
