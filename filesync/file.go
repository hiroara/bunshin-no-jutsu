package filesync

import (
	"io"
	"os"
	"path/filepath"
)

type File struct {
	prefix string
	path   string
}

func NewFile(prefix, path string) *File {
	return &File{prefix, path}
}

func (f *File) Path() string {
	return f.path
}

func (f *File) AbsolutePath() string {
	return filepath.Join(f.prefix, f.path)
}

func (f *File) Dir() *Directory {
	return NewDirectory(f.prefix, filepath.Dir(f.path))
}

func (f *File) Copy(destPrefix string) (Target, error) {
	t := NewFile(destPrefix, f.path)

	err := t.Dir().MkdirAll()
	if err != nil {
		return nil, err
	}

	src, err := os.Open(f.AbsolutePath())
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dest, err := os.Create(t.AbsolutePath())
	if err != nil {
		return nil, err
	}
	defer src.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return nil, err
	}
	return t, nil
}
