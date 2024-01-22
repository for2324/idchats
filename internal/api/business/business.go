package business

import (
	"Open_IM/internal/contract"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/utils"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"math/big"
	"net/http"
	"strings"
	"time"
)

// @Summary		获取服务商的开发者密钥
// @Description	服务商
// @Tags		服务商相关
// @ID			GetBusinessList
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param		req		body	api.UserGetBusinessListReq{}	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UserGetBusinessListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/business/get_business_list [post]
func GetBusinessList(c *gin.Context) {
	params := api.UserGetBusinessListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool

	var errInfo string
	var UserId string

	ok, UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, "AppointmentReq args ", params)
	var response api.UserGetBusinessListResp
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_robot_api").Where("user_id=? and status=1", UserId).Count(&response.Data.Total).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "can't get total size api count"})
		return
	}
	//同步费率
	if response.Data.Total > 0 {
		var resultData []*db.UserRobotAPI
		err := db.DB.MysqlDB.DefaultGormDB().Table("user_robot_api").Where("user_id=? and status=1", UserId).
			Offset(params.PageSize * params.PageIndex).
			Limit(params.PageSize).Order("created_at desc").Find(&resultData).Error
		if err != nil {
			log.NewError(params.OperationID, "GetBusinessList failed ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
			return
		}
		var keyArray []string
		for _, value := range resultData {
			keyArray = append(keyArray, value.APIKey)
		}
		syncApikeyTradeMap := GetApiKeyTradeVolume("", keyArray)
		for _, value := range resultData {
			tradeValue, ok := syncApikeyTradeMap[value.APIKey]
			if !ok {
				tradeValue = decimal.NewFromInt(0)
			}
			response.Data.ApiKeyInfo = append(response.Data.ApiKeyInfo, &api.BusinessApiKeyInfo{
				TradeFee:    value.TradeFee,
				SniperFee:   value.SniperFee,
				Key:         value.APIKey,
				ApiName:     value.APIName,
				TradeVolume: tradeValue.String(),
				CreatedAt:   utils.Int64ToString(value.CreatedAt.Unix()),
			})
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, response)
	return

}
func GetApiKeyTradeVolume(userID string, keyArray []string) (result map[string]decimal.Decimal) {
	result = make(map[string]decimal.Decimal, 0)
	resultDataFromBiBot, err := rocksCache.RequestGetTotalBuyApiTrade(userID, userID, keyArray)
	if err == nil && len(resultDataFromBiBot.Data) > 0 {
		for _, value := range resultDataFromBiBot.Data {
			if value.ApiKey != "" {
				if _, ok := result[value.ApiKey]; ok {
					result[value.ApiKey] = result[value.ApiKey].Add(decimal.NewFromInt(value.SellUsdPrice))
				} else {
					result[value.ApiKey] = decimal.NewFromInt(value.SellUsdPrice)
				}
			}
		}
	}
	return result
}

func SyncNewUserFeeRate(pageIndex int, merchantId string, userID string) (feeRateBibot *rocksCache.UserBiBotFeeRate) {
	feeRateBibot = new(rocksCache.UserBiBotFeeRate)
	feeRateBibot.TradeFeeRate = 0.002
	feeRateBibot.SniperFeeRate = 0.01
	if pageIndex == 0 {
		rocksCache.DeleteUserBiBotFeeRate(userID, userID)
	}
	feeRateBibotNow, err := rocksCache.GetUserBiBotFeeRate(userID, userID)
	if err == nil {
		feeRateBibot.TradeFeeRate = feeRateBibotNow.TradeFeeRate
		feeRateBibot.SniperFeeRate = feeRateBibotNow.SniperFeeRate
	}
	return feeRateBibot
}

type TQueryForBotServiceReq struct {
	MerchantUid string   `form:"merchantUid"`
	ApiKey      []string `form:"apiKey"`
	MerchantId  string   `form:"merchantId"`
}
type TQueryForBotServiceRespData struct {
	Type           string `json:"type"`
	SellAmount     int    `json:"sellAmount"`
	SellUsdPrice   int64  `json:"sellUsdPrice"`
	FeeUsdPrice    int    `json:"feeUsdPrice"`
	ActualSlippage int    `json:"actualSlippage"`
	ApiKey         string `json:"apiKey"`
}
type TQueryForBotServiceResp struct {
	Data    []*TQueryForBotServiceRespData `json:"data"`
	ErrCode int                            `json:"errCode"`
	ErrMsg  string                         `json:"errMsg"`
}

// @Summary		添加开发者密钥
// @Description	服务商
// @Tags  		服务商相关
// @ID			AddBusinessApiKey
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param		req		body	api.UserAddBusinessApiReq{}	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UserAddBusinessApiResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/business/add_business_api_key [post]
func AddBusinessApiKey(c *gin.Context) {
	params := api.UserAddBusinessApiReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	var UserId string
	ok, UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	mutexname := "AddKeyApi:" + UserId
	rs := db.DB.Pool
	mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
	ctx := context.Background()
	if err := mutex.LockContext(ctx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "正在生成密钥"})
		return
	}
	defer mutex.UnlockContext(ctx)
	if utils.VerifySignature(UserId, params.Sign, fmt.Sprintf("apiName=%s&key=%s", params.ApiName, params.Key)) {

		countIsHaveThisName := int64(0)
		db.DB.MysqlDB.DefaultGormDB().Table("user_robot_api").Where("user_id=? and api_name=?  and status=1 ", UserId, params.ApiName).
			Count(&countIsHaveThisName)
		if countIsHaveThisName >= 1 {
			c.JSON(http.StatusOK, gin.H{"errCode": 811, "errMsg": "have this name key"})
			return
		}
		//签名通过的情况
		err := db.DB.MysqlDB.DefaultGormDB().Table("user_robot_api").Create(&db.UserRobotAPI{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
			UserID:      UserId,
			APIKey:      params.Key,
			APISecret:   params.Sign,
			TradeVolume: "0",
			TradeFee:    convertor.ToString(0.0018),
			SniperFee:   convertor.ToString(0.009),
			Status:      1,
			APIName:     params.ApiName,
		}).Error
		if err != nil {
			errMsg := err.Error()
			log.NewError(params.OperationID, errMsg)
			c.JSON(http.StatusOK, gin.H{"errCode": 811, "errMsg": errMsg})
			return
		}
		c.JSON(http.StatusOK, api.CommResp{
			ErrCode: 0,
			ErrMsg:  "success",
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "签名不通过"})
		return
	}
}

// @Summary		管理开发者密钥
// @Description	服务商
// @Tags  		服务商相关
// @ID			UpdateBusinessApiKey
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param		req		body	api.UserUpdateBusinessApiReq{}	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UserUpdateBusinessApiResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/business/update_business_api_key [post]
func UpdateBusinessApiKey(c *gin.Context) {
	params := api.UserUpdateBusinessApiReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	var UserId string
	ok, UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	mutexname := "AddKeyApi:" + UserId
	rs := db.DB.Pool
	mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
	ctx := context.Background()
	if err := mutex.LockContext(ctx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "正在生成密钥"})
		return
	}
	defer mutex.UnlockContext(ctx)
	if !utils.VerifySignature(UserId, params.Sign, fmt.Sprintf("apiName=%s&key=%s", params.ApiName, params.Key)) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "签名不通过"})
		return
	}
	switch params.Method {
	case "update":
		err := db.DB.MysqlDB.DefaultGormDB().Table("user_robot_api").Where("user_id=? and api_key=?", UserId, params.Key).
			Updates(map[string]interface{}{"api_name": params.ApiName}).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
			return
		}
	case "delete":
		err := db.DB.MysqlDB.DefaultGormDB().Table("user_robot_api").Where("user_id=? and api_key=? and api_name=?", UserId, params.Key, params.ApiName).
			Updates(map[string]interface{}{"status": 0}).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "update ok"})
	return
}

// @Summary		获取用户积分奖励
// @Description	服务商
// @Tags  		服务商相关
// @ID			GetUserBusinessTrade
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param		req		body	api.GetUserTradeRewardScoreReq{}	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetUserTradeRewardScoreResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/business/get_user_trade_score [post]
func GetUserBusinessTrade(c *gin.Context) {
	params := api.GetUserTradeRewardScoreReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	var UserId string
	ok, UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	resultResp := new(api.GetUserTradeRewardScoreResp)
	resultResp.Data = new(api.GetUserTradeRewardScoreData)
	err := db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("invitation_code=?", UserId).Count(&resultResp.Data.InviteCount).Error
	var dbData db.UserHistoryTotal
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_history_total").Where("user_id=?", UserId).First(&dbData).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resultResp.Data.RewardFee = decimal.NewFromFloat(rocksCache.GetCurrentTradeParentRewardFee(dbData.SubTotalTradeVolume)).String()
	pendingDeciaml, _ := decimal.NewFromString(dbData.Pending)
	if config.Config.RewardTradeByScore == "" || config.Config.RewardTradeByScore == "0" {
		config.Config.RewardTradeByScore = "0.1"
	}
	bbtPrice, _ := decimal.NewFromString(config.Config.RewardTradeByScore)
	resultResp.Data.Pending = pendingDeciaml.Div(bbtPrice).String()
	pendingDeciaml, _ = decimal.NewFromString(dbData.RakebackPending)
	resultResp.Data.RakebackPending = pendingDeciaml.Div(bbtPrice).String()
	resultResp.Data.Claim = dbData.Claimed
	resultResp.Data.RakebackClaim = dbData.RakebackClaimed
	resultResp.Data.PersonalTradeVolume = dbData.TotalTradeVolume
	resultResp.Data.SubPersonalTradeFee = dbData.SubTotalTradeVolume
	c.JSON(http.StatusOK, resultResp)
	return
}

// @Summary		领取积分
// @Description	服务商
// @Tags  		服务商相关
// @ID			ClaimSelfReward
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param		req		body	api.GetUserTradeRewardScoreReq{}	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UserStakeRewardResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/business/claim_self_reward [post]
func ClaimSelfReward(c *gin.Context) {
	params := api.GetUserTradeRewardScoreReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	var UserId string
	ok, UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	UserId = strings.ToLower(UserId)
	if params.RewardType != "tradeReward" && params.RewardType != "rakebackReward" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "reward type error"})
		return
	}

	mutexname := "trade_volume:" + UserId
	rs := db.DB.Pool
	mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
	if err := mutex.LockContext(context.Background()); err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": "正在更新"})
		return
	}
	defer mutex.UnlockContext(context.Background())
	var dbUserClaim db.UserHistoryTotal
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_history_total").Where("user_id=?", UserId).First(&dbUserClaim).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if params.RewardType == "tradeReward" && dbUserClaim.Pending == "0" {
		c.JSON(http.StatusOK, gin.H{"errCode": 801, "errMsg": "not have tradeReward pending score"})
		return
	}
	if params.RewardType == "rakebackReward" && dbUserClaim.RakebackPending == "0" {
		c.JSON(http.StatusOK, gin.H{"errCode": 801, "errMsg": "not have rakeback pending score"})
		return
	}
	resp := new(api.UserStakeRewardResp)
	nonce, err, _ := GetCurrentUserNonce(UserId, &dbUserClaim)
	if err != nil {
		resp.CommResp.ErrCode = 803
		resp.CommResp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	bytedata, _ := hex.DecodeString(config.Config.RewardKey)
	key := []byte{0x01, 0x02, 0x03, 0x00, 0x30, 0x31, 0x32, 0x33, 0x01, 0x02, 0x03, 0x00, 0x30, 0x31, 0x32, 0x33}
	decrypted := cryptor.AesEcbDecrypt(bytedata, key)
	//解析密钥:
	scannerKey, err := crypto.HexToECDSA(utils.Bytes2string(decrypted))
	if err != nil {
		resp.CommResp.ErrCode = 803
		resp.CommResp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	chainIDint64 := big.NewInt(1)
	amountBigInt, _ := decimal.NewFromString(dbUserClaim.Pending)
	if params.RewardType == "rakebackReward" {
		amountBigInt, _ = decimal.NewFromString(dbUserClaim.RakebackPending)
	}

	bbtPrice, _ := decimal.NewFromString(config.Config.RewardTradeByScore) //bbt的币价按照0.1u算的情况下
	amountBigInt = amountBigInt.Div(bbtPrice)
	amountBigInt = amountBigInt.Shift(18)
	strCustom := params.RewardType + "&" + bbtPrice.String()
	_, signByte, err := contract.SignTransfer(scannerKey, common.HexToAddress(config.Config.RewardChainContractAddress),
		chainIDint64.Int64(), &contract.Eip712Transfer{
			Amount:    amountBigInt.BigInt(),
			Recipient: common.HexToAddress(UserId),
			Nonce:     nonce,
			Custom:    strCustom,
		},
	)
	resp.Data = new(api.SignData)
	resp.Data.Nonce = nonce.String()
	resp.Data.Amount = amountBigInt.String()
	resp.Data.Recipient = common.HexToAddress(UserId).String()
	resp.Data.SignData = common.Bytes2Hex(signByte)
	resp.Data.Custom = strCustom
	c.JSON(http.StatusOK, resp)
	return
}

// @Summary		获取子成员的交易总量
// @Description	获取子成员的交易总量
// @Tags  		服务商相关
// @ID			GetSubPersonalTotalVolume
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param		req		body	api.UserPersonalTotalVolumeListReq{}	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UserPersonalTotalVolumeListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/business/get_total_volume_sub_personal [post]
func GetSubPersonalTotalVolume(c *gin.Context) {
	params := api.UserPersonalTotalVolumeListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	var UserId string
	ok, UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	total, err1 := GetTotalUserCount(UserId, time.Now())
	totalDataDetail, err2 := GetTotalUserPersonDetail(UserId, params.PageSize, params.PageIndex, params.RequestTime)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err1.Error()})
		return
	}
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err2.Error()})
		return
	}
	resp := new(api.UserPersonalTotalVolumeListResp)
	resp.Data = new(api.UserPersonalTotalVolumeListRespList)
	resp.Data.Total = total
	for key := range totalDataDetail {
		resp.Data.Data = append(resp.Data.Data, &api.UserPersonalVolumeListData{
			ID:             totalDataDetail[key].ID,
			FinishTime:     totalDataDetail[key].FinishTime.Unix(),
			Address:        totalDataDetail[key].UserID,
			UsdTradeVolume: totalDataDetail[key].UsdTradeVolume,
		})
	}
	c.JSON(http.StatusOK, resp)
	return
}

type DBUserPersonalVolume struct {
	ID             int64     `json:"ID"`
	FinishTime     time.Time `json:"finishTime"`
	UserID         string    `json:"userID"`
	UsdTradeVolume string    `json:"usdTradeVolume"`
}

func GetTotalUserPersonDetail(invaliteCodeUser string, pageCount, pageIndex int, requestTime int64) (result []*DBUserPersonalVolume, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_history_reward").
		Joins("left join registers on registers.user_id = user_history_reward.user_id").
		Where("registers.invitation_code = ? and user_history_reward.finish_time <= ?", invaliteCodeUser, utils.UnixSecondToTime(requestTime).Format(time.DateTime)).
		Select("user_history_reward.id,user_history_reward.finish_time,user_history_reward.user_id,user_history_reward.usd_trade_volume").
		Order("user_history_reward.finish_time desc").Limit(pageCount).Offset(pageCount * pageIndex).Find(&result).Error
	return
}
func GetTotalUserCount(invaliteCodeUser string, time2 time.Time) (totalCount int64, err error) {
	//select c.* from user_history_reward c
	//	left JOIN registers  on registers.user_id = c.user_id
	//	where registers.invitation_code = '0xad9ebc6b862f65720d2e5329319a8832f4e08d6a'
	//	and c.finish_time <=CURRENT_TIMESTAMP
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_history_reward").
		Joins("left join registers on registers.user_id = user_history_reward.user_id").
		Where("registers.invitation_code = ? and user_history_reward.finish_time<=?", invaliteCodeUser,
			time2.Format("2006-01-02 15:03:04")).Count(&totalCount).Error
	return
}

func GetCurrentUserNonce(userID string, dbData *db.UserHistoryTotal) (*big.Int, error, *ethclient.Client) {
	ethCli := contract.GetRewardRpcClient()
	if ethCli == nil {
		return nil, errors.New("无法链接rpc:" + config.Config.RewardChainRpc), nil
	}
	if dbData.CurrentNonce == "" {
		ptr, _ := contract.NewBBTTradeReward(common.HexToAddress(config.Config.RewardChainContractAddress), ethCli)
		nonce, err := ptr.Nonces(&bind.CallOpts{}, common.HexToAddress(userID))
		if err == nil {
			err = db.DB.MysqlDB.DefaultGormDB().Table("user_history_total").Where("user_id=?", dbData.UserID).Updates(map[string]interface{}{
				"current_nonce": nonce.String(),
			}).Error
		}

		return nonce, err, ethCli
	} else {
		decimalNonce, _ := decimal.NewFromString(dbData.CurrentNonce)
		return decimalNonce.BigInt(), nil, ethCli
	}

}
