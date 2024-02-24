package sourcecollector

type RemoteRepository struct {
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
	Token      string `json:"token"`
	Url        string `json:"url"`
}

type Metadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`
	Stars       int      `json:"stars"`
	Size        int      `json:"size"`
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
