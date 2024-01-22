package order

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ChainHttpGet(uri string, v any) error {
	var resutlbyte []byte
	var err error
	if config.Config.OpenNetProxy.OpenFlag {
		proxyAddress, _ := url.Parse(config.Config.OpenNetProxy.ProxyURL)
		resutlbyte, err = utils.HttpGetWithHeaderWithProxy(uri, map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		}, http.ProxyURL(proxyAddress))
	} else {
		resutlbyte, err = utils.HttpGetWithHeader(uri, map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		})
	}
	if err != nil {
		log.NewError("Failed to send API request: %v", uri, err)
		return err
	}
	body := resutlbyte
	// Parse the JSON response
	if err := json.Unmarshal(body, &v); err != nil {
		log.NewError("Failed to parse API response: %v", uri, err)
		return err
	}
	return nil
}

func NewEthScanResolver(tag string, chainId int64) ScanResolver {
	return &EthScanResolver{
		ChainId: chainId,
		Coin:    tag,
	}
}

type EthScanResolver struct {
	ChainId int64
	Coin    string
}

func (e *EthScanResolver) Scan(fromAddresses, toAddress []string, start, end uint64) ([]*TransferEvent, error) {
	confInfo, ok := config.Config.ChainIdHttpMap[e.ChainId]
	if !ok {
		log.NewError("Failed to get chainId http map")
		return nil, errors.New("Failed to get chainId http map")
	}
	transferEvent := make([]*TransferEvent, 0)
	for _, to := range toAddress {
		txns, err := e.GetAccountTxns(confInfo.EndPoint, to, start, end, confInfo.ApiKey)
		if err != nil {
			log.NewError("Failed to get account txns: %v", err)
			return nil, err
		}
		for _, txn := range txns {
			// Check if the transaction is a transfer
			if strings.EqualFold(txn.To, to) && txn.Input == "0x" {
				// Parse the transaction value
				value := new(big.Int)
				value.SetString(txn.Value, 10)
				blockNumber, err := strconv.ParseInt(txn.BlockNumber, 0, 64)
				if err != nil {
					log.NewError("Failed to parse transaction blockNumber: %v", err)
					return nil, err
				}
				// Create the transfer event

				transferEvent = append(transferEvent, &TransferEvent{
					From:        txn.From,
					To:          txn.To,
					Value:       value,
					BlockNumber: uint64(blockNumber),
					TxHash:      txn.Hash,
				})
			}
		}
	}
	return transferEvent, nil
}

type ScanAccountTxResult struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	TxList  []ScanTransferEventItem `json:"result"`
}

// EtherscanTx represents an Ethereum transaction returned by the Etherscan API.
type ScanTransferEventItem struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Gas         string `json:"gas"`
	GasPrice    string `json:"gasPrice"`
	Input       string `json:"input"`
	BlockNumber string `json:"blockNumber"`
}

func (e *EthScanResolver) GetNodeBlockNumber(endpoint string, apiKey string) (uint64, error) {
	uri := fmt.Sprintf("%s/api?module=proxy&action=eth_blockNumber&apikey=%s", endpoint, apiKey)
	// Parse the JSON response
	var result struct {
		Result string `json:"result"`
	}
	err := ChainHttpGet(uri, &result)
	if err != nil {
		return 0, err
	}
	log.NewInfo("eth GetNodeBlockNumber result", uri, result)
	blockNumber, err := strconv.ParseUint(result.Result, 0, 64)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (e *EthScanResolver) GetAccountTxns(
	endpoint string,
	address string,
	startBlock uint64,
	endBlock uint64,
	apiKey string,
) ([]ScanTransferEventItem, error) {
	page := 100
	offset := 0
	sort := "desc" // or "desc"

	// Create the API request URL
	uri := fmt.Sprintf("%s/api?module=account&action=txlist&address=%s&startblock=%d&endblock=%d&page=%d&offset=%d&sort=%s&apikey=%s",
		endpoint, address, startBlock, endBlock, page, offset, sort, apiKey)
	// Send the API request
	var result ScanAccountTxResult
	err := ChainHttpGet(uri, &result)
	if err != nil {
		return nil, err
	}
	log.NewDebug("eth GetAccountTxns uri:", uri, "result:", result)
	if result.Status == "1" { // result.Status == "0"
		return result.TxList, nil
	}
	return nil, errors.New(result.Message)
}

func (e *EthScanResolver) CompareStartAndEndBanlance(OperationID string, account string, start, end uint64) (bool, error) {
	return CompareStartAndEndBanlance(OperationID, e.ChainId, account, start, end)
}

//	{
//		"status":"1",
//		"message":"OK",
//		"result":{
//		   "ethbtc":"0.07311",
//		   "ethbtc_timestamp":"1636108470",
//		   "ethusd":"4507.38",
//		   "ethusd_timestamp":"1636108466"
//		}
//	 }
type CoinUSDPriceResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Ethbtc          string `json:"ethbtc"`
		EthbtcTimestamp string `json:"ethbtc_timestamp"`
		Ethusd          string `json:"ethusd"`
		EthusdTimestamp string `json:"ethusd_timestamp"`
	} `json:"result"`
}

// 根据 USDPrice 去交易所获取币种价格(1COIN = ? USDT ) gasPrice (coin为单位)
func (e *EthScanResolver) GetCoinUSDPrice(OperationID string) (float64, error) {
	// https://api-goerli.etherscan.io/api?module=stats&action=ethprice&apikey=
	endPointConf := config.Config.ChainIdHttpMap[e.ChainId]
	if endPointConf.EndPoint == "" {
		return 0, errors.New("chainId not support")
	}
	uri := fmt.Sprintf("%s/api?module=stats&action=ethprice&apikey=%s", endPointConf.EndPoint, endPointConf.ApiKey)
	var priceResult CoinUSDPriceResult
	err := ChainHttpGet(uri, &priceResult)
	if err != nil {
		return 0, err
	}
	if priceResult.Status != "1" {
		return 0, errors.New(priceResult.Message)
	}
	Ethusd := priceResult.Result.Ethusd
	priceRate, err := strconv.ParseFloat(Ethusd, 64)
	if err != nil {
		return 0, err
	}
	log.NewInfo(OperationID, utils.GetSelfFuncName(), "e.ChainId", e.ChainId, "resp", priceResult)
	return priceRate, nil
}

// func (e *EthScanResolver) GetBlock(OperationID string, ctx context.Context, chainId int64, number *big.Int) (*ChainBlock, error) {
// 	client, err := iutils.GetEthClient(chainId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	block, err := client.BlockByNumber(ctx, number)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &ChainBlock{
// 		BlockNumber: block.Number().Uint64(),
// 		Time:        time.Unix(int64(block.Time()), 0),
// 	}, nil
// }

func (e *EthScanResolver) GetBlock(OperationID string, ctx context.Context, chainId int64, number *big.Int) (*ChainBlock, error) {
	endPointConf := config.Config.ChainIdHttpMap[e.ChainId]
	if endPointConf.EndPoint == "" {
		return nil, errors.New("chainId not support")
	}
	if number == nil {
		uri := fmt.Sprintf("%s/api?module=proxy&action=eth_blockNumber&apikey=%s", endPointConf.EndPoint, endPointConf.ApiKey)
		var blockNumberResult struct {
			Result string `json:"result"`
		}
		err := ChainHttpGet(uri, &blockNumberResult)
		if err != nil {
			return nil, err
		}
		number = new(big.Int)
		number.SetString(blockNumberResult.Result, 0)
	}
	bnTag := fmt.Sprintf("0x%x", number)
	uri := fmt.Sprintf("%s/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=%s", endPointConf.EndPoint, bnTag, endPointConf.ApiKey)
	block := BlockResultJson{}
	err := ChainHttpGet(uri, &block)
	if err != nil {
		return nil, err
	}
	timeUnix := new(big.Int)
	timeUnix.SetString(block.Result.Timestamp, 0)
	return &ChainBlock{
		BlockNumber: number.Uint64(),
		Time:        time.Unix(int64(timeUnix.Uint64()), 0),
	}, nil
}
