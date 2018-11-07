package idService

import (
	"crypto/rsa"
	"crypto/x509"
	"github.com/multiformats/go-multihash"
)

func IDFromPublicKey(pub *rsa.PublicKey) ID {

	pubData := x509.MarshalPKCS1PublicKey(pub)

	hash, _ := multihash.Sum(pubData, multihash.SHA2_256, -1)

	return ID(hash)
}
