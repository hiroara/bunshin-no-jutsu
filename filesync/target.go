package filesync

import (
	"os"
	"path/filepath"
	"sort"
)

type Target interface {
	Prefix() string
	Path() string
	String() string
	AbsolutePath() string
	Delete(dryrun bool) error
	createCopy(destPrefix string, dryrun bool) (Target, error)
	checksum() ([]byte, error)
}

func NewTarget(prefix, path string) (Target, error) {
	stat, err := os.Lstat(filepath.Join(prefix, path))
	if err != nil {
		return nil, err
	}
	if stat.Mode()&os.ModeSymlink != 0 {
		return NewSymlink(prefix, path), nil
	} else if stat.IsDir() {
		return NewDirectory(prefix, path), nil
	} else {
		return NewFile(prefix, path), nil
	}
}

func Copy(src Target, destPrefix string, dryrun bool) (bool, Target, error) {
	destFp, destExists, err := newFingerprint(destPrefix, src.Path())
	if err != nil {
		return false, nil, err
	}
	if destExists {
		d, err := copyWithDiff(src, destFp, dryrun)
		return false, d, err
	} else {
		d, err := src.createCopy(destPrefix, dryrun)
		return true, d, err
	}
}

func ContainsSymlink(prefix, path string) bool {
	for _, cmp := range splitComponents(path) {
		stat, err := os.Lstat(filepath.Join(prefix, cmp))
		if err != nil {
			return false
		}
		if stat.Mode()&os.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

func copyWithDiff(src Target, destFingerprint *fingerprint, dryrun bool) (Target, error) {
	srcFp, _, err := newFingerprint(src.Prefix(), src.Path())
	if err != nil {
		return nil, err
	}
	return newDiff(srcFp, destFingerprint).copyTarget(dryrun)
}

func syncMode(src, dest Target) error {
	stat, err := os.Stat(src.AbsolutePath())
	if err != nil {
		return err
	}
	return os.Chmod(dest.AbsolutePath(), stat.Mode())
}

func stat(path string) (os.FileInfo, bool, error) {
	stat, err := os.Lstat(path)
	if err == nil {
		return stat, true, nil
	}
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	return nil, false, err
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
