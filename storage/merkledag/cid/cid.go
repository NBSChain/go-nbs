package cid

import (
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multihash"
)

const (
	Raw = 0x55

	DagProtobuf = 0x70
	DagCBOR     = 0x71

	GitRaw = 0x78

	EthBlock           = 0x90
	EthBlockList       = 0x91
	EthTxTrie          = 0x92
	EthTx              = 0x93
	EthTxReceiptTrie   = 0x94
	EthTxReceipt       = 0x95
	EthStateTrie       = 0x96
	EthAccountSnapshot = 0x97
	EthStorageTrie     = 0x98
	BitcoinBlock       = 0xb0
	BitcoinTx          = 0xb1
	ZcashBlock         = 0xc0
	ZcashTx            = 0xc1
	DecredBlock        = 0xe0
	DecredTx           = 0xe1
)

// Codecs maps the name of a codec to its type
var Codecs = map[string]uint64{
	"v0":                   DagProtobuf,
	"raw":                  Raw,
	"protobuf":             DagProtobuf,
	"cbor":                 DagCBOR,
	"git-raw":              GitRaw,
	"eth-block":            EthBlock,
	"eth-block-list":       EthBlockList,
	"eth-tx-trie":          EthTxTrie,
	"eth-tx":               EthTx,
	"eth-tx-receipt-trie":  EthTxReceiptTrie,
	"eth-tx-receipt":       EthTxReceipt,
	"eth-state-trie":       EthStateTrie,
	"eth-account-snapshot": EthAccountSnapshot,
	"eth-storage-trie":     EthStorageTrie,
	"bitcoin-block":        BitcoinBlock,
	"bitcoin-tx":           BitcoinTx,
	"zcash-block":          ZcashBlock,
	"zcash-tx":             ZcashTx,
	"decred-block":         DecredBlock,
	"decred-tx":            DecredTx,
}

// CodecToStr maps the numeric codec to its name
var CodecToStr = map[uint64]string{
	Raw:                "raw",
	DagProtobuf:        "protobuf",
	DagCBOR:            "cbor",
	GitRaw:             "git-raw",
	EthBlock:           "eth-block",
	EthBlockList:       "eth-block-list",
	EthTxTrie:          "eth-tx-trie",
	EthTx:              "eth-tx",
	EthTxReceiptTrie:   "eth-tx-receipt-trie",
	EthTxReceipt:       "eth-tx-receipt",
	EthStateTrie:       "eth-state-trie",
	EthAccountSnapshot: "eth-account-snapshot",
	EthStorageTrie:     "eth-storage-trie",
	BitcoinBlock:       "bitcoin-block",
	BitcoinTx:          "bitcoin-tx",
	ZcashBlock:         "zcash-block",
	ZcashTx:            "zcash-tx",
	DecredBlock:        "decred-block",
	DecredTx:           "decred-tx",
}


var (
	// ErrVarintBuffSmall means that a buffer passed to the cid parser was not
	// long enough, or did not contain an invalid cid
	ErrVarintBuffSmall = errors.New("reading varint: buffer too small")

	// ErrVarintTooBig means that the varint in the given cid was above the
	// limit of 2^64
	ErrVarintTooBig = errors.New("reading varint: varint bigger than 64bits" +
		" and not supported")

	// ErrCidTooShort means that the cid passed to decode was not long
	// enough to be a valid Cid
	ErrCidTooShort = errors.New("cid too short")

	// ErrInvalidEncoding means that selected encoding is not supported
	// by this Cid version
	ErrInvalidEncoding = errors.New("invalid base encoding")
)

type Cid struct {
	Version  uint64
	Code     uint64
	HashType uint64
	HashLen  int
	Hash     multihash.Multihash
}

func (c *Cid) HashLength() int {
	if c.HashLen != -1 {
		return c.HashLen
	}

	if len(c.Hash) == 0{
		return -1
	}

	dec, _ := multihash.Decode(c.Hash)
	c.HashLen = dec.Length

	return c.HashLen

}
func (c *Cid) String() string {

	switch c.Version {
	case 0:
		return c.Hash.B58String()
	case 1:
		mbstr, err := multibase.Encode(multibase.Base58BTC, c.bytesV1())
		if err != nil {
			panic("should not error with hardcoded mbase: " + err.Error())
		}

		return mbstr
	default:
		panic("not possible to reach this point")
	}
}

func (c *Cid) Bytes() []byte {
	switch c.Version {
	case 0:
		return c.bytesV0()
	case 1:
		return c.bytesV1()
	default:
		panic("not possible to reach this point")
	}
}

func (c *Cid) Type() uint64 {
	return c.Code
}

func (c *Cid) bytesV0() []byte {
	return []byte(c.Hash)
}

func (c *Cid) bytesV1() []byte {

	// two 8 bytes (max) numbers plus hash
	buf := make([]byte, 2*binary.MaxVarintLen64+len(c.Hash))

	n := binary.PutUvarint(buf, c.Version)
	n += binary.PutUvarint(buf[n:], c.Code)

	cn := copy(buf[n:], c.Hash)
	if cn != len(c.Hash) {
		panic("copy hash length is inconsistent")
	}

	return buf[:n+len(c.Hash)]
}

func (c *Cid) Sum(data []byte) error {

	hash, err := multihash.Sum(data, c.HashType, c.HashLen)
	if err != nil {
		return err
	}

	c.Hash = hash

	return nil
}

func Decode(v string) (*Cid, error) {
	if len(v) < 2 {
		return nil, ErrCidTooShort
	}

	if len(v) == 46 && v[:2] == "Qm" {
		hash, err := multihash.FromB58String(v)
		if err != nil {
			return nil, err
		}

		return &Cid{
			Version: 	0,
			Code:   	DagProtobuf,
			Hash:    	hash,
			HashLen:	-1,
			HashType: 	multihash.SHA2_256,
		}, nil
	}

	_, data, err := multibase.Decode(v)
	if err != nil {
		return nil, err
	}

	return Cast(data)
}

func Cast(data []byte) (*Cid, error) {
	if len(data) == 34 && data[0] == 18 && data[1] == 32 {
		h, err := multihash.Cast(data)
		if err != nil {
			return nil, err
		}

		return &Cid{
			Code:   	DagProtobuf,
			Version: 	0,
			Hash:    	h,
			HashLen:	-1,
			HashType: 	multihash.SHA2_256,
		}, nil
	}

	vers, n := binary.Uvarint(data)
	if err := uvError(n); err != nil {
		return nil, err
	}

	if vers != 0 && vers != 1 {
		return nil, fmt.Errorf("invalid cid version number: %d", vers)
	}

	codec, cn := binary.Uvarint(data[n:])
	if err := uvError(cn); err != nil {
		return nil, err
	}

	rest := data[n+cn:]
	h, err := multihash.Cast(rest)
	if err != nil {
		return nil, err
	}

	return &Cid{
		Version: 	vers,
		Code:   	codec,
		Hash:    	h,
		HashLen:	-1,
	}, nil
}

func uvError(read int) error {
	switch {
	case read == 0:
		return ErrVarintBuffSmall
	case read < 0:
		return ErrVarintTooBig
	default:
		return nil
	}
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

