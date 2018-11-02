package account

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/NBSChain/go-nbs/storage/application/dataStore"
	"github.com/NBSChain/go-nbs/thirdParty/account/pb"
	"github.com/NBSChain/go-nbs/thirdParty/idService"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/crypto"
	"github.com/gogo/protobuf/proto"
	"sync"
	"time"
)

type Account interface {
	UnlockAccount(password string) error
	GetPeerID() string
	Reload() error
}

type nbsAccount struct {
	accountId     *idService.Identity
	privateKey    *rsa.PrivateKey
	cleanKeyTimer chan struct{}
	dataStore     dataStore.DataStore
}

var instance *nbsAccount
var once sync.Once
var logger = utils.GetLogInstance()

const ParameterKeyForAccount = "keys_for_account_local_info"

func GetAccountInstance() Account {

	once.Do(func() {
		instance, err := newAccount()
		if err != nil {
			panic(err)
		}
		go instance.monitorPrivateKey()
	})

	return instance
}

func newAccount() (*nbsAccount, error) {

	store := dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeLocalParam)
	if store == nil {
		return nil, fmt.Errorf("no available data store for account module")
	}

	obj := &nbsAccount{
		cleanKeyTimer: make(chan struct{}),
		dataStore:     store,
	}

	return obj, nil
}

func (account *nbsAccount) UnlockAccount(encodedKey string) (err error) {

	if account.accountId == nil {
		return fmt.Errorf("the account is not initialized yet")
	}

	privateKeyData, err := base64.StdEncoding.DecodeString(account.accountId.PrivateKey)
	if err != nil {
		return err
	}

	var decryptedData []byte

	if !account.accountId.Encrypted {
		decryptedData = privateKeyData
	} else {
		key, _ := hex.DecodeString(encodedKey)

		decryptedData = crypto.DecryptAES(privateKeyData, key)
	}

	account.privateKey, err = x509.ParsePKCS1PrivateKey(decryptedData)

	account.cleanKeyTimer <- struct{}{}

	return nil
}
func (account *nbsAccount) GetPeerID() string {
	return account.accountId.PeerID
}

func (account *nbsAccount) Reload() error {

	accountData, err := account.dataStore.Get(ParameterKeyForAccount)
	if err != nil {
		return err
	}

	accountInfo := &account_pb.Account{}
	if err := proto.Unmarshal(accountData, accountInfo); err != nil {
		return err
	}

	account.accountId.PeerID = accountInfo.PeerID
	account.accountId.PrivateKey = accountInfo.PrivateKey
	account.accountId.Encrypted = accountInfo.Encrypted

	return nil
}

func (account *nbsAccount) monitorPrivateKey() {

	for {
		select {
		case <-account.cleanKeyTimer:
			time.AfterFunc(5*time.Minute, func() {
				account.privateKey = nil
			})
		}
	}
}

func CreateAccount(password string) error {

	id, err := idService.GetInstance().GenerateId(password)
	if err != nil {
		return err
	}

	keyStore := dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeLocalParam)

	accountInfo := &account_pb.Account{
		PeerID:     id.PeerID,
		Encrypted:  id.Encrypted,
		PrivateKey: id.PrivateKey,
	}

	accountData, err := proto.Marshal(accountInfo)
	if err != nil {
		return err
	}

	if err := keyStore.Put(ParameterKeyForAccount, accountData); err != nil {
		return err
	}

	if err := GetAccountInstance().Reload(); err != nil {
		return err
	}

	return nil
}
