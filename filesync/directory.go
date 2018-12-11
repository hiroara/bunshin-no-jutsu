package filesync

import (
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

func (d *Directory) Path() string {
	return d.path + "/"
}

func (d *Directory) AbsolutePath() string {
	return filepath.Join(d.prefix, d.path) + "/"
}

func (d *Directory) Copy(destPrefix string, dryrun bool) (Target, error) {
	dest := NewDirectory(destPrefix, d.path)
	exists, err := dest.exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, nil
	}

	if dryrun {
		return dest, nil
	}

	err = dest.MkdirAll()
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func (d *Directory) Delete() error {
	return os.Remove(d.AbsolutePath())
}

func (d *Directory) exists() (bool, error) {
	_, err := os.Stat(d.AbsolutePath())
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (d *Directory) MkdirAll() error {
	return os.MkdirAll(d.AbsolutePath(), 0777)
}

func (d *Directory) ListTargets() ([]Target, error) {
	exists, err := d.exists()
	if err != nil {
		return nil, err
	}
	if !exists {
		return []Target{}, nil
	}

	fis, err := ioutil.ReadDir(d.AbsolutePath())
	if err != nil {
		return nil, err
	}
	targets := make([]Target, 0)
	for _, fi := range fis {
		tg, err := NewTarget(d.prefix, filepath.Join(d.Path(), fi.Name()))
		if err != nil {
			return nil, err
		}
		switch t := tg.(type) {
		case *Directory:
			targets = append(targets, t)
			ts, err := t.ListTargets()
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
