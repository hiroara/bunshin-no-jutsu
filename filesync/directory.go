package filesync

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Directory struct {
	path string
}

func NewDirectory(path string) *Directory {
	return &Directory{path}
}

func (d Directory) Sync(destDir *Directory) error {
	ts, err := d.listTargets()
	if err != nil {
		return err
	}
	for _, t := range ts {
		fmt.Println(t.Path())
	}
	return nil
}

func (d *Directory) Path() string {
	return d.path
}

func (d *Directory) listTargets() ([]Target, error) {
	fis, err := ioutil.ReadDir(d.Path())
	if err != nil {
		return nil, err
	}
	ts := make([]Target, len(fis))
	for idx, fi := range fis {
		t, err := NewTarget(filepath.Join(d.Path(), fi.Name()))
		if err != nil {
			return nil, err
		}
		ts[idx] = t
	}
	return ts, nil
}
