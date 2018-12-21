package filesync

import (
	"crypto/md5"
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

func (f *File) Prefix() string {
	return f.prefix
}

func (f *File) Path() string {
	return f.path
}

func (f *File) AbsolutePath() string {
	return filepath.Join(f.prefix, f.path)
}

func (f *File) createCopy(destPrefix string, dryrun bool) (Target, error) {
	t := NewFile(destPrefix, f.path)

	if dryrun {
		return t, nil
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

	stat, err := os.Stat(f.AbsolutePath())
	if err != nil {
		return nil, err
	}
	return t, os.Chmod(t.AbsolutePath(), stat.Mode())
}

func (f *File) Delete() error {
	return os.Remove(f.AbsolutePath())
}

func (f *File) checksum() ([]byte, error) {
	d, err := os.Open(f.AbsolutePath())
	if err != nil {
		return nil, err
	}
	h := md5.New()
	if _, err := io.Copy(h, d); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
