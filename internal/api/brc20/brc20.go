package brc20

import (
	"Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/log"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/shopspring/decimal"
	"net/http"
)

// @Summary		获取Brc20Token列表
// @Description	Brc20
// @Tags		Brc20相关
// @ID			GetBrc20Tokens
// @Accept		json
// @Param		req		body	api.AllBrc20TokensReq{}	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.AllBrc20TokensResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/brc20/brc20_all_tokens [post]
func GetBrc20Tokens(c *gin.Context) {
	params := api.AllBrc20TokensReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	requestOkLink := &OKLinkBrc20AllTokenListReq{
		T:      int64(params.Timestamp),
		Limit:  params.PageSize,
		Type:   params.Type,
		Offset: params.PageIndex * params.PageSize,
		Sort:   params.Sort,
	}
	var responseOkLink OKLinkBrc20AllTokenListResp
	err := gout.GET("https://www.oklink.com/api/explorer/v2/btc/inscription/token/list").
		SetHeader(gout.H{
			"X-Apikey":     utils.GetOkLinkXAPIKey(),
			"Content-Type": "application/json",
		}).Debug(true).SetQuery(requestOkLink).BindJSON(&responseOkLink).Do()
	if err != nil {
		log.NewError(params.OperationID, "gout.GET failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if responseOkLink.Code != 0 {
		log.NewError(params.OperationID, "gout.GET failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": responseOkLink.DetailMsg})
		return
	}
	var response api.AllBrc20TokensResp
	response.CommResp.ErrCode = 0
	response.CommResp.ErrMsg = ""
	response.Data = new(api.HitsBrc20Response)
	response.Data.Total = responseOkLink.Data.Total
	for key := range responseOkLink.Data.Hits {
		var item api.ApiHitsBrc20
		item.Rank = responseOkLink.Data.Hits[key].Rank
		item.LogoUrl = responseOkLink.Data.Hits[key].LogoUrl
		item.Name = responseOkLink.Data.Hits[key].Name
		item.DisplayName = responseOkLink.Data.Hits[key].DisplayName
		supply := decimal.NewFromFloat(responseOkLink.Data.Hits[key].Supply.(float64))
		item.Supply = supply.String()
		mint := decimal.NewFromFloat(responseOkLink.Data.Hits[key].Minted.(float64))
		item.Minted = mint.String()
		item.MintRate = responseOkLink.Data.Hits[key].MintRate
		item.TransactionCount = responseOkLink.Data.Hits[key].TransactionCount
		item.HolderCount = responseOkLink.Data.Hits[key].HolderCount
		item.DeployTime = responseOkLink.Data.Hits[key].DeployTime
		item.InscriptionNumber = responseOkLink.Data.Hits[key].InscriptionNumber
		item.InscriptionId = responseOkLink.Data.Hits[key].InscriptionId
		item.Price = responseOkLink.Data.Hits[key].Price
		item.TickId = responseOkLink.Data.Hits[key].TickId
		item.TokenType = responseOkLink.Data.Hits[key].TokenType
		response.Data.Hits = append(response.Data.Hits, &item)
	}
	log.NewInfo(params.OperationID, "GetBrc20Tokens return ", response)
	c.JSON(http.StatusOK, response)
}

// https://www.oklink.com/api/explorer/v2/btc/inscription/list?t=1701079336235&offset=20&limit=20&type=BRC20&sortKey=0
// @Summary		查找指定用户的brc20
// @Description	Brc20
// @Tags		Brc20相关
// @ID			GetPersonalBrc20Tokens
// @Accept		json
// @Param		req		body	api.PersonalAllBrc20TokenReq{}	true	"请求体RequestType :balance，deploy ,inscription"
// @Produce		json
// @Success		0	{object}	api.PersonalAllBrc20TokenRsp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/brc20/brc20_personal_tokens [post]
func GetPersonalBrc20Tokens(c *gin.Context) {
	//获取自己个人的内容
	params := api.PersonalAllBrc20TokenReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	requestOkLink := &OKLinkBrc20PersonalReq{
		T:       int64(params.Timestamp),
		Limit:   params.PageSize,
		Type:    params.Type,
		Offset:  params.PageIndex * params.PageSize,
		Address: params.Address,
	}
	requestUrl := "https://www.oklink.com/api/explorer/v2/btc/inscription/token-balance/list"
	if params.RequestType == "deploy" {
		requestUrl = "https://www.oklink.com/api/explorer/v2/btc/inscription/token-deploy/list"
	} else if params.RequestType == "inscription" {
		requestUrl = "https://www.oklink.com/api/explorer/v2/btc/inscription/list"
	}

	var responseOkLink OKLinkBrc20PersonalResp
	err := gout.GET(requestUrl).SetHeader(gout.H{
		"X-Apikey":     utils.GetOkLinkXAPIKey(),
		"Content-Type": "application/json",
	}).SetQuery(requestOkLink).Debug(true).BindJSON(&responseOkLink).Do()
	if err != nil {
		log.NewError(params.OperationID, "gout.GET failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if responseOkLink.Code != 0 {
		log.NewError(params.OperationID, "gout.GET failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": responseOkLink.DetailMsg})
		return
	}
	var response api.PersonalAllBrc20TokenRsp
	response.CommResp.ErrCode = 0
	response.CommResp.ErrMsg = ""
	response.Data = new(api.HitsPersonalBrc20Response)
	response.Data.Total = responseOkLink.Data.Total
	response.Data.Extend = responseOkLink.Data.Extend
	for key := range responseOkLink.Data.Hits {
		var item api.HitsPersonalBrc20
		item.LogoUrl = responseOkLink.Data.Hits[key].LogoUrl
		item.Name = responseOkLink.Data.Hits[key].Name
		item.DisplayName = responseOkLink.Data.Hits[key].DisplayName
		item.Balance = responseOkLink.Data.Hits[key].Balance
		item.AvailableBalance = responseOkLink.Data.Hits[key].AvailableBalance
		item.TransferableBalance = responseOkLink.Data.Hits[key].TransferableBalance
		item.InscriptionId = responseOkLink.Data.Hits[key].InscriptionId
		item.InscriptionNumber = responseOkLink.Data.Hits[key].InscriptionNumber
		item.TokenInscriptionNumber = responseOkLink.Data.Hits[key].TokenInscriptionNumber
		item.InscriptionAmount = responseOkLink.Data.Hits[key].InscriptionAmount
		item.Price = responseOkLink.Data.Hits[key].Price
		item.UsdValue = responseOkLink.Data.Hits[key].UsdValue
		item.Type = responseOkLink.Data.Hits[key].Type
		item.DeployTime = responseOkLink.Data.Hits[key].DeployTime
		item.Action = responseOkLink.Data.Hits[key].Action
		item.Owner = responseOkLink.Data.Hits[key].Owner
		item.Status = responseOkLink.Data.Hits[key].Status
		item.ErrorMsg = responseOkLink.Data.Hits[key].ErrorMsg
		response.Data.Hits = append(response.Data.Hits, &item)
	}
	log.NewInfo(params.OperationID, "GetBrc20Tokens return ", response)
	c.JSON(http.StatusOK, response)
}

//https://www.oklink.com/api/explorer/v2/btc/inscription/token-deploy/list?t=1701136532842&offset=0&limit=20&address=bc1qn7vswjfwp0uvgqr3m4r43ukkuu2whl4xjr9ax6&type=BRC20

//https://www.oklink.com/api/explorer/v2/btc/inscription/token-balance/list?t=1701136536393&offset=0&limit=20&address=bc1qn7vswjfwp0uvgqr3m4r43ukkuu2whl4xjr9ax6&type=BRC20

type OKLinkBrc20AllTokenListResp struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	DetailMsg string `json:"detailMsg"`
	Data      struct {
		Total int              `json:"total"`
		Hits  []*api.HitsBrc20 `json:"hits"`
	} `json:"data"`
}
type OKLinkBrc20PersonalResp struct {
	Code      int                            `json:"code"`
	Msg       string                         `json:"msg"`
	DetailMsg string                         `json:"detailMsg"`
	Data      *api.HitsPersonalBrc20Response `json:"data"`
}
type OKLinkBrc20AllTokenListReq struct {
	T      int64  `query:"t"`
	Limit  int    `query:"limit"`
	Type   string `query:"type"`
	Offset int    `query:"offset"`
	Sort   string `query:"sort,omitempty"`
}

// t=1701136773487&offset=0&limit=20&address=bc1qn7vswjfwp0uvgqr3m4r43ukkuu2whl4xjr9ax6&type=BRC20
type OKLinkBrc20PersonalReq struct {
	T       int64  `query:"t"`
	Limit   int    `query:"limit"`
	Type    string `query:"type"`
	Offset  int    `query:"offset"`
	Address string `query:"address"`
}
