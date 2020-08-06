package gbloom

// Hasher is an interface provided to BloomFilter to generate hashes for the
// payload.
type Hasher interface {
	BuildHashes(v interface{}) ([][]byte, error)
}

// NewHasher builds a Hasher implementation with the given callback function.
func NewHasher(f func(v interface{}) ([][]byte, error)) Hasher {
	return hasher{f}
}

type hasher struct {
	buildHashes func(v interface{}) ([][]byte, error)
}

var _ Hasher = (*hasher)(nil)

func (h hasher) BuildHashes(v interface{}) ([][]byte, error) {
	return h.buildHashes(v)
}
