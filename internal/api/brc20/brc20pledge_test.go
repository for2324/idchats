package brc20

import (
	"Open_IM/internal/api/brc20/services/unisat_wallet"
	"testing"
)

func TestGetBrc20TransferableList(t *testing.T) {
	tempdata := new(unisat_wallet.UnisatWeb)
	tempdata.NetParam = "testnet"
	tempdata.LastScanTime = 1703217600
	tempdata.Ticker = "dead"
	tempdata.ScanUnisatWallet()
}
