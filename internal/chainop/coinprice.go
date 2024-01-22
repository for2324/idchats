package chainop

import (
	"Open_IM/pkg/common/config"
	utils2 "Open_IM/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type CoinPriceReq struct {
	CoinName string `json:"coinName"  bind:"require"`
}
type CoinPriceResp struct {
	ErrCode int
	ErrMsg  string
	Data    interface{}
}

// CoinPrice 查询币价的合约 包括母币和swap
func CoinPrice(c *gin.Context) {
	var req CoinPriceReq
	var responseData CoinPriceResp

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	//spot 现货交易
	client := resty.New()
	url := fmt.Sprintf(
		"https://api.coinmarketcap.com/data-api/v3/cryptocurrency/market-pairs/latest?slug=%s&start=1&limit=10&category=spot&centerType=all",
		strings.ToLower(req.CoinName))
	if config.Config.OpenNetProxy.OpenFlag == true {
		// Setting a Proxy URL and Port
		client.SetProxy("http://proxy.idchats.com:7890")
	}
	resp, err := client.R().Get(url)
	if err != nil || resp.StatusCode() != 200 {
		responseData.ErrCode = 501
		responseData.ErrMsg = "操作错误"
		c.JSON(http.StatusOK, responseData)
		return
	} else {
		logrus.Info(url)
		var Tdata TCoinMarketCap
		err := json.Unmarshal(resp.Body(), &Tdata)
		if err != nil {
			logrus.Info(resp.String())
			responseData.ErrCode = 502
			responseData.ErrMsg = "操作错误" + err.Error()
			c.JSON(http.StatusOK, responseData)
			return
		}
		if Tdata.Status.ErrorCode != "0" {
			responseData.ErrCode = utils2.StringToInt(Tdata.Status.ErrorCode)
			responseData.ErrMsg = Tdata.Status.ErrorMessage
			return
		}
		if len(Tdata.Data.MarketPairs) == 0 {
			responseData.ErrCode = 503
			responseData.ErrMsg = "操作错误"
			c.JSON(http.StatusOK, responseData)
			return
		}
		responseData.ErrCode = 0
		responseData.ErrMsg = "ok"
		nowValue := Tdata.Data.MarketPairs[0].Price
		for _, value := range Tdata.Data.MarketPairs {
			if value.Price < nowValue {
				nowValue = value.Price
			}
		}
		responseData.Data = nowValue

		c.JSON(http.StatusOK, responseData)
	}
	return
}

type TCoinMarketCap struct {
	Data struct {
		Id             int    `json:"id"`
		Name           string `json:"name"`
		Symbol         string `json:"symbol"`
		NumMarketPairs int    `json:"numMarketPairs"`
		MarketPairs    []struct {
			ExchangeId          int         `json:"exchangeId"`
			ExchangeName        string      `json:"exchangeName"`
			ExchangeSlug        string      `json:"exchangeSlug"`
			ExchangeNotice      string      `json:"exchangeNotice"`
			OutlierDetected     int         `json:"outlierDetected"`
			PriceExcluded       int         `json:"priceExcluded"`
			VolumeExcluded      int         `json:"volumeExcluded"`
			MarketId            int         `json:"marketId"`
			MarketPair          string      `json:"marketPair"`
			Category            string      `json:"category"`
			MarketUrl           string      `json:"marketUrl"`
			MarketScore         string      `json:"marketScore"`
			MarketReputation    interface{} `json:"marketReputation"`
			BaseSymbol          string      `json:"baseSymbol"`
			BaseCurrencyId      int         `json:"baseCurrencyId"`
			QuoteSymbol         string      `json:"quoteSymbol"`
			QuoteCurrencyId     int         `json:"quoteCurrencyId"`
			Price               float64     `json:"price"`
			VolumeUsd           float64     `json:"volumeUsd"`
			EffectiveLiquidity  float64     `json:"effectiveLiquidity"`
			Liquidity           interface{} `json:"liquidity"`
			LastUpdated         time.Time   `json:"lastUpdated"`
			Quote               float64     `json:"quote"`
			VolumeBase          float64     `json:"volumeBase"`
			VolumeQuote         float64     `json:"volumeQuote"`
			FeeType             string      `json:"feeType"`
			DepthUsdNegativeTwo float64     `json:"depthUsdNegativeTwo"`
			DepthUsdPositiveTwo float64     `json:"depthUsdPositiveTwo"`
			ReservesAvailable   interface{} `json:"reservesAvailable"`
			PorAuditStatus      interface{} `json:"porAuditStatus"`
		} `json:"marketPairs"`
	} `json:"data"`
	Status struct {
		Timestamp    time.Time `json:"timestamp"`
		ErrorCode    string    `json:"error_code"`
		ErrorMessage string    `json:"error_message"`
		Elapsed      string    `json:"elapsed"`
		CreditCount  int       `json:"credit_count"`
	} `json:"status"`
}
