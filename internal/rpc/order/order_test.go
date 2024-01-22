package order

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

func TestEncrypt(t *testing.T) {
	// 加密
	plaintext := []byte("example plaintext")
	key := []byte("examplekey123456")
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	fmt.Printf("%x\n", ciphertext)

	// 解密
	block, err = aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	openData, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(openData))
}
