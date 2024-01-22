package order

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"errors"
	"fmt"
	"strconv"
)

func NewBNBScanResolver(tag string, chainId int64) ScanResolver {
	return &BNBScanResolver{
		ChainId: chainId,
		Coin:    tag,
		EthScanResolver: EthScanResolver{
			ChainId: chainId,
			Coin:    tag,
		},
	}
}

type BNBScanResolver struct {
	EthScanResolver
	ChainId int64
	Coin    string
}

func (e *BNBScanResolver) GetCoinUSDPrice(OperationID string) (float64, error) {
	endPointConf := config.Config.ChainIdHttpMap[e.ChainId]
	if endPointConf.EndPoint == "" {
		return 0, errors.New("chainId not support")
	}
	uri := fmt.Sprintf("%s/api?module=stats&action=bnbprice&apikey=%s", endPointConf.EndPoint, endPointConf.ApiKey)
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
