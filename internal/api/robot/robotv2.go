package robot

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/utils"
	"Open_IM/pkg/walletRdbService"
	client2 "Open_IM/pkg/walletRdbService/swaprobotservice"
	"Open_IM/pkg/web3util"
	"context"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
	"time"
)

// CreateRobotV2
// @Summary		创建机器人V2
// @Description	创建机器人V2
// @Tags			机器人
// @ID				CreateRobotV2
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.CreateRobotV2Req	true	"参数"
// @Produce		json
// @Success		0	{object}    api.CreateRobotResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/robot/v2/create_robot [post]
func CreateRobotV2(c *gin.Context) {
	var (
		req api.CreateRobotV2Req
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
	var walletPtr *web3util.Wallet
	mnemonicstring := ""
	if !config.Config.WalletService.OpenFlag {
		byteEntropy, _ := web3util.NewEntropy(128)
		mnemonicstring, _ = web3util.NewMnemonicFromEntropy(byteEntropy)
		walletPtr, _ = web3util.NewFromMnemonic(mnemonicstring)
	} else {
		client, err := walletRdbService.GetRdbService()

		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "create robot failed", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
			return
		}
		reqHeader := metadata.New(map[string]string{"userid": userId})
		ctx := metadata.NewOutgoingContext(context.Background(), reqHeader)
		walletPtrResp, err := client.CreateUserRobot(ctx, &client2.CreateWalletMnemonicReq{
			UserID:   userId,
			CreateAt: time.Now().Format(time.DateTime),
			Status:   0,
			Ex:       "",
			FeeRate:  int32(feeRate),
			Sign:     req.Sign,
		})
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "create robot failed", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
			return
		}
		if walletPtrResp.BaseResp.StatusCode != 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "creator rebot error restat"})
			return
		}
		walletPtr, _ = web3util.NewFromMnemonic(walletPtrResp.Mnemonic)

	}
	ptAccount := &accounts.Account{
		URL: accounts.URL{
			Scheme: "",
			Path:   "m/44'/60'/0'/0/0",
		},
	}
	ethpublicAddress, _ := walletPtr.PublicKey(*ptAccount)
	ptAccount = &accounts.Account{
		URL: accounts.URL{
			Scheme: "",
			Path:   "m/84'/0'/0'/0/0",
		},
	}
	btcpublicAddress, _ := walletPtr.PublicKeyHexBtc(*ptAccount)

	//ptAccount = &accounts.Account{
	//	URL: accounts.URL{
	//		Scheme: "",
	//		Path:   "m/44'/714'/0'/0/0",
	//	},
	//}
	//bnbPublicAddressStr, _ := walletPtr.PublicKeyHexBnB(*ptAccount)
	ptAccount = &accounts.Account{
		URL: accounts.URL{
			Scheme: "",
			Path:   "m/44'/195'/0'/0/0",
		},
	}
	tronPublicKey, _ := walletPtr.PublicKeyHexTron(*ptAccount)
	dbRobot := &db.Robot{
		UserID:       userId,
		EthAddress:   crypto.PubkeyToAddress(*ethpublicAddress).String(),
		BnbAddress:   crypto.PubkeyToAddress(*ethpublicAddress).String(),
		BtcAddress:   btcpublicAddress,
		Status:       0,
		TronAddress:  tronPublicKey,
		CreateAt:     time.Now(),
		FeeRate:      feeRate,
		SnipeFeeRate: snipeFeeRate,
		Mnemonic:     mnemonicstring,
	}
	err := imdb.CreateRobot(dbRobot)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "create robot failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "",
		"data": map[string]string{
			"eth":  crypto.PubkeyToAddress(*ethpublicAddress).String(),
			"btc":  btcpublicAddress,
			"bsc":  crypto.PubkeyToAddress(*ethpublicAddress).String(),
			"tron": dbRobot.TronAddress,
		}})
}

// CheckExpireKey
// @Summary		检查密钥是否过期V2
// @Description	检查密钥是否过期V2
// @Tags			机器人
// @ID				CheckExpireKey
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetRobotReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.GetRobotResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/robot/v2/check_expire_key [post]
func CheckExpireKey(c *gin.Context) {
	var (
		req api.GetRobotReq
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
	mutexname := "CreateRobot:" + userId
	rs := db.DB.Pool
	mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
	ctx := context.Background()
	if err := mutex.LockContext(ctx); err != nil {
		return
	}
	defer mutex.UnlockContext(ctx)
	if _, err := imdb.GetUserRobot(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "un create robot"})
		return
	}
	client, err := walletRdbService.GetRdbService()

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "create robot failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	reqHeader := metadata.New(map[string]string{"userid": userId})
	ctx2 := metadata.NewOutgoingContext(context.Background(), reqHeader)
	resutlData, err := client.CheckIsExpireTime(ctx2, &client2.CheckIsNeedSign{
		UserID: userId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "un create robot"})
		return
	}
	if resutlData != nil && resutlData.IsNeedSign == false {
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": map[string]bool{"needSign": false}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 800, "errMsg": "", "data": map[string]bool{"needSign": true}})
	return
}

// ReLoadKey
// @Summary		重新签名密钥
// @Description	重新签名密钥
// @Tags			机器人
// @ID				ReLoadKey
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.CreateRobotV2Req	true	"参数"
// @Produce		json
// @Success		0	{object}    api.GetRobotResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/robot/v2/reloadKey [post]
func ReLoadKey(c *gin.Context) {
	var (
		req api.CreateRobotV2Req
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

	client, err := walletRdbService.GetRdbService()

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "create robot failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	reqHeader := metadata.New(map[string]string{"userid": userId})
	ctx2 := metadata.NewOutgoingContext(context.Background(), reqHeader)
	resutldata, err := client.ReloadMnemonic(ctx2, &client2.RequestUserID{
		UserID: userId,
		Sign:   req.Sign,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "un create robot"})
		return
	}
	if resutldata != nil && resutldata.BaseResp.StatusCode == 0 {
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "重新签名成功"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 800, "errMsg": "重新签名失败"})
	return
}

// ExportsWalletV2
// @Summary		导出钱包V2
// @Description	导出钱包V2
// @Tags			机器人
// @ID				ExportsWalletV2
// @Accept			json
// @Param			req		body	api.DelegateCallReq	true	"method=exportWallet,params里面sign不能为空,robotAddress为签名内容（要导出的钱包地址"
// @Produce		json
// @Success		0	{object}    api.DelegateCallResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/robot/v2/export_wallet [post]
func ExportsWalletV2(c *gin.Context) {
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
		//var valueExportRobotAddressStr string
		//if valueExportRobotAddress, ok := req.Params["robotAddress"]; !ok {
		//	c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": "export RobotAddress empty"})
		//	return
		//} else {
		//	valueExportRobotAddressStr = valueExportRobotAddress.(string)
		//	if valueExportRobotAddress == "" {
		//		c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": "export RobotAddress not string"})
		//		return
		//	}
		//}

		userId = strings.ToLower(userId)
		var valueExportRobotAddressStr = userId
		if !utils.VerifySignature(userId, valueSign.(string), valueExportRobotAddressStr) {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": "sign error "})
			return
		}
		//if _, err := imdb.GetRobotAddressIsUserRobot(userId, valueExportRobotAddressStr); err != nil {
		//	c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": "You not own this robot error "})
		//	return
		//}
		client, err := walletRdbService.GetRdbService()

		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "create robot failed", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
			return
		}
		reqHeader := metadata.New(map[string]string{"userid": userId})
		ctx2 := metadata.NewOutgoingContext(context.Background(), reqHeader)
		resutldata, err := client.GetMnemonic(ctx2, &client2.RequestUserID{
			UserID: userId,
			Sign:   valueSign.(string),
		})

		if err == nil && resutldata.BaseResp.StatusCode == 0 {
			resultdata, err := utils.WalletEncrypt(valueSign.(string), resutldata.Mnemonic)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resultdata})
				return
			}
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "method error"})
	return

}

// ImportWalletV2
// @Summary		导入钱包V2
// @Description	导入钱包V2
// @Tags			机器人
// @ID				ImportWalletV2
// @Accept			json
// @Param			req		body	api.ImportWalletMnemonic	true	"token 不能为空"
// @Param			token	header	string				true	"im token"
// @Produce		json
// @Success		0	{object}    api.DelegateCallResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/robot/v2/import_wallet [post]
func ImportWalletV2(c *gin.Context) {
	var (
		req api.ImportWalletMnemonic
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

	userId = strings.ToLower(userId)
	if !utils.VerifySignature(userId, req.Sign, userId) {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": "sign error "})
		return
	}
	client, err := walletRdbService.GetRdbService()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "create robot failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	reqHeader := metadata.New(map[string]string{"userid": userId})
	ctx2 := metadata.NewOutgoingContext(context.Background(), reqHeader)
	resutldata, err := client.ImportWallet(ctx2, &client2.ImportWalletMnemonicReq{
		UserID:   userId,
		CreateAt: time.Now().Format(time.DateTime),
		Ex:       "",
		Sign:     req.Sign,
		S:        req.S,
		I:        req.I,
		C:        req.C,
	})
	if err == nil && resutldata.BaseResp.StatusCode == 0 {
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
		walletPtr, _ := web3util.NewFromMnemonic(resutldata.Mnemonic)
		ptAccount := &accounts.Account{
			URL: accounts.URL{
				Scheme: "",
				Path:   "m/44'/60'/0'/0/0",
			},
		}
		ethpublicAddress, _ := walletPtr.PublicKey(*ptAccount)
		ptAccount = &accounts.Account{
			URL: accounts.URL{
				Scheme: "",
				Path:   "m/84'/0'/0'/0/0",
			},
		}
		btcpublicAddress, _ := walletPtr.PublicKeyHexBtc(*ptAccount)
		ptAccount = &accounts.Account{
			URL: accounts.URL{
				Scheme: "",
				Path:   "m/44'/195'/0'/0/0",
			},
		}
		tronPublicKey, _ := walletPtr.PublicKeyHexTron(*ptAccount)
		dbRobot := &db.Robot{
			UserID:       userId,
			EthAddress:   crypto.PubkeyToAddress(*ethpublicAddress).String(),
			BnbAddress:   crypto.PubkeyToAddress(*ethpublicAddress).String(),
			BtcAddress:   btcpublicAddress,
			TronAddress:  tronPublicKey,
			Status:       0,
			CreateAt:     time.Now(),
			FeeRate:      feeRate,
			SnipeFeeRate: snipeFeeRate,
			Mnemonic:     "",
		}
		err = imdb.CreateRobot(dbRobot)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "",
				"data": map[string]string{
					"eth":  dbRobot.EthAddress,
					"btc":  dbRobot.BtcAddress,
					"bsc":  dbRobot.BnbAddress,
					"tron": dbRobot.TronAddress,
				}})
			return
		}
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "s,i,c,error" + err.Error()})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "s,i,c,error" + resutldata.BaseResp.StatusMessage})
	}

	return

}
