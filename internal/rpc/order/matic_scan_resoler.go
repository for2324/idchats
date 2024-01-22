package order

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"errors"
	"fmt"
	"strconv"
)

type MaticScanResolver struct {
	EthScanResolver
	ChainId int64
	Coin    string
}

func NewMaticScanResolver(tag string, chainId int64) ScanResolver {
	return &MaticScanResolver{
		ChainId: chainId,
		Coin:    tag,
		EthScanResolver: EthScanResolver{
			ChainId: chainId,
			Coin:    tag,
		},
	}
}

// func (e *MaticScanResolver) Scan(fromAddresses, toAddress []string, start, end uint64) ([]*TransferEvent, error) {
// 	// 创建以太坊客户端
// 	client, err := utils.GetEthClient("EthScanResolver", e.ChainId)
// 	if err != nil {
// 		log.Error("get eth client failed, err: ", err.Error())
// 		return nil, err
// 	}
// 	defer client.Close()
// 	// 监听指定地址的转账信息
// 	Addresses := []common.Hash{}
// 	for _, address := range fromAddresses {
// 		Addresses = append(Addresses, common.HexToHash(address))
// 	}
// 	query := ethereum.FilterQuery{
// 		FromBlock: big.NewInt(int64(start)),
// 		ToBlock:   big.NewInt(int64(end)),
// 		Topics: [][]common.Hash{
// 			{
// 				common.HexToHash("0x4dfe1bbbcf077ddc3e01291eea2d5c70c2b422b415d95645b9adcfd678cb1d63"),
// 			}, {
// 				common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000001010"),
// 			},
// 			Addresses,
// 		},
// 		Addresses: []common.Address{common.HexToAddress("0x0000000000000000000000000000000000001010")},
// 	}
// 	events, err := client.FilterLogs(context.Background(), query)
// 	if err != nil {
// 		log.Error("filter logs failed, err: ", err.Error())
// 		return nil, err
// 	}
// 	transferEvent := make([]*TransferEvent, 0)
// 	for _, event := range events {
// 		if event.Removed {
// 			continue
// 		}
// 		if len(event.Topics) < 3 {
// 			continue
// 		}
// 		log.Debug("Transfer evnet",
// 			common.HexToAddress(event.Topics[0].Hex()).String(),
// 			common.HexToAddress(event.Topics[1].Hex()).String(),
// 			common.HexToAddress(event.Topics[2].Hex()).String(),
// 		)
// 		from := common.HexToAddress(event.Topics[2].Hex()).String()
// 		// 获取交易数据
// 		tx, _, err := client.TransactionByHash(context.Background(), event.TxHash)
// 		if err != nil {
// 			log.Error("get transaction by hash failed, err: ", err.Error())
// 			continue
// 		}
// 		for _, address := range toAddress {
// 			if strings.EqualFold(address, tx.To().String()) {
// 				transferEvent = append(transferEvent, &TransferEvent{
// 					From:        from,
// 					To:          tx.To().String(),
// 					Value:       tx.Value(),
// 					BlockNumber: event.BlockNumber,
// 					TxHash:      event.TxHash.String(),
// 				})
// 			}
// 		}

// 	}
// 	return transferEvent, nil
// }

// func (e *MaticScanResolver) CompareStartAndEndBanlance(OperationID string, account string, start, end uint64) (bool, error) {
// 	return CompareStartAndEndBanlance(OperationID, e.ChainId, account, start, end)
// }

type CoinMaticUSDPriceResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Maticbtc          string `json:"maticbtc"`
		MaticbtcTimestamp string `json:"maticbtc_timestamp"`
		Maticusd          string `json:"maticusd"`
		MaticusdTimestamp string `json:"maticusd_timestamp"`
	} `json:"result"`
}

func (e *MaticScanResolver) GetCoinUSDPrice(OperationID string) (float64, error) {
	// https://api-goerli.etherscan.io/api?module=stats&action=ethprice&apikey=
	endPointConf := config.Config.ChainIdHttpMap[e.ChainId]
	if endPointConf.EndPoint == "" {
		return 0, errors.New("chainId not support")
	}
	uri := fmt.Sprintf("%s/api?module=stats&action=maticprice&apikey=%s", endPointConf.EndPoint, endPointConf.ApiKey)
	var priceResult CoinMaticUSDPriceResult
	err := ChainHttpGet(uri, &priceResult)
	if err != nil {
		return 0, err
	}
	if priceResult.Status != "1" {
		return 0, errors.New(priceResult.Message)
	}
	Ethusd := priceResult.Result.Maticusd
	priceRate, err := strconv.ParseFloat(Ethusd, 64)
	if err != nil {
		return 0, err
	}
	log.NewInfo(OperationID, utils.GetSelfFuncName(), "e.ChainId", e.ChainId, "resp", priceResult)

	return priceRate, nil
}
