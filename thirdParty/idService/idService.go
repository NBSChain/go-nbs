package idService

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"sync"
)

type IdService interface {
	GenerateId(encodedKey string) (*Identity, error)
}

type nbsIdService struct {

}

var instance 	IdService
var once 	sync.Once
var logger	= utils.GetLogInstance()

func GetInstance() IdService {

	once.Do(func() {
		instance = newIdService()
	})

	return instance
}

func newIdService() *nbsIdService{

	service := &nbsIdService{

	}

	return service
}

func (service *nbsIdService) GenerateId(encodedKey string) (*Identity, error){

	pri, pub, err := crypto.GenerateRSAKeyPair()
	if err != nil{
		return nil , err
	}

	id := IDFromPublicKey(pub)

	identity := &Identity{
		PeerID: id.Pretty(),
	}

	privateKeyData := x509.MarshalPKCS1PrivateKey(pri)
	if encodedKey != ""{

		key, _ := hex.DecodeString(encodedKey)
		privateKeyData = crypto.EncryptAES(privateKeyData, key)

		identity.Encrypted = true
	}

	identity.PrivateKey = base64.StdEncoding.EncodeToString(privateKeyData)

	logger.Info(identity)

	return identity, nil
}
