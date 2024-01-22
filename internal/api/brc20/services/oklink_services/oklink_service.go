package oklink_services

import (
	"Open_IM/internal/api/brc20/services"
	"Open_IM/internal/utils"
	"Open_IM/pkg/common/config"
	"errors"
	"fmt"
	"strconv"
)

// Get brc20Balance-list
func MainNetGetAddressBrc20BalanceListResult(address, tick string, netParam string, page, limit int64) ([]*services.BalanceListItem, error) {
	var (
		url    string
		result string
		resp   *OklinkResp
		data   *OklinkBrc20BalanceList
		err    error
		query  map[string]string = map[string]string{
			"address": address,
			"token":   tick,
			"page":    strconv.FormatInt(page, 10),
			"limit":   strconv.FormatInt(limit, 10),
		}
		headers map[string]string = map[string]string{
			"Ok-Access-Key": config.Config.OklinkKey,
		}
	)

	url = fmt.Sprintf("%s/api/v5/explorer/btc/address-balance-list", config.Config.OklinkDomain)
	result, err = utils.GetUrl(url, query, headers)
	if err != nil {
		return nil, err
	}
	//fmt.Println(result)
	if err = utils.JsonToObject(result, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("Get request err:%s", err))
	}

	if resp.Code != "0" {
		return nil, errors.New(fmt.Sprintf("Msg:%s", resp.Msg))
	}

	if err = utils.JsonToAny(resp.Data, &data); err != nil {
		return nil, errors.New(fmt.Sprintf("Get request err:%s", err))
	}
	return data.BalanceList, nil
}

type OklinkResp struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type OklinkBrc20BalanceDetails struct {
	Page                string         `json:"page"`
	Limit               string         `json:"limit"`
	TotalPage           string         `json:"totalPage"`
	Token               string         `json:"token"`
	TokenType           string         `json:"tokenType"`
	Balance             string         `json:"balance"`
	AvailableBalance    string         `json:"availableBalance"`
	TransferBalance     string         `json:"transferBalance"`
	TransferBalanceList []*BalanceItem `json:"transferBalanceList"`
}

type BalanceItem struct {
	InscriptionId     string `json:"inscriptionId"`
	InscriptionNumber string `json:"inscriptionNumber"`
	Amount            string `json:"amount"`
}

type OklinkBrc20BalanceList struct {
	Page        string                      `json:"page"`
	Limit       string                      `json:"limit"`
	TotalPage   string                      `json:"totalPage"`
	BalanceList []*services.BalanceListItem `json:"balanceList"`
}
