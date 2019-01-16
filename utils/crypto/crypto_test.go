package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"github.com/mr-tron/base58/base58"
	"github.com/multiformats/go-multihash"
	"testing"
)

var privateKeyStr string
var privateKey *rsa.PrivateKey

func TestEncryptAES(t *testing.T) {

	pri, _, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatal("GenerateRSAKeyPair failed.")
	}
	privateKey = pri

	privateKeyData := x509.MarshalPKCS1PrivateKey(privateKey)

	password := MD5SB("Wesley")

	encryptedPrivateKeyData := EncryptAES(privateKeyData, []byte(password))

	privateKeyStr = base64.StdEncoding.EncodeToString(encryptedPrivateKeyData)

	t.Log("EncryptAES 测试通过")
}

func TestDecryptAES(t *testing.T) {

	privateKeyData, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		t.Fatal("DecodeString private key failed.")
	}

	password := MD5SB("Wesley")

	decryptedData := DecryptAES(privateKeyData, []byte(password))

	privateKeyTest, err := x509.ParsePKCS1PrivateKey(decryptedData)

	if err != nil {
		t.Fatal("ParsePKCS1PrivateKey failed.")
	}

	if privateKey.D.Cmp(privateKeyTest.D) != 0 {
		t.Fatalf("decrypt failed:%x-----%x", privateKey.D, privateKeyTest.D)
	}

	t.Log("DecryptAES 测试通过")
}

func TestPeerID(t *testing.T) {
	_, pub, err := GenerateRSAKeyPair()
	if err != nil {
		t.Fatal("GenerateRSAKeyPair failed.", err)
	}

	pubData := x509.MarshalPKCS1PublicKey(pub)

	hash, _ := multihash.Sum(pubData, multihash.SHA2_256, -1)

	id := string(hash)

	t.Log("PeerID 测试通过:", base58.Encode([]byte(id)))
}
