package rocksCache

import (
	"Open_IM/pkg/common/config"
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/shopspring/decimal"
)

type TQuestTradeTotalData struct {
	Type           string  `json:"type"`
	SellAmount     int64   `json:"sellAmount"`
	SellUsdPrice   int64   `json:"sellUsdPrice"`
	FeeUsdPrice    int     `json:"feeUsdPrice"`
	ActualSlippage float64 `json:"actualSlippage"`
	ApiKey         string  `json:"apiKey"`
}
type TQuestTradeTotalResp struct {
	Data    []*TQuestTradeTotalData `json:"data"`
	ErrCode int                     `json:"errCode"`
	ErrMsg  string                  `json:"errMsg"`
}

func RequestGetTotalTrade(merchantUid string, merchantId string, apiKey []string) (result *TQuestTradeTotalResp, err error) {
	var intPutRequestData = make(gout.H, 0)
	if len(apiKey) > 0 {
		intPutRequestData["apiKey"] = apiKey
	}

	intPutRequestData["merchantUid"] = merchantUid
	if merchantId != merchantUid {
		intPutRequestData["merchantId"] = merchantId
	}

	err = gout.GET(config.Config.UniswapRobot.BibotUri + "/api/swap/record/statistics").SetQuery(intPutRequestData).Debug(!config.Config.IsPublicEnv).BindJSON(&result).Do()
	return
}
func RequestGetTotalBuyApiTrade(merchantUid string, merchantId string, apiKey []string) (result *TQuestTradeTotalResp, err error) {
	var intPutRequestData = make(gout.H, 0)
	if len(apiKey) == 0 {
		result = new(TQuestTradeTotalResp)
		result.Data = nil
		return nil, nil
	} else {
		intPutRequestData["apiKey"] = apiKey
	}

	err = gout.GET(config.Config.UniswapRobot.BibotUri + "/api/swap/record/apiKey/statistics").
		SetQuery(intPutRequestData).Debug(!config.Config.IsPublicEnv).BindJSON(&result).Do()
	return
}

// 万分之，避免0.95*2  超过整数部分
func GetCurrentFee(merchantUid string, merchantId string, apiKey []string) (tradeFeeRate, sniperFeeRate float64) {
	if merchantUid == "" && merchantId == "" {
		return 0.002, 0.01
	}
	resultDataBiBotServer, err := RequestGetTotalTrade(merchantUid, merchantId, nil)
	if err != nil {
		return 0.002, 0.01
	}
	if resultDataBiBotServer.ErrCode != 0 {
		return 0.002, 0.01
	}
	subTotalTrade, _ := decimal.NewFromString("0")
	for _, value := range resultDataBiBotServer.Data {
		subTotalTrade = subTotalTrade.Add(decimal.NewFromInt(value.SellUsdPrice).Shift(-6))
	}
	fmt.Println("subTotalTrade>>>>>>>>>>>>>", subTotalTrade.String())
	if subTotalTrade.GreaterThan(decimal.NewFromInt(300_000)) {
		return 0.002 * 0.9, 0.01 * 0.9
	}
	if subTotalTrade.GreaterThan(decimal.NewFromInt(100_000)) {
		return 0.002 * 0.95, 0.01 * 0.95
	}
	return 0.002, 0.01
}
func GetCurrentTradeParentRewardFee(strDecimalData string) (tradeFeeRate float64) {
	decimalData, _ := decimal.NewFromString(strDecimalData)
	decimalData = decimalData.Shift(-6)
	if decimalData.GreaterThan(decimal.NewFromInt(100_000)) {
		return 0.2
	}
	if decimalData.GreaterThan(decimal.NewFromInt(10_000)) {
		return 0.1
	}
	return 0.05
}
