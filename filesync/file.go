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

func (f *File) Path() string {
	return f.path
}

func (f *File) AbsolutePath() string {
	return filepath.Join(f.prefix, f.path)
}

func (f *File) Copy(destPrefix string, dryrun bool) (Target, error) {
	t := NewFile(destPrefix, f.path)

	canSkip, err := f.compare(t)
	if err != nil {
		return nil, err
	}
	if canSkip {
		return nil, nil
	}

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
	return t, nil
}

func (f *File) Delete() error {
	return os.Remove(f.AbsolutePath())
}

func (f *File) compare(t *File) (bool, error) {
	ffp, err := f.fingerprint()
	if err != nil {
		return false, err
	}
	tfp, err := t.fingerprint()
	if err != nil {
		return false, err
	}
	return ffp.match(tfp), nil
}

func (f *File) fingerprint() (*fingerprint, error) {
	exists, size, err := f.size()
	if err != nil {
		return nil, err
	}
	if !exists {
		return &fingerprint{false, 0, nil}, nil
	}
	sum, err := f.checksum()
	if err != nil {
		return nil, err
	}
	return &fingerprint{true, size, sum}, nil
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

func (f *File) size() (bool, int64, error) {
	stat, err := os.Stat(f.AbsolutePath())
	if err == nil {
		return true, stat.Size(), nil
	}
	if os.IsNotExist(err) {
		return false, 0, nil
	}
	return false, 0, err
}
