package account

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/NBSChain/go-nbs/thirdParty/idService"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"sync"
	"time"
)

type Account interface {
	UnlockAccount(password string) error
}

type nbsAccount struct {
	accountId	*idService.Identity
	privateKey	*rsa.PrivateKey
	cleanKeyTimer	chan struct{}
}

var instance 	*nbsAccount
var once 	sync.Once
var logger	= utils.GetLogInstance()

func GetAccountInstance() Account{

	once.Do(func() {
		instance = loadAccount()
		go instance.MonitorPrivateKey()
	})

	return instance
}

func loadAccount() *nbsAccount{
	obj := &nbsAccount{
		cleanKeyTimer:make(chan struct{}),
	}
	return obj
}

func (account *nbsAccount) UnlockAccount(encodedKey string) (err error){

	if account.accountId == nil{
		return fmt.Errorf("the account is not initialized yet")
	}

	privateKeyData, err := base64.StdEncoding.DecodeString(account.accountId.PrivateKey)
	if err != nil{
		return err
	}

	var decryptedData []byte

	if !account.accountId.Encrypted {
		decryptedData = privateKeyData
	}else{
		key, _ := hex.DecodeString(encodedKey)

		decryptedData = crypto.DecryptAES(privateKeyData, key)
	}


	account.privateKey, err = x509.ParsePKCS1PrivateKey(decryptedData)

	account.cleanKeyTimer <- struct{}{}

	return nil
}

func (account *nbsAccount) MonitorPrivateKey()  {

	for{
		select {
		case <-account.cleanKeyTimer:
			time.AfterFunc( 5 * time.Minute, func(){
				account.privateKey = nil
			})
		}
	}
}