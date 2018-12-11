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
	targets := make([]Target, 0)
	for _, fi := range fis {
		tg, err := NewTarget(filepath.Join(d.Path(), fi.Name()))
		if err != nil {
			return nil, err
		}
		switch t := tg.(type) {
		case *Directory:
			ts, err := t.listTargets()
			if err != nil {
				return nil, err
			}
			if len(ts) > 0 {
				targets = append(targets, ts...)
			} else {
				targets = append(targets, t)
			}
		case *File:
			targets = append(targets, t)
		}
	}
	return targets, nil
}
