package utils

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	data, _ := AesEncrypt([]byte("70009118e567d4715d1d39caf0112b25849e2a7fa93750b1e81ca24c0f571a54"), []byte("U2FsdGVkX1+4xoFd+2jiqf+m16e3EdEQ"))
	fmt.Println(base64.StdEncoding.EncodeToString(data))

}
