package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestWalletDecrypt(t *testing.T) {
	//ttdata := &EncryptedVault{
	//	Salt:                 "GGojlmLewrU=",
	//	InitializationVector: "6c9N7t65C/YwEhiB",
	//	CipherText:           "fwiwqicrUj7N6kOHjzkdlac7y4ehAZ6P/GzNNME6931hMdT/DA59sfoaivHy6mXOFa2uAjLRG3TZn3GxPVUI3n6BZPWDhZkrj52BMkZoKB400x6HGgtgo9XRdS8wS4o2bTGf",
	//}
	//tt := VerifySignature("0xfc004e9052fd1740a662fac99c61c9cc73d41db8",
	//	"0x3903a7a394a9c13dcfe6fc895e45d692c1c21ee042e72b170523c1e4e4842e7e41bf441aaedf017184b87b35077e6018d6769d1f1cf8202dfb94560d18660b261b",
	//	"0xfc004e9052fd1740a662fac99c61c9cc73d41db8")
	//fmt.Println(tt)
	//ttt := WalletDecrypt("0x9086f9624c0afbfe75ba35c8d21ceb8871955b950bee331d00862f74fe2ef07571bf6f08a0425ddb93149b225602e47fcd3a688bd2c23a510138b54126058c691c", *ttdata)
	ttdata := &EncryptedVault{
		Salt:                 "ag4mFbF1hk4=",
		InitializationVector: "qJBs2nxvF0VRLZcz",
		CipherText:           "dwwtf3gj80asAuZcaga0OafghqCH/6UTx6ohNK2b085g4e5NDONXieSHZp3ArMoeaalsmqCUsGbKPS/+Q3Ukk4gYUWoxZAk3Z6tAABLlz5Xkkh71jQWBDMIIcqLJnoyl5nKKdVeaohZiaZaqmB/pN/K6xdJrhO0KNSKLh8KFUBOSeZ4TlUInQcMqq+aFQbkXz2pTYQ==",
	}
	ttt := WalletDecrypt("0x5d203073a0f19db7a5c35a9a3e1340972a4e294fdfaaaf0d9b84e5b742f17b1f1fffb1340d0ec9e0b4d296021f11d21f05ef864fdff99e5bf47612a1f2e5f3ba1c", *ttdata)
	fmt.Println(ttt)
}
func TestWalletEncrypt(t *testing.T) {
	salt, _ := base64.StdEncoding.DecodeString("ag4mFbF1hk4=")
	fmt.Println(salt)
	key, salt := deriveKey("0x5d203073a0f19db7a5c35a9a3e1340972a4e294fdfaaaf0d9b84e5b742f17b1f1fffb1340d0ec9e0b4d296021f11d21f05ef864fdff99e5bf47612a1f2e5f3ba1c",
		nil)
	fmt.Println(salt)
	iv, _ := base64.StdEncoding.DecodeString("qJBs2nxvF0VRLZcz")
	b, errr := aes.NewCipher(key)
	if errr != nil {
		//return EncryptedVault{}, errr
		return
	}
	aesgcm, errr := cipher.NewGCM(b)
	if errr != nil {
		//return EncryptedVault{}, errr
		return
	}
	data := aesgcm.Seal(nil, iv, []byte("great mansion mushroom session correct leopard knife black drastic system friend tape"), nil)
	tt := EncryptedVault{
		Salt:                 base64.StdEncoding.EncodeToString(salt),
		InitializationVector: base64.StdEncoding.EncodeToString(iv),
		CipherText:           base64.StdEncoding.EncodeToString(data),
	}
	//	tt, err := WalletEncrypt("", "great mansion mushroom session correct leopard knife black drastic system friend tape")
	//ttdata := &EncryptedVault{
	//	Salt:                 "ag4mFbF1hk4=",
	//	InitializationVector: "qJBs2nxvF0VRLZcz",
	//	CipherText:           "dwwtf3gj80asAuZcaga0OafghqCH/6UTx6ohNK2b085g4e5NDONXieSHZp3ArMoeaalsmqCUsGbKPS/+Q3Ukk4gYUWoxZAk3Z6tAABLlz5Xkkh71jQWBDMIIcqLJnoyl5nKKdVeaohZiaZaqmB/pN/K6xdJrhO0KNSKLh8KFUBOSeZ4TlUInQcMqq+aFQbkXz2pTYQ==",
	//}
	//ttt := WalletDecrypt("0x5d203073a0f19db7a5c35a9a3e1340972a4e294fdfaaaf0d9b84e5b742f17b1f1fffb1340d0ec9e0b4d296021f11d21f05ef864fdff99e5bf47612a1f2e5f3ba1c", *ttdata)
	fmt.Println(tt)
	//fmt.Println(err)
	//fmt.Println(ttt)
}

// {\"userID\":\"0x68efbe4733b6afe100d3163ac76ab376538264e5\",\"createAt\":\"2023-10-16 15:29:17\",\"sign\":\"0x5d203073a0f19db7a5c35a9a3e1340972a4e294fdfaaaf0d9b84e5b742f17b1f1fffb1340d0ec9e0b4d296021f11d21f05ef864fdff99e5bf47612a1f2e5f3ba1c\",\"s\":\"ag4mFbF1hk4=\",\"i\":\"qJBs2nxvF0VRLZcz\",\"c\":\"dwwtf3gj80asAuZcaga0OafghqCH/6UTx6ohNK2b085g4e5NDONXieSHZp3ArMoeaalsmqCUsGbKPS/+Q3Ukk4gYUWoxZAk3Z6tAABLlz5Xkkh71jQWBDMIIcqLJnoyl5nKKdVeaohZiaZaqmB/pN/K6xdJrhO0KNSKLh8KFUBOSeZ4TlUInQcMqq+aFQbkXz2pTYQ==\"}","duration":"926.2ms","level":"slow","span":"5618ce48265b6dc9","trace":"8ce3effc46b8bbe11ffa831015d8e90a"}
