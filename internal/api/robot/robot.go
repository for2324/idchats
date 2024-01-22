package robot

import (
	"Open_IM/internal/api/brc20/services/unisat_wallet"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbSwapRobot "Open_IM/pkg/proto/swaprobot"
	"Open_IM/pkg/utils"
	"Open_IM/pkg/web3util"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"github.com/sourcegraph/conc"
	"gorm.io/gorm"
	"k8s.io/utils/strings/slices"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CreateRobot
// @Summary		创建机器人
// @Description	创建机器人
// @Tags			机器人
// @ID				CreateRobot
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.CreateRobotReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.CreateRobotResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/robot/create_robot [post]
func CreateRobot(c *gin.Context) {
	var (
		req api.CreateRobotReq
		//resp api.CreateRobotResp
	)

	//fmt.Println(utils.StructToJsonString(c.Request.Header))
	//fmt.Println(c.Request.Header.Get("Referer"))

	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	var ok bool
	var errInfo string
	ok, userId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	mutexname := "CreateRobot:" + userId
	rs := db.DB.Pool
	mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
	ctx := context.Background()
	if err := mutex.LockContext(ctx); err != nil {
		return
	}
	defer mutex.UnlockContext(ctx)
	if robotInfo, err := imdb.GetUserRobot(userId); err == nil {
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "",
			"data": map[string]string{
				"eth": robotInfo.EthAddress,
				"btc": robotInfo.BtcAddress,
				"bsc": robotInfo.BnbAddress,
			}})
		return
	}

	//暂时排除并发创建私钥的情况：
	byteEntropy, _ := web3util.NewEntropy(128)
	mnemonicstring, _ := web3util.NewMnemonicFromEntropy(byteEntropy)
	walletPtr, _ := web3util.NewFromMnemonic(mnemonicstring)
	ptAccount := &accounts.Account{
		URL: accounts.URL{
			Scheme: "",
			Path:   "m/44'/60'/0'/0/0",
		},
	}
	ethprivateKeyHex, _ := walletPtr.PrivateKeyHex(*ptAccount)
	ethpublicAddress, _ := walletPtr.PublicKey(*ptAccount)
	privateKeyCrypt, err := utils.AesEncrypt([]byte(ethprivateKeyHex), []byte("U2FsdGVkX1+4xoFd+2jiqf+m16e3EdEQ"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	ptAccount = &accounts.Account{
		URL: accounts.URL{
			Scheme: "",
			Path:   "m/84'/0'/0'/0/0",
		},
	}
	btcprivateKeyHex, _ := walletPtr.PrivateKeyHexBtc(*ptAccount)
	btcprivateKeyHexStr, err := utils.AesEncrypt([]byte(btcprivateKeyHex), []byte("U2FsdGVkX1+4xoFd+2jiqf+m16e3EdEQ"))
	btcpublicAddress, _ := walletPtr.PublicKeyHexBtc(*ptAccount)

	ptAccount = &accounts.Account{
		URL: accounts.URL{
			Scheme: "",
			Path:   "m/44'/714'/0'/0/0",
		},
	}

	feeRate := 2
	snipeFeeRate := 10
	if str := checkIsProRobot(c); str != "" {
		feeRateMap := config.Config.UniswapRobot.FeeRateMap[str]
		for _, value := range feeRateMap {
			if value.Method == "swap" {
				feeRate = value.FeeRate
			}
			if value.Method == "snipe" {
				snipeFeeRate = value.FeeRate
			}
		}
	}
	dbRobot := &db.Robot{
		UserID:        userId,
		EthAddress:    crypto.PubkeyToAddress(*ethpublicAddress).String(),
		EthPrivateKey: base64.StdEncoding.EncodeToString(privateKeyCrypt),
		BnbAddress:    "",
		BnbPrivateKey: "",
		Status:        0,
		Mnemonic:      mnemonicstring,
		BtcAddress:    btcpublicAddress,
		BtcPrivateKey: base64.StdEncoding.EncodeToString(btcprivateKeyHexStr),
		CreateAt:      time.Now(),
		FeeRate:       feeRate,
		SnipeFeeRate:  snipeFeeRate,
	}
	err = imdb.CreateRobot(dbRobot)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "create robot failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "",
		"data": map[string]string{
			"eth": crypto.PubkeyToAddress(*ethpublicAddress).String(),
			"btc": btcpublicAddress,
			"bsc": crypto.PubkeyToAddress(*ethpublicAddress).String(),
		}})
}

func checkIsProRobot(c *gin.Context) string {
	stringRefer := c.Request.Header.Get("Referer")
	at, err := url.Parse(stringRefer)
	if err == nil && strings.HasPrefix(at.Host, "pro.") {
		return "pro"
	}
	return "bibot"
}

// GetRobot
// @Summary		获取用户机器人信息
// @Description	获取用户机器人信息
// @Tags			机器人
// @ID				GetRobot
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetRobotReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.GetRobotResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/robot/get_robot [post]
func GetRobot(c *gin.Context) {
	var (
		req  api.GetRobotReq
		resp api.GetRobotResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	var ok bool
	var errInfo string
	ok, userId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	//并发打db 会导致db死
	robot, err := imdb.GetUserRobot(userId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 0, "errMsg": err.Error()})
		return
	} else if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "无钱包地址"})
		return
	}
	robot.EthPrivateKey = ""
	robot.BtcPrivateKey = ""
	robot.Mnemonic = ""
	resp.Data = make(map[string]string, 0)
	resp.Data["eth"] = robot.EthAddress
	resp.Data["btc"] = robot.BtcAddress
	resp.Data["bsc"] = robot.BnbAddress
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp.Data})
}

// GetRobotAllTask
// @Summary		获取机器人当前的任务
// @Description	获取机器人当前的任务
// @Tags			机器人
// @ID				GetRobotAllTask
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.CreateRobotReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.DelegateCallResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/robot/get_user_all_robot_task [post]
func GetRobotAllTask(c *gin.Context) {

}

// DelegateCall
// @Summary		机器人委托调用
// @Description	机器人委托调用
// @Tags			机器人
// @ID				DelegateCall
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.DelegateCallReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.DelegateCallResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/robot/delegate_call [post]
func DelegateCall(c *gin.Context) {
	var (
		req api.DelegateCallReq
		// resp api.DelegateCallResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	// get user id from token
	var ok bool
	var errInfo string
	ok, userId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok && !slices.Contains([]string{"exploredPair", "tokenPrice", "quote", "newTokens", "tokensPrice"}, req.Method) {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	bibotKey := c.Request.Header.Get("bibotKey")
	// 创建监听任务
	etcdConn := getcdv3.GetDefaultConn(
		config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.SwapRobotPort,
		req.OperationID,
	)
	if strings.EqualFold(req.Method, "withdraw") {
		var valueSign interface{}
		var ok bool
		var toAddress interface{}

		if valueSign, ok = req.Params["sign"]; !ok {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": "sign empty"})
			return
		}
		if toAddress, ok = req.Params["recipientAddress"]; !ok {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": "recipientAddress empty"})
			return
		}
		toAddressStr := strings.ToLower(toAddress.(string))
		userId = strings.ToLower(userId)
		if !utils.VerifySignature(userId, valueSign.(string), toAddressStr) {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": "sign error "})
			return
		}
	}

	client := pbSwapRobot.NewSwaprobotClient(etcdConn)
	paramsStr, err := json.Marshal(req.Params)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "json.Marshal failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resultData, err := client.BotOperation(c.Request.Context(), &pbSwapRobot.BotOperationReq{
		OperatorID: req.OperationID,
		Method:     req.Method,
		BiBotKey:   bibotKey,
		UserID:     userId,
		Params:     paramsStr,
	})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "client.BotOperation failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if resultData.CommonResp.ErrCode != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": resultData.CommonResp.ErrCode,
			"errMsg": resultData.CommonResp.ErrMsg})
		return
	}
	resp := api.DelegateCallResp{}
	json.Unmarshal(resultData.Data, &resp.Data)
	c.JSON(http.StatusOK, resp)
}

type ParamReBackData struct {
	OrdId string `json:"ordId"`
}

type Param struct {
	FromSymbol  string  `json:"fromSymbol"`
	Amount      float64 `json:"amount"`
	ToSymbol    string  `json:"toSymbol"`
	Tp          float64 `json:"tp"`
	Sl          float64 `json:"sl"`
	OrdStatus   string  `json:"ordStatus"`
	MinimumOut  float64 `json:"minimumOut"`
	DeadlineDay int64   `json:"deadlineDay"`
	SearchBy    string  `json:"searchBy"`
}
type SwapBotInfoParam struct {
	OrdId     string `json:"ordId"`
	Method    string `json:"method"`
	OrdStatus string `json:"ordStatus"`
	Params    Param  `json:"params"`
}
type SwapBotInfoParamReq struct {
	OperationID string `json:"operationID"`
	TxHash      string `json:"txHash"`
	PrivateKey  string `json:"privateKey"`
	SwapBotInfoParam
}

type RobotRunTaskPostReq struct {
	PrivateKey string                 `json:"privateKey"`
	Params     map[string]interface{} `json:"params"`
	Address    string                 `json:"address"`
	TimeStamp  int64                  `json:"timeStamp"`
	Method     string                 `json:"method"`
}
type RobotRunTaskResp struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

// CheckOrderStatus
// @Summary		查询到交易状态
// @Description	查询到交易状态
// @Tags			机器人
// @ID				CheckOrderStatus
// @Accept			json
// @Param			req		body	SwapBotInfoParamReq	true	"method=getOrd,ordid"
// @Produce		json
// @Success		0	{object}    api.DelegateCallResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/robot/check_ord_status [post]
func CheckOrderStatus(c *gin.Context) {
	var (
		tempResultData SwapBotInfoParamReq
	)
	if err := c.BindJSON(&tempResultData); err != nil {
		log.NewError(tempResultData.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	if tempResultData.Method == "" ||
		tempResultData.OrdId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrInternal.ErrCode,
			"errMsg": "参数错误"})
		return
	}

	if tempResultData.Method == "getOrd" {

	}
	c.JSON(http.StatusOK, gin.H{"errCode": 801, "errMsg": "失败"})
	return
}

// ChangeOrderStatus
// @Summary		交易状态通知
// @Description	交易状态通知
// @Tags			机器人
// @ID				ChangeOrderStatus
// @Accept			json
// @Param			req		body	SwapBotInfoParamReq	true	"method=statusTask,ordid和status 不能为空"
// @Produce		json
// @Success		0	{object}    api.DelegateCallResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/robot/change_ord_status [post]
func ChangeOrderStatus(c *gin.Context) {
	var (
		tempResultData SwapBotInfoParamReq
	)
	if err := c.BindJSON(&tempResultData); err != nil {
		log.NewError(tempResultData.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	// if tempResultData.Method == "" ||
	// 	tempResultData.OrdId == "" ||
	// 	tempResultData.OrdStatus == "" {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrInternal.ErrCode,
	// 		"errMsg": "参数错误"})
	// 	return
	// }
	// errormsg := "失败"
	// if tempResultData.Method == "statusTask" {
	// 	if err := db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
	// 		err := tx.Table("swap_robot_task").
	// 			Where("ord_id =?", tempResultData.OrdId).Updates(map[string]interface{}{
	// 			"order_status": tempResultData.OrdStatus,
	// 			"tx_hash":      tempResultData.TxHash,
	// 		}).Error
	// 		if err == nil {
	// 			err = tx.Table("swap_robot_task_log").Create(&db.RoBotTaskLog{
	// 				CreatedAt:   time.Now(),
	// 				UpdatedAt:   time.Now(),
	// 				OrdID:       tempResultData.OrdId,
	// 				OrderStatus: tempResultData.OrdStatus,
	// 				Method:      "statusTask",
	// 				TxHash:      tempResultData.TxHash,
	// 			}).Error
	// 		}
	// 		return err

	// 	}); err == nil {
	// 		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "通知成功"})
	// 	}

	// }
	// c.JSON(http.StatusOK, gin.H{"errCode": 801, "errMsg": errormsg})
}

// TokenPrice
// @Summary		查询token价格汇率
// @Description	查询token价格汇率
// @Tags			机器人
// @ID				TokenPrice
// @Accept			json
// @Param			req		body	api.TokenPriceReq	true	"method=statusTask,ordid和status 不能为空"
// @Produce		json
// @Success		0	{object}    api.TokenPriceResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/robot/token_price [post]
func TokenPrice(c *gin.Context) {
	var (
		req api.TokenPriceReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	queryParams := url.Values{}
	for _, v := range req.Tokens {
		queryParams.Add("tokens", v)
	}
	uri := config.Config.UniswapRobot.ThorSwapEndpoint + "/api/thor/token_price?" + queryParams.Encode()
	data, err := utils.HttpGet(uri)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "HttpGet failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrInternal.ErrCode,
			"errMsg": constant.ErrInternal.ErrMsg})
		return
	}
	var resp api.TokenPriceResp
	json.Unmarshal(data, &resp.Data)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "data": resp.Data})
}

// ExportsWallet
// @Summary		导出钱包废弃
// @Description	导出钱包废弃
// @Tags			机器人
// @ID				ExportsWallet
// @Accept			json
// @Param			req		body	api.DelegateCallReq	true	"method=exportWallet,params里面sign 不能为空"
// @Produce		json
// @Success		0	{object}    api.DelegateCallResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/robot/export_wallet [post]
func ExportsWallet(c *gin.Context) {
	var (
		req api.DelegateCallReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, userId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	if req.Method == "exportWallet" {
		var valueSign interface{}
		if valueSign, ok = req.Params["sign"]; !ok {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": "sign empty"})
			return
		}
		userId = strings.ToLower(userId)
		if !utils.VerifySignature(userId, valueSign.(string), userId) {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": "sign error "})
			return
		}
		robot, err := imdb.GetUserRobot(userId)
		if err == nil {
			resultdata, err := utils.WalletEncrypt(valueSign.(string), robot.Mnemonic)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resultdata})
				return
			}
		}

	}
	c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "method error"})
	return

}

// TradingVolume
// @Summary		查询交易量
// @Description	查询交易量
// @Tags			机器人
// @ID				TradingVolume
// @Accept			json
// @Param			req		body	TFeeRobotRequest	true	"不能为空"
// @Produce		json
// @Success		0	{object}    TFeeRobotRewardResponse
// sponse
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/robot/v2/total_volume [post]
func TradingVolume(c *gin.Context) {
	var (
		req TFeeRobotRequest
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req.MerchantId = checkIsProRobot(c)
	if req.EndTimestamp == 0 {
		req.EndTimestamp = time.Now().Unix()
	}
	fmt.Println(config.Config.UniswapRobot.BibotUri + "/api/swap/record/statistics")
	resultUrl := fmt.Sprintf("merchantId=%s&merchantUid=%s&startTimestamp=%d&endTimestamp=%d",
		req.MerchantId, req.MerchantUid, req.StartTimestamp, req.EndTimestamp)
	if req.MerchantUid == "" {
		resultUrl = fmt.Sprintf("merchantId=%s&merchantUid=%s&startTimestamp=%d&endTimestamp=%d&isFake=true",
			req.MerchantId, req.MerchantUid, req.StartTimestamp, req.EndTimestamp)
	}
	resp, err := resty.New().R().SetQueryString(resultUrl).
		Get(config.Config.UniswapRobot.BibotUri + "/api/swap/record/statistics")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "param error"})
		return
	}
	var dataRobotFeeResp TFeeRobotRewardResponse
	json.Unmarshal(resp.Body(), &dataRobotFeeResp)
	if dataRobotFeeResp.ErrCode != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": dataRobotFeeResp.ErrMsg})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": dataRobotFeeResp.Data})
		return
	}
}

// 请求七天交易量的数据
// TradingVolume7Days
// @Summary		查询7天交易量
// @Description	查询7天交易量
// @Tags			机器人
// @ID				TradingVolume7Days
// @Accept			json
// @Param			req		body	TTotalVolumeData	true	"operation有用,类型trade,pledge,lpPledge"
// @Produce		json
// @Success		0	{object}    TFeeRobotTotalVolumeResponse
// sponse
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/robot/v2/total_7days_volume [post]
func TradingVolume7Days(c *gin.Context) {
	var (
		req TTotalVolumeData
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var dataRobotFeeResp TFeeRobotTotalVolumeResponse
	dataRobotFeeResp.Data = make([]*TFeeRobotDayTotalVolumeData, 0, 8)
	switch req.VolumeType {
	case "pledge":
		GetPledgeVolumn(&dataRobotFeeResp)
	case "lpPledge":
		GetLpPledgeVolumn(&dataRobotFeeResp)
	case "trade":
		GetTradeVolumn(&dataRobotFeeResp)
	case "brc20_obbt":
		GetBrc20ObbtVolumn(&dataRobotFeeResp)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "参数错误"})
		return

	}
	c.JSON(http.StatusOK, dataRobotFeeResp)
}
func GetLpPledgeVolumn(dataRobotFeeResp *TFeeRobotTotalVolumeResponse) {
	now := time.Now()
	eightDaysAgo := now.AddDate(0, 0, -7)
	formattedDate := eightDaysAgo.Format("2006-01-02")
	var dbData []*db.BLPPledgeLog
	err := db.DB.MysqlDB.DefaultGormDB().Raw(`SELECT  a.id,a.chain,a.created_at,a.contract,a.total_lock,a.block_height,a.block_date FROM ( select max(id) as id,block_date from blp_pledge_log WHERE block_date>=? and block_date<= ? GROUP BY block_date ORDER BY block_date) c left join blp_pledge_log  a  on  c.id= a.id`, formattedDate, now.Format("2006-01-02")).Scan(&dbData).
		Error
	if err != nil {
		return
	}
	for i := 0; i <= 7; i++ {
		checkDay := i - 7
		checkDayTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
			AddDate(0, 0, checkDay)
		targetDay := &TFeeRobotDayTotalVolumeData{
			Timestamp:   checkDayTime.Unix(),
			TotalVolume: "0",
		}
		for _, value := range dbData {
			if value.BlockDate.Unix() == checkDayTime.Unix() {
				targetDay.TotalVolume = value.TotalLock
			}
		}
		dataRobotFeeResp.Data = append(dataRobotFeeResp.Data, targetDay)
	}
	return
}

func GetPledgeVolumn(dataRobotFeeResp *TFeeRobotTotalVolumeResponse) {
	//var concPool conc.WaitGroup
	now := time.Now()
	eightDaysAgo := now.AddDate(0, 0, -7)
	formattedDate := eightDaysAgo.Format("2006-01-02")
	var dbData []*db.BbtPledgeLog
	err := db.DB.MysqlDB.DefaultGormDB().Raw(`SELECT  a.id,a.chain,a.created_at,a.contract,a.total_lock,a.block_height,a.block_date FROM ( select max(id) as id,block_date from bbt_pledge_log WHERE block_date>=? and block_date<= ? GROUP BY block_date ORDER BY block_date) c left join bbt_pledge_log  a  on  c.id= a.id`, formattedDate, now.Format("2006-01-02")).Scan(&dbData).
		Error
	if err != nil {
		return
	}
	for i := 0; i <= 7; i++ {
		checkDay := i - 7
		checkDayTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
			AddDate(0, 0, checkDay)
		targetDay := &TFeeRobotDayTotalVolumeData{
			Timestamp:   checkDayTime.Unix(),
			TotalVolume: "0",
		}
		for _, value := range dbData {
			if value.BlockDate.Unix() == checkDayTime.Unix() {
				targetDay.TotalVolume = value.TotalLock
			}
		}
		dataRobotFeeResp.Data = append(dataRobotFeeResp.Data, targetDay)
	}
	return
	//concPool.Wait()
}
func GetTradeVolumn(dataRobotFeeResp *TFeeRobotTotalVolumeResponse) {
	var concPool conc.WaitGroup
	now := time.Now()
	for i := 0; i <= 7; i++ {
		checkDay := i - 7
		checkDayTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, checkDay)
		dataRobotFeeResp.Data = append(dataRobotFeeResp.Data, &TFeeRobotDayTotalVolumeData{
			Timestamp:   checkDayTime.Unix(),
			TotalVolume: "0",
		})
		nIndex := i
		concPool.Go(func() {
			if config.Config.IsPublicEnv && db.DB.RDB.Exists(context.Background(), "dayTotalVolume:"+checkDayTime.Format("20060102")).Val() == 1 {
				totalVolume, _ := db.DB.RDB.Get(context.Background(), "dayTotalVolume:"+checkDayTime.Format("20060102")).Result()
				dataRobotFeeResp.Data[nIndex].TotalVolume = totalVolume
			} else {
				beginTime, endTime := getTargetDay(checkDay)
				totalVolume, _ := readFromRobotServerToGetTradingVolume(beginTime, endTime)
				fmt.Println(totalVolume)
				dataRobotFeeResp.Data[nIndex].TotalVolume = totalVolume

				if checkDay != 0 {
					db.DB.RDB.Set(context.Background(), "dayTotalVolume:"+checkDayTime.Format("20060102"), totalVolume, 0)
				}
			}
		})
	}
	concPool.Wait()
}

func GetBrc20ObbtVolumn(dataRobotFeeResp *TFeeRobotTotalVolumeResponse) {
	now := time.Now()
	eightDaysAgo := now.AddDate(0, 0, -7)
	formattedDate := eightDaysAgo.Format("2006-01-02")
	var dbData []*db.BbtPledgeLog
	err := db.DB.MysqlDB.DefaultGormDB().Raw(`SELECT  a.id,a.chain,a.created_at,a.contract,a.total_lock,a.block_height,a.block_date 
FROM ( select max(id) as id,block_date from obbt_pledge_log_day_report 
                                       WHERE block_date>=? and block_date<= ? GROUP BY block_date ORDER BY block_date) 
    c left join obbt_pledge_log_day_report  a  on  c.id= a.id`, formattedDate, now.Format("2006-01-02")).Scan(&dbData).
		Error
	if err != nil {
		return
	}
	for i := 0; i <= 7; i++ {
		checkDay := i - 7
		checkDayTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
			AddDate(0, 0, checkDay)
		targetDay := &TFeeRobotDayTotalVolumeData{
			Timestamp:   checkDayTime.Unix(),
			TotalVolume: "0",
		}
		for _, value := range dbData {
			if value.BlockDate.Unix() == checkDayTime.Unix() {
				targetDay.TotalVolume = value.TotalLock
			}
		}
		dataRobotFeeResp.Data = append(dataRobotFeeResp.Data, targetDay)
	}
	tempUnsatPtr := new(unisat_wallet.UnisatWeb)
	totalZhiYa, _ := tempUnsatPtr.GetTotalStake()
	dataRobotFeeResp.Data[len(dataRobotFeeResp.Data)-1].TotalVolume = totalZhiYa
	return
}

func getTargetDay(day int) (daytimeZeroHourTimeStap, daytimeLastHourTimeStap int64) {
	now := time.Now()
	daytimeZeroHourTimeStap = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, day).Unix()
	daytimeLastHourTimeStap = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, day+1).Unix() - 1
	return
}
func readFromRobotServerToGetTradingVolume(beginTime, endTime int64) (totalVolume string, err error) {
	//resultUrl := fmt.Sprintf("merchantId=&merchantUid=&startTimestamp=%d&endTimestamp=%d", beginTime, endTime)
	resultUrl := fmt.Sprintf("merchantId=&merchantUid=&startTimestamp=%d&endTimestamp=%d&isFake=true", beginTime, endTime)

	fmt.Println(resultUrl)
	resp, err := resty.New().R().SetQueryString(resultUrl).
		Get(config.Config.UniswapRobot.BibotUri + "/api/swap/record/statistics")
	if err != nil {
		return
	}
	if resp.StatusCode() == http.StatusOK {
		var result TFeeRobotResponse
		if err = json.Unmarshal(resp.Body(), &result); err == nil && result.ErrCode == 0 && len(result.Data) > 0 {
			returnDecimal, _ := decimal.NewFromString("0")
			for key := range result.Data {
				returnDecimal = returnDecimal.Add(decimal.NewFromFloat(result.Data[key].SellUsdPrice))
			}
			return returnDecimal.String(), nil
		}
	}
	return
}

type TFeeRobotDayTotalVolumeData struct {
	Timestamp   int64  `json:"timestamp"`
	TotalVolume string `json:"totalVolume"`
}
type TFeeRobotTotalVolumeResponse struct {
	Data    []*TFeeRobotDayTotalVolumeData `json:"data"`
	ErrCode int                            `json:"errCode"`
	ErrMsg  string                         `json:"errMsg"`
}
type TFeeRobotRewardResponse struct {
	Data []*struct {
		SellAmount     float64 `json:"sellAmount"`
		SellUsdPrice   float64 `json:"sellUsdPrice"`
		FeeUsdPrice    float64 `json:"feeUsdPrice"`
		ActualSlippage float64 `json:"actualSlippage"`
		Type           string  `json:"type"`
	} `json:"data"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}
type TFeeRobotRequest struct {
	StartTimestamp int64  `json:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp"`
	MerchantId     string `json:"merchantId"`
	MerchantUid    string `json:"merchantUid"`
	OperationID    string `json:"operationID"`
}
type TTotalVolumeData struct {
	OperationID string `json:"operationID"`
	VolumeType  string `json:"volumeType"` //trade pledge ,
}
type TFeeRobotResponse struct {
	Data []*struct {
		SellAmount     float64 `json:"sellAmount"`
		SellUsdPrice   float64 `json:"sellUsdPrice"`
		FeeUsdPrice    float64 `json:"feeUsdPrice"`
		ActualSlippage float64 `json:"actualSlippage"`
		Type           string  `json:"type"`
	} `json:"data"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}
