package filesync

import (
	"os"
)

type Target interface {
	Path() string
	AbsolutePath() string
	Copy(destPrefix string) (Target, error)
}

func NewTarget(prefix, path string) (Target, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return NewDirectory(prefix, path), nil
	} else {
		return NewFile(prefix, path), nil
	}
}
