package order

import (
	"Open_IM/pkg/common/config"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbOrder "Open_IM/pkg/proto/order"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestScanBlock(t *testing.T) {
	ScanBlock()
}

func TestCreateERC20ScanTask(t *testing.T) {
	s := &OrderServer{}
	resp, _ := s.CreatePayScanBlockTask(context.TODO(), &pbOrder.CreatePayScanBlockTaskReq{
		FormAddress: "0xcbd033ea3c05dc9504610061c86c7ae191c5c913",
		OrderId:     "123",
		USD:         1000000,
		TxnType:     "ERC20_MATIC_USDT",
	})
	log.NewInfo("TestCreateTask", "resp", resp)
}
func TestCreateMaticScanTask(t *testing.T) {
	s := &OrderServer{}
	resp, _ := s.CreatePayScanBlockTask(context.TODO(), &pbOrder.CreatePayScanBlockTaskReq{
		FormAddress: "0xcbd033ea3c05dc9504610061c86c7ae191c5c913",
		OrderId:     "123",
		USD:         0,
		TxnType:     "MATIC",
	})
	log.NewInfo("TestCreateTask", "resp", resp)
}
func TestCreateETHScanTask(t *testing.T) {
	s := &OrderServer{}
	resp, _ := s.CreatePayScanBlockTask(context.TODO(), &pbOrder.CreatePayScanBlockTaskReq{
		FormAddress: "0xcbd033ea3c05dc9504610061c86c7ae191c5c913",
		OrderId:     "123",
		USD:         1000000,
		TxnType:     "ETH",
	})
	log.NewInfo("TestCreateTask", "resp", resp)
}
func TestCreateBNBScanTask(t *testing.T) {
	s := &OrderServer{}
	resp, _ := s.CreatePayScanBlockTask(context.TODO(), &pbOrder.CreatePayScanBlockTaskReq{
		FormAddress: "0xcbd033ea3c05dc9504610061c86c7ae191c5c913",
		OrderId:     "345",
		USD:         0,
		TxnType:     "BNB",
	})
	log.NewInfo("TestCreateTask", "resp", resp)
}
func TestCreateARBScanTask(t *testing.T) {
	s := &OrderServer{}
	resp, _ := s.CreatePayScanBlockTask(context.TODO(), &pbOrder.CreatePayScanBlockTaskReq{
		FormAddress: "0xcbd033ea3c05dc9504610061c86c7ae191c5c913",
		OrderId:     "567",
		USD:         0,
		TxnType:     "ARB",
	})
	log.NewInfo("TestCreateTask", "resp", resp)
}
func TestCreateOPScanTask(t *testing.T) {
	s := &OrderServer{}
	resp, _ := s.CreatePayScanBlockTask(context.TODO(), &pbOrder.CreatePayScanBlockTaskReq{
		FormAddress: "0xcbd033ea3c05dc9504610061c86c7ae191c5c913",
		OrderId:     "567",
		USD:         0,
		TxnType:     "OP",
	})
	log.NewInfo("TestCreateTask", "resp", resp)
}

func TestLoopScanBlock(t *testing.T) {
	LoopScanBlock()
}

func TestScanMaticEthBlock(t *testing.T) {
	s := NewEthScanResolver("MATIC", 80001)
	s.Scan(
		[]string{"0xfC004E9052Fd1740A662Fac99c61C9cC73D41Db8"},
		[]string{"0xE5A9748caB1A4A5756563C0Eb4a281A9345e0acD"},
		36629087,
		36629087,
	)
}
func TestCompareStartAndEndBanlance(t *testing.T) {
	s := NewEthScanResolver("MATIC", 80001)
	// 9170172
	var start uint64 = 9170166
	var end uint64 = 9170236
	s.CompareStartAndEndBanlance("ETH", "0xcbd033ea3c05dc9504610061c86c7ae191c5c913", start, end)
}

func TestScanBNBEthBlock(t *testing.T) {
	// s := NewEthScanResolver("BSC", 97)
	// s.Scan(
	// 	[]string{"0xf474cf03cceff28abc65c9cbae594f725c80e12d"},
	// 	[]string{"0xcbd033ea3c05dc9504610061c86c7ae191c5c913"},
	// 	30538763,
	// 	30538763,
	// )
	cli, err := utils.GetEthClient(97)
	if err != nil {
		t.Errorf("Failed to connect to Ethereum network: %v", err)
		return
	}
	txHas := "0x6a3445e9dd85da0c208df3a59bccfab2fd9f0cb076fccaed38165a3e9f171135"
	tx, _, err := cli.TransactionByHash(context.TODO(), common.HexToHash(txHas))
	if err != nil {
		t.Errorf("Failed to get transaction: %v", err)
		return
	}
	log.NewInfo("TestScanBNBEthBlock", "tx", tx.To().Hex())
}

func TestScanETHBlock(t *testing.T) {
	// formAddresses:  [0xcbd033ea3c05dc9504610061c86c7ae191c5c913] , toAddreses:  [0xe5a9748cab1a4a5756563c0eb4a281a9345e0acd] , start:  9186267 , end:  9186268]
	s := NewEthScanResolver("ETH", 5)
	s.Scan(
		[]string{"0xcbd033ea3c05dc9504610061c86c7ae191c5c913"},
		[]string{"0xe5a9748cab1a4a5756563c0eb4a281a9345e0acd"},
		9186267,
		9186268,
	)
}

func TestGetBalanceETH(t *testing.T) {
	same, err := CompareStartAndEndBanlance(
		"CompareStartAndEndBanlance", 5, "0xE5A9748caB1A4A5756563C0Eb4a281A9345e0acD", 9174890, 9174910,
	)
	if err != nil {
		t.Errorf("Failed to get balance: %v", err)
		return
	}
	log.NewInfo("TestGetBalanceETH", "same", same)
}

// func TestERC20GasPrice(t *testing.T) {
// 	constract := "0x48c04ed5691981C42154C6167398f95e8f38a7fF"
// 	gasfee, err := SuggestERC20GasFee(1, constract)
// 	if err != nil {
// 		t.Errorf("Failed to get gas fee: %v", err)
// 		return
// 	}
// 	log.NewInfo("TestERC20GasPrice", "gasfee", gasfee)
// }

func TestSendOrderToKafka(t *testing.T) {
	producerToOrder = kafka.NewKafkaProducer(config.Config.Kafka.MsgOrder.Addr, config.Config.Kafka.MsgOrder.Topic)
	OperationID := "TestSendOrderToKafka"
	orderInfo, err := imdb.GetPayScanBlockTaskByOrderId("2")
	if err != nil {
		log.Error(OperationID, utils.GetSelfFuncName(), "get order info failed", "err", err.Error())
		return
	}
	key := fmt.Sprintf("%s:%s", orderInfo.Mark, "123")
	pid, offset, err := producerToOrder.SendMessage(
		&pbOrder.MsgDataToOrderByMQ{
			OperationID: OperationID,
			USD:         orderInfo.USDPrice,
			OrderId:     orderInfo.OrderId,
			UserId:      orderInfo.FormAddress,
			TxnType:     orderInfo.TxnType,
			Mark:        orderInfo.Mark,
			Value:       orderInfo.Value,
			Decimal:     orderInfo.Decimal,
			TxnHash:     orderInfo.TxnHash,
		},
		key, OperationID)
	if err != nil {
		log.Error(OperationID, utils.GetSelfFuncName(), "kafka send failed", "send data", orderInfo.Id, "pid", pid, "offset", offset, "err", err.Error(), "key", key)
		return
	}
	log.NewInfo(OperationID, "kafka send success", "send data", orderInfo.Id, "pid", pid, "offset", offset, "key", key)
}
