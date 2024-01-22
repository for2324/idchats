package main

import (
	"Open_IM/internal/contract"

	"github.com/supermigo/xlog"
)

func main() {
	xlog.InitLevel(xlog.DefaultOption())
	xlog.CInfo("123123123123123123123")
	contract.InitLocalDB("./scanblocksqlite.db")
	go contract.StartScanBlockFilterQuery()
	select {}
}
