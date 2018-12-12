package filesync

import (
	"os"
)

type fingerprint struct {
	existence bool
	stat      os.FileInfo
	checksum  []byte
}

func newFingerprint(t Target) (*fingerprint, error) {
	stat, exists, err := stat(t)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &fingerprint{false, nil, nil}, nil
	}
	sum, err := t.checksum()
	if err != nil {
		return nil, err
	}
	return &fingerprint{true, stat, sum}, nil
}

func (fp *fingerprint) match(other *fingerprint) bool {
	return fp.existence == other.existence && fp.matchStat(other) && fp.matchChecksum(other)
}

func (fp *fingerprint) matchStat(other *fingerprint) bool {
	return fp.stat.IsDir() == other.stat.IsDir() && fp.stat.Mode() == other.stat.Mode() && fp.stat.Size() == other.stat.Size()
}

func (fp *fingerprint) matchChecksum(other *fingerprint) bool {
	if len(fp.checksum) != len(other.checksum) {
		return false
	}
	for idx, fc := range fp.checksum {
		if fc != other.checksum[idx] {
			return false
		}
	}
	return true
}
