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

func Copy(srcPrefix, destPrefix, path string, dryrun bool) (Target, bool, error) {
	srcFp, _, err := newFingerprint(srcPrefix, path)
	if err != nil {
		return nil, false, err
	}
	destFp, destExists, err := newFingerprint(destPrefix, path)
	if err != nil {
		return nil, false, err
	}
	if destExists {
		dest, err := newDiff(srcFp, destFp).copyTarget(dryrun)
		return dest, false, err
	}
	dest, err := srcFp.target.createCopy(destPrefix, dryrun)
	if err != nil {
		return nil, false, err
	}
	return dest, true, nil
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
