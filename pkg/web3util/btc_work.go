package web3util

import (
	"Open_IM/pkg/web3util/types"
	"github.com/btcsuite/btcd/chaincfg"
	"strings"
)

type BtcWorker struct {
}

func getNetwork(network string) *chaincfg.Params {
	var defaultNet *chaincfg.Params

	//指定网络 {MainNet：主网，TestNet：测试网，TestNet3：测试网3，SimNet：测试网}
	switch strings.ToLower(network) {
	case "mainnet":
		defaultNet = &chaincfg.MainNetParams
	case "testnet":
		defaultNet = &chaincfg.RegressionNetParams
	case "testnet3":
		defaultNet = &chaincfg.TestNet3Params
	case "simnet":
		defaultNet = &chaincfg.SimNetParams
	}
	return defaultNet
}

// CreateWallet 创建钱包
func (btcWork *BtcWorker) CreateWallet(defaultNet *chaincfg.Params) (*types.Wallet, error) {

	//1.生成私钥，参数：Secp256k1
	//privateKey, err := btcec.NewPrivateKey(btcec.S256())
	//if err != nil {
	//	return nil, err
	//}

	//2.转成wif格式
	//privateKeyWif, err := btcutil.NewWIF(privateKey, defaultNet, true)
	//if err != nil {
	//	return nil, err
	//}

	//return getWalletByPrivateKey(defaultNet, privateKeyWif)
	return nil, nil
}
