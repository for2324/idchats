package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventConsumer struct {
	Address    common.Address //合约地址
	Topic      common.Hash    //topic内容
	ParseEvent EventParser    //解析log的回掉
	Contract   interface{}    //合约ptr 合约指针
}

type IEventConsumer interface {
	GetConsumer() ([]*EventConsumer, error)
	GetFilterQuery() []ethereum.FilterQuery
}

type EventParser func(types.Log, uint64) error

// ITxMonitor -----------------------------------------------------------
type ITxMonitor interface {
	MonitoredAddress() common.Address
	TxParse(t *types.Transaction, from, to, tokenAddr, amount string, data []byte) error
}

func KeyFromBEConsumer(address string, topic string) string {
	return fmt.Sprintf("%s:%s", address, topic)
}
