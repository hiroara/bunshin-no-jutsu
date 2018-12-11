package filesync

import (
	"os"
)

type Target interface {
	Path() string
}

func NewTarget(path string) (Target, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return NewDirectory(path), nil
	} else {
		return NewFile(path), nil
	}
}
