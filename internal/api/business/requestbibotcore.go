package business

import (
	"Open_IM/pkg/common/config"
	"github.com/guonaihong/gout"
	"github.com/shopspring/decimal"
)

// 请求后端服务 获取数据
// http://192.168.1.99:12002/api/swap/record/statistics?merchantUid=0xad9ebc6b862f65720d2e5329319a8832f4e08d6a&merchantId=bibot
// 获取
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
	var intPutRequestData gout.H
	if len(apiKey) > 0 {
		intPutRequestData["apiKey"] = apiKey
	}
	intPutRequestData["merchantUid"] = merchantUid
	intPutRequestData["merchantId"] = merchantId
	err = gout.GET(config.Config.UniswapRobot.BibotUri + "/api/swap/record/statistics").SetQuery(intPutRequestData).Debug(!config.Config.IsPublicEnv).BindJSON(&result).Do()
	return
}

// 万分之，避免0.95*2  超过整数部分
func GetCurrentFee(merchantUid string, merchantId string, apiKey []string) (tradeFeeRate, sniperFeeRate int) {
	resultData, err := RequestGetTotalTrade(merchantUid, merchantId, nil)
	if err != nil {
		return 20, 100
	}
	if resultData.ErrCode != 0 {
		return 20, 100
	}
	subTotalTrade, _ := decimal.NewFromString("0")
	for _, value := range resultData.Data {
		subTotalTrade = subTotalTrade.Add(decimal.NewFromInt(value.SellUsdPrice))
	}
	if subTotalTrade.GreaterThan(decimal.NewFromInt(3000_000)) {
		return 20 * 0.9, 100 * 0.9
	}
	if subTotalTrade.GreaterThan(decimal.NewFromInt(1000_000)) {
		return 20 * 0.95, 100 * 0.95
	}
	return 20, 100
}
