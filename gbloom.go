package gbloom

import (
	"errors"
	"math/big"
)

const (
	defaultMaxBits = 100_000_007
)

// BloomFilter is a bloom filter implementation.
type BloomFilter struct {
	value   *big.Int
	maxBits int
	hasher  Hasher
}

// New initializes and returns a BloomFilter instance.
func New(maxBits int, hasher Hasher) *BloomFilter {
	return initialize(maxBits, hasher)
}

// NewFromBytes initializes a BloomFilter instance and loads the given bytes
// into the bloom filter.
func NewFromBytes(maxBits int, hasher Hasher, bytes []byte) *BloomFilter {
	ret := initialize(maxBits, hasher)
	ret.SetBytes(bytes)

	return ret
}

// NewFromInt initializes a BloomFilter instance and loads the given big.Int
// value into the bloom filter.
func NewFromInt(maxBits int, hasher Hasher, value *big.Int) *BloomFilter {
	ret := initialize(maxBits, hasher)
	ret.SetInt(value)

	return ret
}

func initialize(maxBits int, hasher Hasher) *BloomFilter {
	if maxBits == 0 {
		maxBits = defaultMaxBits
	}

	return &BloomFilter{
		value:   new(big.Int),
		maxBits: maxBits,
		hasher:  hasher,
	}
}

// Add sets the given payload in the bloom filter.
func (b *BloomFilter) Add(v interface{}) error {
	hashes, err := b.hash(v)
	if err != nil {
		return err
	}

	if len(hashes) == 1 {
		b.AddHash(hashes[0])
	} else {
		b.AddHashes(hashes...)
	}

	return nil
}

// AddHash sets a single hash in the bloom filter.
func (b *BloomFilter) AddHash(hash []byte) {
	b.value.SetBit(b.value, b.bitIndex(hash), 1)
}

// AddHashes sets the hashes in the bloom filter.
func (b *BloomFilter) AddHashes(hashes ...[]byte) {
	for _, hash := range hashes {
		b.AddHash(hash)
	}
}

// Bytes exports and returns the bytes of the bloom filter.
func (b *BloomFilter) Bytes() []byte {
	return b.value.Bytes()
}

// MaxBits returns the max allowed bits of the bloom filters.
func (b *BloomFilter) MaxBits() int {
	return b.maxBits
}

// SetBytes loads the bytes into the bloom filter.
func (b *BloomFilter) SetBytes(bytes []byte) {
	if len(bytes) == 0 {
		return
	}

	b.value.SetBytes(bytes)
}

// SetInt loads the big.Int value into the bloom filter.
func (b *BloomFilter) SetInt(value *big.Int) {
	if value == nil {
		value = new(big.Int)
	}

	b.value.Set(value)
}

// Test checks whether the payload is set by Add in the bloom filter. The
// result may include false positives but not false negatives.
func (b *BloomFilter) Test(v interface{}) (bool, error) {
	hashes, err := b.hash(v)
	if err != nil {
		return true, err
	}

	if len(hashes) == 1 {
		return b.TestHash(hashes[0]), nil
	}

	return b.TestHashes(hashes...), nil
}

// TestHash checks whether the hash is set in the bloom filter.
func (b *BloomFilter) TestHash(hash []byte) bool {
	return b.value.Bit(b.bitIndex(hash)) > 0
}

// TestHashes checks whether the hashes are set in the bloom filter.
func (b *BloomFilter) TestHashes(hashes ...[]byte) bool {
	for _, hash := range hashes {
		if !b.TestHash(hash) {
			return false
		}
	}

	return true
}

func (b *BloomFilter) bitIndex(hash []byte) (v int) {
	hash = b.truncateHash(4, hash)

	v |= int(hash[0])
	v |= int(hash[1]) << 8
	v |= int(hash[2]) << 16
	v |= int(hash[3]) << 24

	if v < 0 {
		v = -v
	}

	return v % b.maxBits
}

func (b *BloomFilter) hash(v interface{}) ([][]byte, error) {
	if b.hasher == nil {
		return nil, errors.New("hasher not set")
	}

	hashes, err := b.hasher.BuildHashes(v)
	if err != nil {
		return nil, err
	}

	if len(hashes) == 0 {
		return nil, errors.New("found empty result returned by Hasher")
	}

	return hashes, nil
}

func (b *BloomFilter) truncateHash(size int, hash []byte) []byte {
	p := len(hash) - size
	if p > 0 {
		return hash[p:]
	}

	if p < 0 {
		return append(make([]byte, -p, size), hash...)
	}

	return hash
}
