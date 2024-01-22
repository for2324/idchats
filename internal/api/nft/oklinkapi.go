package nft

import (
	"Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/shopspring/decimal"
	"time"
)

type TOkLinkTokenHits struct {
	HolderAddress        string  `json:"holderAddress"`
	TokenContractAddress string  `json:"tokenContractAddress"`
	IsAddressRisk        bool    `json:"isAddressRisk"`
	Value                float64 `json:"value"`
	Price                float64 `json:"price"`
	EthValue             float64 `json:"ethValue"`
	UsdValue             float64 `json:"usdValue"`
	Symbol               string  `json:"symbol"`
	CoinName             string  `json:"coinName"`
	IsRiskStablecoin     bool    `json:"isRiskStablecoin"`
	IsRiskToken          bool    `json:"isRiskToken"`
	LogoUrl              string  `json:"logoUrl"`
	Rate                 float64 `json:"rate"`
	PercentChange24H     float64 `json:"percentChange24h"`
	ValueChangeRate24H   float64 `json:"valueChangeRate24h"`
}
type TOKLinkApiResponse struct {
	Code      int                     `json:"code"`
	Msg       string                  `json:"msg"`
	DetailMsg string                  `json:"detailMsg"`
	Data      *TOKLinkApiResponseData `json:"data"`
}
type TOKLinkApiResponseData struct {
	Total  int                 `json:"total"`
	Hits   []*TOkLinkTokenHits `json:"hits"`
	Extend struct {
		Total      int     `json:"total"`
		ValueTotal float64 `json:"valueTotal"`
	} `json:"extend"`
}

func GetOKLinkTokenDetail(requestAddress string, chianName string) (userBalanceInfoItems []*api.UserBalanceInfoItems) {
	//https://www.oklink.com/api/explorer/v2/bsc/addresses/0xa3dd654c68d0bf1d9b3495855d5f89d636877e0a/holders/token?t=1697765972553&offset=0&limit=20&tokenAddress=
	requestUrl := fmt.Sprintf("https://www.oklink.com/api/explorer/v2/%s/addresses/%s/holders/token", chianName, requestAddress)
	var totalTokenHits []*TOkLinkTokenHits //总数
	fromOffsetIndex := 0
	limitPageSize := 20
	for {
		var resp TOKLinkApiResponse
		NewRequstJson := struct {
			T            int64  `form:"t"`
			Offset       int    `form:"offset"`
			Limit        int    `form:"limit"`
			TokenAddress string `form:"tokenAddress"`
		}{
			T:            time.Now().UnixMilli(),
			Offset:       fromOffsetIndex,
			Limit:        limitPageSize,
			TokenAddress: "",
		}
		err2 := gout.GET(requestUrl).
			SetHeader(gout.H{
				"X-Apikey":     utils.GetOkLinkXAPIKey(),
				"Content-Type": "application/json",
			}).SetQuery(NewRequstJson).BindJSON(&resp).Do()
		if err2 != nil {
			break
		}
		if resp.Code != 0 || resp.Data == nil || len(resp.Data.Hits) == 0 {
			break
		}
		for _, value := range resp.Data.Hits {
			totalTokenHits = append(totalTokenHits, value)
		}

		//下一页的token信息
		if resp.Data != nil && len(resp.Data.Hits) < limitPageSize {
			break
		}
		if resp.Data != nil && len(resp.Data.Hits) == limitPageSize {
			fromOffsetIndex += limitPageSize
		}
	}
	for _, value := range totalTokenHits {
		if value.TokenContractAddress == "" {
			value.TokenContractAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		userBalanceInfoItems = append(userBalanceInfoItems, &api.UserBalanceInfoItems{
			ContractDecimals:     0,
			ContractName:         value.CoinName,
			ContractTickerSymbol: value.Symbol,
			ContractAddress:      value.TokenContractAddress,
			SupportsErc:          nil,
			LogoURL:              value.LogoUrl,
			LastTransferredAt:    time.Now(),
			NativeToken:          false,
			Type:                 "",
			Balance:              decimal.NewFromFloat(value.Value).String(),
			Balance24H:           decimal.NewFromFloat(value.Value).String(),
			QuoteRate:            value.Price,
			QuoteRate24H:         value.Price,
			Quote:                value.Price,
			Quote24H:             value.Price,
			NftData:              nil,
		})
	}
	return
}
