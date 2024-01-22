package unisat_wallet

import (
	"Open_IM/internal/api/brc20/services"
	"Open_IM/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/shopspring/decimal"
	"strings"
)

const (
	devnetRestUrl  = "https://fullnode.devnet.aptoslabs.com"
	testnetRestUrl = "https://testnet.aptoslabs.com"
	mainnetRestUrl = "https://fullnode.mainnet.aptoslabs.com"
)

func mempoolHost(chainnet string) string {
	switch chainnet {
	case "mainnet", "livenet":
		return "https://mempool.space/api"
	case "testnet":
		return "https://mempool.space/testnet/api"
	}
	return ""
}

func unisatHostApi(chainnet string) (string, error) {
	switch chainnet {
	case "mainnet", "livenet":
		return "https://open-api.unisat.io", nil
	case "testnet":
		return "https://open-api-testnet.unisat.io", nil
	}
	return "", errors.New("Unsupported BTC chainnet")
}
func unisatHost(chainnet string) (string, error) {
	switch chainnet {
	case "mainnet", "livenet":
		return "https://api.unisat.io", nil
	case "testnet":
		return "https://api-testnet.unisat.io", nil
	}
	return "", errors.New("Unsupported BTC chainnet")
}

func unisatRequestHeader(address string) map[string]string {
	return map[string]string{
		"X-Client":   "UniSat Wallet",
		"X-Version":  "1.1.33",
		"x-address":  address,
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	}
}
func NetTestGetAddressBrc20BalanceListResult(address, tick string, netParamStr string, page, limit int64) (resultArray []*services.BalanceListItem, err error) {
	host, err := unisatHost(netParamStr)
	if err != nil {
		return nil, err
	}
	header := unisatRequestHeader(address)
	url := fmt.Sprintf("%v/wallet-v4/brc20/tokens?address=%v&cursor=%v&size=%v", host, address, page, limit)
	result, err := utils.GetUrl(url, nil, header)
	if err != nil {
		return nil, err
	}
	var rawPage TUnisatWallResponse[UnisatWalletBalanceBrc20]
	if err = json.Unmarshal([]byte(result), &rawPage); err != nil {
		return nil, err
	}
	if rawPage.Status != "1" {
		return nil, errors.New(rawPage.Message)
	}
	for _, value := range rawPage.Result.List {
		resultArray = append(resultArray, &services.BalanceListItem{
			Token:            value.Ticker,
			TokenType:        "brc-20",
			Balance:          value.OverallBalance,
			AvailableBalance: value.AvailableBalance,
			TransferBalance:  value.TransferableBalance,
		})
	}
	return
}

// https://api-testnet.unisat.io/wallet-v4/address/btc-utxo?address=tb1q0e36uglptue2uk67qu5um42fq5emyadpnf654w
func GetUnspendUtxoFromUnisat(address string, netParamStr string) (resultUnspend []*services.UnspendUtxo, err error) {
	host, err := unisatHost(netParamStr)
	if err != nil {
		return nil, err
	}
	header := unisatRequestHeader(address)
	url := fmt.Sprintf("%v/wallet-v4/address/btc-utxo?address=%v", host, address)
	result, err := utils.GetUrl(url, nil, header)
	if err != nil {
		return nil, err
	}
	var rawPage TUnspendUtxoFromUnisat
	if err = json.Unmarshal([]byte(result), &rawPage); err != nil {
		return nil, err
	}
	if rawPage.Status != "1" {
		return nil, errors.New(rawPage.Message)
	}

	for _, value := range rawPage.Result {
		if value.Satoshis == 546 || value.Satoshis == 600 || value.Satoshis == 1000 {
			continue
		}
		resultUnspend = append(resultUnspend, value)
	}
	return resultUnspend, err
}

// https://api-testnet.unisat.io/wallet-v4/default/fee-summary
func GetFeeRate(address string, netParamStr string) (feeRate int64, err error) {
	host, err := unisatHost(netParamStr)
	if err != nil {
		return 0, err
	}
	header := unisatRequestHeader(address)
	url := fmt.Sprintf("%v/wallet-v4/default/fee-summary", host)
	result, err := utils.GetUrl(url, nil, header)
	if err != nil {
		return 0, err
	}
	var rawPage TUnisatWallResponse[FeeRate]
	if err = json.Unmarshal([]byte(result), &rawPage); err != nil {
		return 0, err
	}
	if rawPage.Status != "1" {
		return 0, errors.New(rawPage.Message)
	}
	for _, value := range rawPage.Result.List {
		if strings.EqualFold(value.Title, "avg") {
			return int64(value.FeeRate), nil
		}
	}
	return 0, errors.New("无法找到最合适的费率")
}

func CheckIsHaveThisUtxo(address string, tick string, useInstrip []string, netParamStr string) (Amount string, amountList []string, err error) {
	//https://api-testnet.unisat.io/wallet-v4/brc20/token-summary?address=tb1q0e36uglptue2uk67qu5um42fq5emyadpnf654w&ticker=dead
	host, err := unisatHost(netParamStr)
	if err != nil {
		return "", nil, err
	}
	header := unisatRequestHeader(address)
	url := fmt.Sprintf("%v/wallet-v4/brc20/token-summary?address=%v&ticker=%v", host, address, tick)
	result, err := utils.GetUrl(url, nil, header)
	if err != nil {
		return "", nil, err
	}
	var rawPage TTokenSummary
	if err = json.Unmarshal([]byte(result), &rawPage); err != nil {
		return "", nil, err
	}
	if rawPage.Status != "1" {
		return "", nil, errors.New(rawPage.Message)
	}
	var containTransferable map[string]string
	containTransferable = make(map[string]string, 0)
	for _, value := range rawPage.Result.TransferableList {
		containTransferable[value.InscriptionId] = value.Amount
	}
	returnAmount := decimal.Zero
	for _, value := range useInstrip {
		if valueAmount, ok := containTransferable[value]; ok {
			decimalAmount, _ := decimal.NewFromString(valueAmount)
			amountList = append(amountList, decimalAmount.String())
			returnAmount = returnAmount.Add(decimalAmount)
		} else {
			return "", nil, errors.New("提交的票据有问题")
		}
	}

	return returnAmount.String(), amountList, nil
}
func NetGetAddressBrc20TransableListResult(address, tick string, netParamStr string, page, limit int64) (resultArray []*services.TransferAbleInscript, err error) {
	host, err := unisatHost(netParamStr)
	if err != nil {
		return nil, err
	}
	header := unisatRequestHeader(address)
	url := fmt.Sprintf("%v/wallet-v4/brc20/transferable-list?address=%v&ticker=%v&cursor=%v&size=%v", host, address, tick, page, limit)
	fmt.Println(url)
	result, err := utils.GetUrl(url, nil, header)
	if err != nil {
		return nil, err
	}
	var rawPage TUnisatWallResponse[TUnisatWalletTransableBrc20]
	if err = json.Unmarshal([]byte(result), &rawPage); err != nil {
		return nil, err
	}
	if rawPage.Status != "1" {
		return nil, errors.New(rawPage.Message)
	}
	for _, value := range rawPage.Result.List {
		voutIndex, _ := convertor.ToInt(strings.Split(value.InscriptionId, "i")[1])
		resultArray = append(resultArray, &services.TransferAbleInscript{
			InscriptionId:     value.InscriptionId,
			Ticker:            value.Ticker,
			Amount:            value.Amount,
			InscriptionNumber: value.InscriptionNumber,
			UtxoHash:          strings.Split(value.InscriptionId, "i")[0],
			Vout:              int(voutIndex),
		})
	}
	return
}

// Ticker:                 value.Ticker,
// OverallBalance:         value.OverallBalance,
// TransferableBalance:    value.TransferableBalance,
// AvailableBalance:       value.AvailableBalance,
// AvailableBalanceSafe:   value.AvailableBalanceSafe,
// AvailableBalanceUnSafe: value.AvailableBalanceUnSafe,
// Decimal:                value.Decimal,
type TUnisatWallResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		List  []*T `json:"list"`
		Total int  `json:"total"`
	} `json:"result"`
}
type FeeRate struct {
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	FeeRate int    `json:"feeRate"`
}
type TTokenSummary struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		TokenBalance struct {
			Ticker                 string `json:"ticker"`
			AvailableBalance       string `json:"availableBalance"`
			TransferableBalance    string `json:"transferableBalance"`
			OverallBalance         string `json:"overallBalance"`
			AvailableBalanceSafe   string `json:"availableBalanceSafe"`
			AvailableBalanceUnSafe string `json:"availableBalanceUnSafe"`
		} `json:"tokenBalance"`
		HistoryList []struct {
			InscriptionId     string `json:"inscriptionId"`
			InscriptionNumber int    `json:"inscriptionNumber"`
			Amount            string `json:"amount"`
			Ticker            string `json:"ticker"`
		} `json:"historyList"`
		TransferableList []struct {
			InscriptionId     string `json:"inscriptionId"`
			InscriptionNumber int    `json:"inscriptionNumber"`
			Amount            string `json:"amount"`
			Ticker            string `json:"ticker"`
		} `json:"transferableList"`
		TokenInfo struct {
			TotalSupply string `json:"totalSupply"`
			TotalMinted string `json:"totalMinted"`
		} `json:"tokenInfo"`
	} `json:"result"`
}

type TUnspendUtxoFromUnisat struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Result  []*services.UnspendUtxo `json:"result"`
}
type UnisatWalletBalanceBrc20 struct {
	Ticker                 string `json:"ticker"`
	OverallBalance         string `json:"overallBalance"`
	TransferableBalance    string `json:"transferableBalance"`
	AvailableBalance       string `json:"availableBalance"`
	AvailableBalanceSafe   string `json:"availableBalanceSafe"`
	AvailableBalanceUnSafe string `json:"availableBalanceUnSafe"`
	Decimal                int    `json:"decimal"`
}
type TUnisatWalletTransableBrc20 struct {
	InscriptionId     string `json:"inscriptionId"`
	Ticker            string `json:"ticker"`
	Amount            string `json:"amount"`
	InscriptionNumber int    `json:"inscriptionNumber"`
}
