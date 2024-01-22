package web3util

import (
	"Open_IM/pkg/xlog"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestNewEthWorker(t *testing.T) {

	got, err := NewEthWorker(3, "https://data-seed-prebsc-1-s2.bnbchain.org:8545")
	if err != nil {
		xlog.CError(err.Error())
		return
	}
	xlog.InitLevel(xlog.DefaultOption())
	xlog.CInfo("创建钱包")
	//wallet, err := got.CreateWallet()
	//if err != nil {
	//	xlog.CError(err.Error())
	//	return
	//}
	//xlog.CInfo(utils.StructToJsonString(wallet))
	amount := new(big.Int)
	amount.SetString("1000000000000000", 10) // 0.001 ETH = 1e15 wei
	//提现母币
	_, hash, _, err := got.Transfer("261975ba2e1f7af2b43c2fbed32759862ae197da148f33afc493c52e0b6f0310",
		"0xC60e1c8FFbd885EBf13E577F80EcBa41D235A455",
		amount, 0, "")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(hash)
	}

	//提取erc20
	time.Sleep(3 * time.Second)

}
