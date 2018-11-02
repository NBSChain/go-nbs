package idService

import (
	"crypto/rsa"
	"crypto/x509"
	"gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
)

func IDFromPublicKey(pub *rsa.PublicKey) ID{

	pubData := x509.MarshalPKCS1PublicKey(pub)

	hash, _ := multihash.Sum(pubData, multihash.ID, -1)

	return ID(hash)
}