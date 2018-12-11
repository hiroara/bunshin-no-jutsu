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
