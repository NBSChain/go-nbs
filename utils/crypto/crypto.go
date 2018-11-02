package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)

	return md5str
}


const (
	nBitsForKeyPairDefault = 2048
)

func GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, nBitsForKeyPairDefault)
	if err != nil{
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}


func padding(src []byte, blockSize int) []byte {

	padNum 	:= blockSize - len(src) % blockSize

	pad	:= bytes.Repeat([]byte{byte(padNum)}, padNum)

	return append(src,pad...)
}

func unPadding(src []byte) []byte {

	n := len(src)

	unPadNum :=int(src[n - 1])

	return src[:n - unPadNum]
}

func EncryptAES(src []byte, key []byte) []byte {

	block, _ := aes.NewCipher(key)

	src = padding(src, block.BlockSize())
	dst := make ([]byte, 0, len(src))

	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(dst, src)

	return dst
}

func DecryptAES(src []byte, key []byte) []byte {

	block, _ := aes.NewCipher(key)

	blockMode:=cipher.NewCBCDecrypter(block, key)

	dst := make([]byte, 0, len(src))

	blockMode.CryptBlocks(dst, src)

	dst = unPadding(dst)

	return dst
}

//TODO::account password logic check
func CheckPassword(password string) error{
	if len(password) < 6 {
		return fmt.Errorf("password is too short, must be longer than 6")
	}
	return nil
}