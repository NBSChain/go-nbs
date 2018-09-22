package cid

import (
	"encoding/base32"
	"fmt"
	"github.com/multiformats/go-multihash"
)


var ErrPossiblyInsecureHashFunction = fmt.Errorf("potentially insecure hash functions not allowed")
var ErrBelowMinimumHashLength = fmt.Errorf("hashes must be at %d least bytes long", minimumHashLength)

const minimumHashLength = 20

var goodset = map[uint64]bool{
	multihash.SHA2_256:     true,
	multihash.SHA2_512:     true,
	multihash.SHA3_224:     true,
	multihash.SHA3_256:     true,
	multihash.SHA3_384:     true,
	multihash.SHA3_512:     true,
	multihash.SHAKE_256:    true,
	multihash.DBL_SHA2_256: true,
	multihash.KECCAK_224:   true,
	multihash.KECCAK_256:   true,
	multihash.KECCAK_384:   true,
	multihash.KECCAK_512:   true,
	multihash.ID:           true,

	multihash.SHA1: true, // not really secure but still useful
}

func IsGoodHash(code uint64) bool {
	good, found := goodset[code]
	if good {
		return true
	}

	if !found {
		if code >= multihash.BLAKE2B_MIN+19 && code <= multihash.BLAKE2B_MAX {
			return true
		}
		if code >= multihash.BLAKE2S_MIN+19 && code <= multihash.BLAKE2S_MAX {
			return true
		}
	}

	return false
}

func ValidateCid(c *Cid) error {

	if !IsGoodHash(c.HashType) {
		return ErrPossiblyInsecureHashFunction
	}

	if c.HashType != multihash.ID && c.HashType < minimumHashLength {
		return ErrBelowMinimumHashLength
	}

	return nil
}

func NewKeyFromBinary(rawKey []byte) string {
	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	buf := make([]byte, 1 + encoder.EncodedLen(len(rawKey)))
	buf[0] = '/'
	encoder.Encode(buf[1:], rawKey)
	return string(buf)
}

func BinaryFromDsKey(k string) ([]byte, error) {
	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	return encoder.DecodeString(k[1:])
}

func CidToDsKey(k *Cid) string {
	return NewKeyFromBinary(k.Bytes())
}

func DsKeyToCid(dsKey string) (*Cid, error) {
	kb, err := BinaryFromDsKey(dsKey)
	if err != nil {
		return nil, err
	}
	return Cast(kb)
}