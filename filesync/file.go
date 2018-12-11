package filesync

type File struct {
	path string
}

func NewFile(path string) *File {
	return &File{path}
}

func (f *File) Path() string {
	return f.path
}
