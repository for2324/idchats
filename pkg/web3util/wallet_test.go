package web3util

import (
	"Open_IM/pkg/utils"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestWallet_PublicKeyHex(t *testing.T) {

	//strbyte, _ := base64.StdEncoding.DecodeString("ho/1lmxPbhkRxPAkL+6lIlbpSo5aUFrw4Pu/QxApORtI0S6tejVO+ncb0CP/l3ReY0S/MR6I8KaYC60xVCa9mHwyuKQ8EOfU4S570xLC1aQ=")
	//stdm, _ := utils.AesDecrypt(strbyte, []byte("U2FsdGVkX1+4xoFd+2jiqf+m16e3EdEQ"))
	//fmt.Println(string(stdm))
	//return
	//暂时排除并发创建私钥的情况：
	//byteEntropy, _ := NewEntropy(128)
	//mnemonicstring, _ := NewMnemonicFromEntropy([]byte("emotion gloom swift door inform whip narrow donate drip weather maid hen"))
	walletPtr, err := NewFromMnemonic("repair upper super draw turn similar fever carry fiber build hat unfair")
	//fmt.Println(mnemonicstring)
	if err != nil {
		fmt.Println("123123123123")
		return
	}
	ptAccount := &accounts.Account{
		URL: accounts.URL{
			Scheme: "",
			Path:   "m/44'/60'/0'/0/0",
		},
	}
	ethprivateKeyHex, _ := walletPtr.PrivateKeyHex(*ptAccount)
	ethpublicAddress, _ := walletPtr.PublicKey(*ptAccount)
	fmt.Println("ethprivateKeyHex>>>>>>", string(ethprivateKeyHex))
	privateKeyCrypt, err := utils.AesEncrypt([]byte(ethprivateKeyHex), []byte("U2FsdGVkX1+4xoFd+2jiqf+m16e3EdEQ"))
	fmt.Println(crypto.PubkeyToAddress(*ethpublicAddress).String())
	fmt.Println(base64.StdEncoding.EncodeToString(privateKeyCrypt))
	ptAccount = &accounts.Account{
		URL: accounts.URL{
			Scheme: "",
			Path:   "m/84'/0'/0'/0/0",
		},
	}
	btcprivateKeyHex, _ := walletPtr.PrivateKeyHexBtc(*ptAccount)
	fmt.Println("btcprivateKeyHex >>>>>", btcprivateKeyHex)
	btcprivateKeyHexStr, err := utils.AesEncrypt([]byte(btcprivateKeyHex), []byte("U2FsdGVkX1+4xoFd+2jiqf+m16e3EdEQ"))
	btcpublicAddress, _ := walletPtr.PublicKeyHexBtc(*ptAccount)

	fmt.Println(base64.StdEncoding.EncodeToString(btcprivateKeyHexStr))
	fmt.Println(btcpublicAddress)

	ptAccount = &accounts.Account{
		URL: accounts.URL{
			Scheme: "",
			Path:   "m/44'/714'/0'/0/0",
		},
	}
	//ethprivateKeyHex, _ = walletPtr.PrivateKeyHexBnB(*ptAccount)
	//fmt.Println(ethprivateKeyHex)
	//
	//fmt.Println("bnbkey>>>>", ethprivateKeyHex)
	//ethpublicAddressStr, _ := walletPtr.PublicKeyHexBnB(*ptAccount)
	//fmt.Println(ethpublicAddressStr)
	//ptAccount = &accounts.Account{
	//	URL: accounts.URL{
	//		Scheme: "",
	//		Path:   "m/44'/195'/0'/0/0",
	//	},
	//}

	tronprivateKeyHex, _ := walletPtr.PrivateKeyHexTron(*ptAccount)
	tronpublicAddress, _ := walletPtr.PublicKeyHexTron(*ptAccount)

	fmt.Println("tronprivateKeyHex:", tronprivateKeyHex)
	fmt.Println("tronpublicAddress:", tronpublicAddress)

}
