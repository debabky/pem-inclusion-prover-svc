package models

import (
	"hash"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

type PoseidonHash struct {
	inner hash.Hash
}

func NewPoseidonHash(frameSize int) (*PoseidonHash, error) {
	p, err := poseidon.New(frameSize)

	if err != nil {
		return nil, err
	}

	return &PoseidonHash{p}, nil
}

func (ph *PoseidonHash) Hash(data ...[]byte) []byte {
	var hash []byte
	if len(data) == 1 {
		hash = poseidon.Sum(data[0])
	} else {
		concatDataLen := 0
		for _, d := range data {
			concatDataLen += len(d)
		}
		concatData := make([]byte, concatDataLen)
		curOffset := 0
		for _, d := range data {
			copy(concatData[curOffset:], d)
			curOffset += len(d)
		}
		hash = poseidon.Sum(concatData)
	}

	return hash
}

func (ph *PoseidonHash) HashLength() int {
	return ph.inner.Size()
}

func (ph *PoseidonHash) HashName() string {
	return "poseidon"
}
