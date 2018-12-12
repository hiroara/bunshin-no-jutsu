package filesync

import (
	"os"
	"path/filepath"
)

type Target interface {
	Path() string
	AbsolutePath() string
	Copy(destPrefix string, dryrun bool) (Target, error)
	Delete() error
	checksum() ([]byte, error)
}

func NewTarget(prefix, path string) (Target, error) {
	stat, err := os.Stat(filepath.Join(prefix, path))
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return NewDirectory(prefix, path), nil
	} else {
		return NewFile(prefix, path), nil
	}
}

func stat(t Target) (os.FileInfo, bool, error) {
	stat, err := os.Stat(t.AbsolutePath())
	if err == nil {
		return stat, true, nil
	}
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	return nil, false, err
}
