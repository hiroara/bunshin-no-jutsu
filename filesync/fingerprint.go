package filesync

type fingerprint struct {
	existence bool
	size      int64
	checksum  []byte
}

func (fp *fingerprint) match(other *fingerprint) bool {
	return fp.existence == other.existence && fp.size == other.size && fp.matchChecksum(other)
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
