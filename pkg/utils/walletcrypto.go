package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
)

type EncryptedVault struct {
	Salt                 string `json:"salt"`
	InitializationVector string `json:"initializationVector"`
	CipherText           string `json:"cipherText"`
}

func deriveKey(passphrase string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		rand.Read(salt)
	}
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}

func WalletEncrypt(passphrase, plaintext string) (EncryptedVault, error) {
	key, salt := deriveKey(passphrase, nil)
	iv := make([]byte, 12)
	rand.Read(iv)
	b, errr := aes.NewCipher(key)
	if errr != nil {
		return EncryptedVault{}, errr
	}
	aesgcm, errr := cipher.NewGCM(b)
	if errr != nil {
		return EncryptedVault{}, errr
	}
	data := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return EncryptedVault{
		Salt:                 base64.StdEncoding.EncodeToString(salt),
		InitializationVector: base64.StdEncoding.EncodeToString(iv),
		CipherText:           base64.StdEncoding.EncodeToString(data),
	}, nil
}

func WalletDecrypt(passphrase string, ciphertext EncryptedVault) string {
	salt, _ := base64.StdEncoding.DecodeString(ciphertext.Salt)
	iv, _ := base64.StdEncoding.DecodeString(ciphertext.InitializationVector)
	data, _ := base64.StdEncoding.DecodeString(ciphertext.CipherText)
	key, _ := deriveKey(passphrase, salt)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data, _ = aesgcm.Open(nil, iv, data, nil)
	return string(data)
}
