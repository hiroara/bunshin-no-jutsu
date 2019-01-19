package filesync

import (
	"crypto/md5"
	"os"
	"path/filepath"
)

type Symlink struct {
	prefix string
	path   string
}

func NewSymlink(prefix, path string) *Symlink {
	return &Symlink{prefix, path}
}

func (s *Symlink) Prefix() string {
	return s.prefix
}

func (s *Symlink) String() string {
	return filepath.Join(s.prefix, s.path)
}

func (s *Symlink) Path() string {
	return s.path
}

func (s *Symlink) AbsolutePath() string {
	return filepath.Join(s.prefix, s.path)
}

func (s *Symlink) Delete(dryrun bool) error {
	if dryrun {
		return nil
	}
	return os.Remove(s.AbsolutePath())
}

func (s *Symlink) createCopy(destPrefix string, dryrun bool) (Target, error) {
	d := NewSymlink(destPrefix, s.path)
	if dryrun {
		return d, nil
	}
	path, err := os.Readlink(s.AbsolutePath())
	if err != nil {
		return nil, err
	}
	err = os.Symlink(path, d.AbsolutePath())
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *Symlink) checksum() ([]byte, error) {
	l, err := os.Readlink(s.AbsolutePath())
	if err != nil {
		return nil, err
	}
	return md5.New().Sum([]byte(l)), nil
}
