package filesync

import (
	"os"
	"path/filepath"
)

type fingerprint struct {
	target   Target
	stat     os.FileInfo
	checksum []byte
	symlink  bool
}

func newFingerprint(prefix, path string) (*fingerprint, bool, error) {
	absPath := filepath.Join(prefix, path)
	stat, exists, err := stat(absPath)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return &fingerprint{nil, nil, nil, false}, false, nil
	}
	t, err := NewTarget(prefix, path)
	if err != nil {
		return nil, false, err
	}
	sum, err := t.checksum()
	if err != nil {
		return nil, false, err
	}
	return &fingerprint{t, stat, sum, stat.Mode()&os.ModeSymlink != 0}, true, nil
}
