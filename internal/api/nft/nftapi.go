package nft

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/msg"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	rpc "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sourcegraph/conc"
)

// GetShowBalance
// @Summary		展示用户的资产
// @Description	展示用户的资产
// @Tags			用户资产相关
// @ID				GetShowBalance
// @Accept			json
// @Param			req		body	api.GetGlobalUserProfileReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UserBalanceInfoResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/nft/get_show_balance [post]
// d45a5575ca5249a4bd2c228e11e8aafe
func GetShowBalance(c *gin.Context) {
	var (
		req  api.GetGlobalUserProfileReq
		resp api.UserBalanceInfoResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var wg conc.WaitGroup
	var resultBackFromOkLink []*api.UserBalanceInfoItems
	resp.UserBalanceInfoDetail = new(api.UserBalanceInfoDetail)
	resp.UserBalanceInfoDetail.Items = make([]*api.UserBalanceInfoItems, 0)
	if req.ChainNet != "bsc-mainnet" {
		wg.Go(func() {
			resultbyte, err := utils.HttpGet(fmt.Sprintf("https://open-platform.nodereal.io/%s/covalenthq/v1/%s/address/%s/balances_v2/",
				"d45a5575ca5249a4bd2c228e11e8aafe", req.ChainNet, req.UserID))
			if err != nil {
				resp.ErrCode = constant.ErrInternal.ErrCode
				resp.ErrMsg = "balance error"
				log.NewError(req.OperationID, "bind req", utils.GetSelfFuncName(), err.Error())

			} else if len(resultbyte) > 0 {
				json.Unmarshal(resultbyte, &resp)
				if resp.Error {
					resp.ErrCode = constant.ErrInternal.ErrCode
					resp.UserBalanceInfoDetail = nil
					_, ok := resp.ErrorMessage.(string)
					if ok {
						resp.ErrMsg = resp.ErrorMessage.(string)
					} else {
						resp.ErrMsg = "数据错误"
					}
				}
			}
		})
	}
	if req.ChainNet == "bsc-mainnet" {
		resp.UserBalanceInfoDetail.ChainID = 56
		wg.Go(func() {
			resultBackFromOkLink = GetOKLinkTokenDetail(req.UserID, "bsc")
		})
	}
	wg.Wait()
	if len(resultBackFromOkLink) > 0 {
		contain := map[string]struct{}{}
		for _, value := range resp.UserBalanceInfoDetail.Items {
			contain[strings.ToLower(value.ContractAddress)] = struct{}{}
		}
		for index := 0; index < len(resultBackFromOkLink); index++ {
			if _, ok := contain[resultBackFromOkLink[index].ContractAddress]; !ok {
				resp.UserBalanceInfoDetail.Items = append(resp.UserBalanceInfoDetail.Items, resultBackFromOkLink[index])
			}
		}
	}
	c.JSON(http.StatusOK, resp)
}

// GetBtcShowBalance
// @Summary		展示某个BTC钱包地址资产
// @Description	展示某个BTC钱包地址资产
// @Tags			用户资产相关
// @ID				GetBtcShowBalance
// @Accept			json
// @Param			req		body	api.GetGlobalUserProfileReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UserBtcBalanceInfoResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/nft/get_show_btc_balance [post]
func GetBtcShowBalance(c *gin.Context) {
	var (
		req  api.GetGlobalUserProfileReq
		resp api.UserBtcBalanceInfoResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var BtcBalance []byte
	var err error
	proxyAddress, _ := url.Parse("http://proxy.idchats.com:7890")
	resp.Data = new(api.TBtcWalletDetail)
	var wg conc.WaitGroup
	wg.Go(func() {
		if !config.Config.IsPublicEnv {
			BtcBalance, err = utils.HttpGetWithHeaderWithProxy(fmt.Sprintf("https://unisat.io/wallet-api-v4/address/balance?address=%s", req.UserID), map[string]string{
				"X-Channel": "store",
				"X-Client":  "UniSat Wallet",
				"X-Udid":    randomdata.RandStringRunes(12),
				"X-Version": "1.1.25",
			}, http.ProxyURL(proxyAddress))
		} else {
			BtcBalance, err = utils.HttpGetWithHeader(fmt.Sprintf("https://unisat.io/wallet-api-v4/address/balance?address=%s", req.UserID), map[string]string{
				"X-Channel": "store",
				"X-Client":  "UniSat Wallet",
				"X-Udid":    randomdata.RandStringRunes(12),
				"X-Version": "1.1.25",
			})
		}
		if len(BtcBalance) > 0 {
			type uniSatStruct struct {
				Status  string                `json:"status"`
				Message string                `json:"message"`
				Result  *api.TBtcWalletDetail `json:"result"`
			}
			var unData uniSatStruct
			err = json.Unmarshal(BtcBalance, &unData)
			if err == nil && unData.Status == "1" {
				resp.Data.ConfirmAmount = unData.Result.ConfirmAmount
				resp.Data.PendingAmount = unData.Result.PendingAmount
				resp.Data.Amount = unData.Result.Amount
				resp.Data.ConfirmBtcAmount = unData.Result.ConfirmBtcAmount
				resp.Data.PendingBtcAmount = unData.Result.PendingBtcAmount
				resp.Data.BtcAmount = unData.Result.BtcAmount
				resp.Data.ConfirmInscriptionAmount = unData.Result.ConfirmInscriptionAmount
				resp.Data.PendingInscriptionAmount = unData.Result.PendingInscriptionAmount
				resp.Data.InscriptionAmount = unData.Result.InscriptionAmount
				resp.Data.UsdValue = unData.Result.UsdValue
			}
		}
	})
	wg.Go(func() {
		var BtcBalanceData []byte
		if !config.Config.IsPublicEnv {
			BtcBalance, err = utils.HttpGetWithHeaderWithProxy(fmt.Sprintf("https://unisat.io/wallet-api-v4/brc20/tokens?address=%s&cursor=0&size=100", req.UserID), map[string]string{
				"X-Channel": "store",
				"X-Client":  "UniSat Wallet",
				"X-Udid":    randomdata.RandStringRunes(12),
				"X-Version": "1.1.25",
			}, http.ProxyURL(proxyAddress))
		} else {
			BtcBalance, err = utils.HttpGetWithHeader(fmt.Sprintf("https://unisat.io/wallet-api-v4/brc20/tokens?address=%s&cursor=0&size=100", req.UserID), map[string]string{
				"X-Channel": "store",
				"X-Client":  "UniSat Wallet",
				"X-Udid":    randomdata.RandStringRunes(12),
				"X-Version": "1.1.25",
			})
		}

		if len(BtcBalanceData) > 0 {
			type uniSatStruct struct {
				Status  string `json:"status"`
				Message string `json:"message"`
				Result  struct {
					List  []api.TBtcWalletDetailTokenList `json:"list"`
					Total int                             `json:"total"`
				} `json:"result"`
			}
			var unData uniSatStruct
			err := json.Unmarshal(BtcBalanceData, &unData)
			if err == nil && unData.Status == "1" {
				for key, _ := range unData.Result.List {
					resp.Data.List = append(resp.Data.List, &api.TBtcWalletDetailTokenList{
						Ticker:              unData.Result.List[key].Ticker,
						OverallBalance:      unData.Result.List[key].OverallBalance,
						TransferableBalance: unData.Result.List[key].TransferableBalance,
						AvailableBalance:    unData.Result.List[key].AvailableBalance,
					})
				}
				resp.Data.Total = unData.Result.Total
			}
		}
	})
	// panics with a nice stacktrace
	wg.Wait()

	c.JSON(http.StatusOK, resp)
}

// UpdateLikeShowNft
// @Summary		为NFT点赞
// @Description	为NFT点赞
// @Tags			用户相关
// @ID				UpdateLikeShowNft
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.LikeActionNftReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.LikeActionNftResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/nft/like_unlike_nft_id [post]
func UpdateLikeShowNft(c *gin.Context) {
	var (
		req   api.LikeActionNftReq
		resp  api.LikeActionNftResp
		reqPb pbChat.SendLikeMsgReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.MsgData = new(server_api_params.LikeRewordReq)
	utils.CopyStructFields(reqPb.MsgData, &req)
	reqPb.MsgData.ContentType = "nft"
	var ok bool
	var errInfo string
	ok, reqPb.MsgData.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImMsgName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbChat.NewMsgClient(etcdConn)
	respPb, err := client.SendLikeAction(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.ErrCode
	resp.CommResp.ErrMsg = respPb.ErrMsg
	c.JSON(http.StatusOK, resp)
	return
}

// GetLikeShowNft
// @Summary		展示NFT的内容
// @Description	展示NFT的内容
// @Tags			用户相关
// @ID				GetLikeShowNft
// @Accept			json
// @Param			req		body	api.GetLikeShowNftCountReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetShowNftResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/nft/get_like_status_nft_id [post]
func GetLikeShowNft(c *gin.Context) {
	var (
		req   api.GetLikeShowNftCountReq
		resp  api.LikeActionNftResp
		reqPb rpc.RpcLikeShowNftStatusReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.ArticleID = req.ArticleID
	_, reqPb.UserID, _ = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.GetShowNftLikeStatus(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = 0
	resp.CommResp.ErrMsg = ""

	resp.NCount = respPb.NftLikeCount
	resp.IsLike = respPb.NftIsLike
	c.JSON(http.StatusOK, resp)
	return
}
