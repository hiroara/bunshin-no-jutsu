package filesync

import (
	"os"
	"path/filepath"
)

type Target interface {
	Prefix() string
	Path() string
	AbsolutePath() string
	Delete() error
	createCopy(destPrefix string, dryrun bool) (Target, error)
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

func copyWithDiff(src Target, destFingerprint *fingerprint, dryrun bool) (Target, error) {
	srcFp, _, err := newFingerprint(src.Prefix(), src.Path())
	if err != nil {
		return nil, err
	}
	return newDiff(srcFp, destFingerprint).copyTarget(dryrun)
}

func stat(path string) (os.FileInfo, bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return stat, true, nil
	}
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	return nil, false, err
}
