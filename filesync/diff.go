package filesync

import (
	"os"
)

type diff struct {
	srcFingerprint  *fingerprint
	destFingerprint *fingerprint
	content         bool
	mode            bool
}

func newDiff(srcFingerprint, destFingerprint *fingerprint) *diff {
	return &diff{
		srcFingerprint, destFingerprint,
		matchChecksum(srcFingerprint.checksum, destFingerprint.checksum),
		srcFingerprint.stat.Mode() == destFingerprint.stat.Mode(),
	}
}

func (d *diff) copyTarget(dryrun bool) (Target, error) {
	if !d.content {
		return d.overwrite(dryrun)
	}
	if !d.mode {
		return d.syncMode(dryrun)
	}
	return nil, nil
}

func matchChecksum(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for idx, c := range left {
		if c != right[idx] {
			return false
		}
	}
	return true
}

func (d *diff) syncMode(dryrun bool) (Target, error) {
	dest := d.destFingerprint.target
	if dryrun {
		return dest, nil
	}
	err := os.Chmod(dest.AbsolutePath(), d.srcFingerprint.stat.Mode())
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func (d *diff) overwrite(dryrun bool) (Target, error) {
	dest := d.destFingerprint.target
	if dryrun {
		return dest, nil
	}
	err := dest.Delete()
	if err != nil {
		return nil, err
	}
	_, err = d.srcFingerprint.target.createCopy(dest.Prefix(), dryrun)
	if err != nil {
		return nil, err
	}
	return NewTarget(dest.Prefix(), d.srcFingerprint.target.Path())
}
