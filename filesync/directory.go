package filesync

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type Directory struct {
	prefix string
	path   string
}

func NewDirectory(prefix, path string) *Directory {
	return &Directory{prefix, path}
}

func (d *Directory) Prefix() string {
	return d.prefix
}

func (d *Directory) Path() string {
	return d.path + "/"
}

func (d *Directory) AbsolutePath() string {
	return filepath.Join(d.prefix, d.path) + "/"
}

func (d *Directory) createCopy(destPrefix string, dryrun bool) (Target, error) {
	dest := NewDirectory(destPrefix, d.path)
	_, exists, err := stat(dest.AbsolutePath())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, nil
	}

	if dryrun {
		return dest, nil
	}

	err = d.MkdirAll(destPrefix, dryrun)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func (d *Directory) Delete() error {
	return os.Remove(d.AbsolutePath())
}

func (d *Directory) MkdirAll(destPrefix string, dryrun bool) error {
	if dryrun {
		return nil
	}
	for _, part := range splitComponents(d.path) {
		path := filepath.Join(destPrefix, part)
		srcStat, err := os.Stat(filepath.Join(d.prefix, part))
		if err != nil {
			return err
		}
		destStat, err := os.Stat(path)
		doesNotExist := os.IsNotExist(err)
		if doesNotExist {
			err = os.Mkdir(path, 0777)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		if doesNotExist || destStat.Mode() != srcStat.Mode() {
			os.Chmod(path, srcStat.Mode())
		}
	}
	return nil
}

func (d *Directory) checksum() ([]byte, error) {
	return []byte{}, nil
}

func splitComponents(path string) []string {
	targets := []string{path}
	for {
		last := targets[len(targets)-1]
		parent := filepath.Dir(last)
		if parent == last {
			sort.Sort(sort.Reverse(sort.StringSlice(targets)))
			return targets[:len(targets)-1]
		}
		targets = append(targets, parent)
	}
}

func (d *Directory) ListTargets() ([]Target, error) {
	_, exists, err := stat(d.AbsolutePath())
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
