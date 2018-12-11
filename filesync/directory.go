package filesync

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Directory struct {
	prefix string
	path   string
}

func NewDirectory(prefix, path string) *Directory {
	return &Directory{prefix, path}
}

func (d Directory) Sync(dest string) error {
	ts, err := d.listTargets()
	if err != nil {
		return err
	}
	for _, t := range ts {
		dest, err := t.Copy(dest)
		if err != nil {
			return err
		}
		fmt.Printf("%s => %s\n", t.AbsolutePath(), dest.AbsolutePath())
	}
	return nil
}

func (d *Directory) Path() string {
	return d.path + "/"
}

func (d *Directory) AbsolutePath() string {
	return filepath.Join(d.prefix, d.path) + "/"
}

func (d *Directory) Copy(destPrefix string) (Target, error) {
	dest := NewDirectory(destPrefix, d.path)
	err := dest.MkdirAll()
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func (d *Directory) MkdirAll() error {
	return os.MkdirAll(d.AbsolutePath(), 0777)
}

func (d *Directory) listTargets() ([]Target, error) {
	fis, err := ioutil.ReadDir(d.Path())
	if err != nil {
		return nil, err
	}
	if len(fis) == 0 {
		return []Target{d}, nil
	}
	targets := make([]Target, 0)
	for _, fi := range fis {
		tg, err := NewTarget(d.prefix, filepath.Join(d.Path(), fi.Name()))
		if err != nil {
			return nil, err
		}
		switch t := tg.(type) {
		case *Directory:
			ts, err := t.listTargets()
			if err != nil {
				return nil, err
			}
			targets = append(targets, ts...)
		case *File:
			targets = append(targets, t)
		}
	}
	return targets, nil
}
