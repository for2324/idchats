package register

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	http2 "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/auth"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ParamsLogin struct {
	UserID      string `json:"userID"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Platform    int32  `json:"platform"`
	OperationID string `json:"operationID" binding:"required"`
	AreaCode    string `json:"areaCode"`
}
type LoginMetaMaskRequest struct {
	PublicAddress  string `json:"publicAddress" validate:"required"`
	InvitationCode string `json:"invitationCode"`
	Signature      string `json:"signature" validate:"required"`
}

func PostNonceVerify(c *gin.Context) {
	var params LoginMetaMaskRequest
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	accountKeyNonce, err := db.DB.GetAccountNonce(params.PublicAddress)
	if err != nil || len(params.Signature) == 0 {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NotRegistered, "errMsg": "地址未授权"})
	}
	if sigValid := utils.VerifySignatureEip4361(params.PublicAddress, params.Signature, accountKeyNonce); !sigValid {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.SignErrorWeb3, "errMsg": "地址未授权"})
		return
	}
	chainid := c.Request.Header.Get("chainId")
	ip := c.Request.Header.Get("X-Forward-For")
	if ip == "" {
		ip = c.ClientIP()
	}
	OperationID := utils.OperationIDGenerator()
	ok, opUserID, _ := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), OperationID)
	if !ok || !utils.IsContain(opUserID, config.Config.Manager.AppManagerUid) {
		Limited, LimitError := im_mysql_model.IsLimitRegisterIp(ip)
		if LimitError != nil {
			log.Error(OperationID, utils.GetSelfFuncName(), LimitError, ip)
			c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.ErrDB.ErrCode, "errMsg": LimitError.Error()})
			return
		}
		if Limited {
			log.NewInfo(OperationID, utils.GetSelfFuncName(), "is limited", ip, "params:", params)
			c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.RegisterLimit, "errMsg": "limited"})
			return
		}
	}

	var account, userID string
	account = params.PublicAddress
	url := ""
	var bMsg []byte
	if _, err := im_mysql_model.GetRegisterInfo(account); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.SignErrorWeb3, "errMsg": err.Error()})
		return
	} else if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		if params.InvitationCode != "" {
			if resultInvitationUserInfo, _ := im_mysql_model.GetRegisterInfo(params.InvitationCode); resultInvitationUserInfo.Account == "" {
				c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.InvitationError, "errMsg": "邀请人未注册"})
				return
			}
		}
		_, err2 := im_mysql_model.GetUserByUserID(account)
		if errors.Is(err2, gorm.ErrRecordNotFound) {
			//未注册的用户先进行注册
			openIMRegisterReq := api.UserRegisterReq{}
			userID = params.PublicAddress
			url = config.Config.Demo.ImAPIURL + "/auth/user_register"
			openIMRegisterReq.OperationID = OperationID
			openIMRegisterReq.Platform = 5
			openIMRegisterReq.UserID = userID
			openIMRegisterReq.Nickname = account
			openIMRegisterReq.Secret = config.Config.Secret
			int32ChainId, _ := strconv.Atoi(chainid)
			openIMRegisterReq.ChainId = int32(int32ChainId)
			fmt.Println("当前的ChainId 为：", int32ChainId, chainid)
			openIMRegisterResp := api.UserRegisterResp{}
			log.NewDebug(OperationID, utils.GetSelfFuncName(), "register req:", utils.Interface2JsonString(openIMRegisterReq))
			bMsg, err = http2.Post(url, openIMRegisterReq, 2)
			if err != nil {
				log.NewError(OperationID, "request openIM register error", account, "err", err.Error())
				c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": err.Error()})
				return
			}
			err = json.Unmarshal(bMsg, &openIMRegisterResp)
			if err != nil || openIMRegisterResp.ErrCode != 0 {
				log.NewError(OperationID, "request openIM register error", account, "err", "resp: ", openIMRegisterResp.ErrCode)
				if err != nil {
					log.NewError(OperationID, utils.GetSelfFuncName(), err.Error())
					c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": "register limit"})
					return
				}
				if openIMRegisterResp.ErrCode != 0 {
					c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": "register failed: " + openIMRegisterResp.ErrMsg})
					return
				}
			}
			userID = account
			log.Info(OperationID, "begin store mysql", account)
			err = im_mysql_model.SetPasswordWithInvitationCode(account, "xxxxxxxxx111", "", userID, chainid, ip, params.InvitationCode)
			if err != nil {
				log.NewError(OperationID, "set phone number password error", account, "err", err.Error())
				c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": err.Error()})
				return
			}
			if err := im_mysql_model.InsertIpRecord(userID, ip); err != nil {
				log.NewError(OperationID, utils.GetSelfFuncName(), userID, ip, err.Error())
			}
			log.Info(OperationID, "end  setuserInfo", account)
			// demo onboarding
			if config.Config.Demo.OnboardProcess {
				select {
				case Ch <- OnboardingProcessReq{
					OperationID: OperationID,
					UserID:      userID,
					NickName:    account,
					FaceURL:     "",
					PhoneNumber: "",
					Email:       "",
				}:
				case <-time.After(time.Second * 2):
					log.NewWarn(OperationID, utils.GetSelfFuncName(), "to ch timeOut")
				}
			}

			select {
			case ChOnNewLogin <- OnboardingProcessReq{
				OperationID: OperationID,
				UserID:      userID,
				NickName:    account,
				FaceURL:     "",
				PhoneNumber: "",
				Email:       "",
			}:
			case <-time.After(time.Second * 2):
				log.NewWarn(OperationID, utils.GetSelfFuncName(), "to ch timeOut")
			}
			// register add friend
			select {
			case ChImportFriend <- &pbFriend.ImportFriendReq{
				OperationID: OperationID,
				FromUserID:  userID,
				OpUserID:    config.Config.Manager.AppManagerUid[0],
			}:
			case <-time.After(time.Second * 2):
				log.NewWarn(OperationID, utils.GetSelfFuncName(), "to ChImportFriend timeOut")
			}
		} else {
			userID = account
		}
	} else {
		resultdata2, err2 := im_mysql_model.GetUserByUserID(account)

		if errors.Is(err2, gorm.ErrRecordNotFound) {
			//未注册的用户先进行注册
			openIMRegisterReq := api.UserRegisterReq{}
			userID = params.PublicAddress
			url = config.Config.Demo.ImAPIURL + "/auth/user_register"
			openIMRegisterReq.OperationID = OperationID
			openIMRegisterReq.Platform = 5
			openIMRegisterReq.UserID = userID
			openIMRegisterReq.Nickname = account
			openIMRegisterReq.Secret = config.Config.Secret
			int32ChainId, _ := strconv.Atoi(chainid)
			openIMRegisterReq.ChainId = int32(int32ChainId)
			fmt.Println("当前的ChainId 为：", int32ChainId, chainid)
			openIMRegisterResp := api.UserRegisterResp{}
			log.NewDebug(OperationID, utils.GetSelfFuncName(), "register req:", utils.Interface2JsonString(openIMRegisterReq))
			bMsg, err = http2.Post(url, openIMRegisterReq, 2)
			if err != nil {
				log.NewError(OperationID, "request openIM register error", account, "err", err.Error())
				c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": err.Error()})
				return
			}
			err = json.Unmarshal(bMsg, &openIMRegisterResp)
			if err != nil || openIMRegisterResp.ErrCode != 0 {
				log.NewError(OperationID, "request openIM register error", account, "err", "resp: ", openIMRegisterResp.ErrCode)
				if err != nil {
					log.NewError(OperationID, utils.GetSelfFuncName(), err.Error())
					c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": "register limit"})
					return
				}
				if openIMRegisterResp.ErrCode != 0 {
					c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": "register failed: " + openIMRegisterResp.ErrMsg})
					return
				}
			}
		}

		userID = resultdata2.UserID
		im_mysql_model.UpdateUserInfo(db.User{
			UserID:  userID,
			Chainid: utils.StringToInt32(chainid),
		})
	}

	fmt.Println("user id is :>>>>", userID)
	mutexname := "AddKeyApi:" + userID
	rs := db.DB.Pool
	mutex := rs.NewMutex(mutexname, redsync.WithTries(3), redsync.WithRetryDelay(time.Second*1), redsync.WithExpiry(time.Second*10))
	ctx := context.Background()
	if err := mutex.LockContext(ctx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "正在生成密钥"})
		return
	}
	defer mutex.UnlockContext(ctx)
	if isHave, _ := im_mysql_model.IsUserHaveApiKey(userID); !isHave {
		im_mysql_model.InsertUserApiKey(userID)
	}

	req := &rpc.UserTokenReq{Platform: 5, FromUserID: userID, OperationID: OperationID, LoginIp: ip}
	log.NewInfo(req.OperationID, "UserToken args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + " getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewAuthClient(etcdConn)
	reply, err := client.UserToken(context.Background(), req)
	if err != nil {
		errMsg := req.OperationID + " UserToken failed " + err.Error() + " req: " + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	resp := api.UserTokenResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg},
		UserToken: api.UserTokenInfo{UserID: req.FromUserID, Token: reply.Token, ExpiredTime: reply.ExpiredTime}}
	log.NewInfo(req.OperationID, "UserToken return ", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "验证通过，", "data": resp.UserToken})
}

func Login(c *gin.Context) {
	params := ParamsLogin{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	var account string
	if params.Email != "" {
		account = params.Email
	} else if params.PhoneNumber != "" {
		account = params.PhoneNumber
	} else {
		account = params.UserID
	}

	r, err := im_mysql_model.GetRegister(account, params.AreaCode, params.UserID)
	if err != nil {
		log.NewError(params.OperationID, "user have not register", params.Password, account, err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NotRegistered, "errMsg": "Mobile phone number is not registered"})
		return
	}
	if r.Password != params.Password {
		log.NewError(params.OperationID, "password  err", params.Password, account, r.Password, r.Account)
		c.JSON(http.StatusOK, gin.H{"errCode": constant.PasswordErr, "errMsg": "password err"})
		return
	}
	var userID string
	if r.UserID != "" {
		userID = r.UserID
	} else {
		userID = r.Account
	}
	ip := c.Request.Header.Get("X-Forward-For")
	if ip == "" {
		ip = c.ClientIP()
	}
	url := fmt.Sprintf("%s/auth/user_token", config.Config.Demo.ImAPIURL)
	openIMGetUserToken := api.UserTokenReq{}
	openIMGetUserToken.OperationID = params.OperationID
	openIMGetUserToken.Platform = params.Platform
	openIMGetUserToken.Secret = config.Config.Secret
	openIMGetUserToken.UserID = userID
	openIMGetUserToken.LoginIp = ip
	loginIp := c.Request.Header.Get("X-Forward-For")
	if loginIp == "" {
		loginIp = c.ClientIP()
	}
	openIMGetUserToken.LoginIp = loginIp
	openIMGetUserTokenResp := api.UserTokenResp{}
	bMsg, err := http2.Post(url, openIMGetUserToken, 2)
	if err != nil {
		log.NewError(params.OperationID, "request openIM get user token error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.GetIMTokenErr, "errMsg": err.Error()})
		return
	}
	err = json.Unmarshal(bMsg, &openIMGetUserTokenResp)
	if err != nil || openIMGetUserTokenResp.ErrCode != 0 {
		log.NewError(params.OperationID, "request get user token", account, "err", "")
		if openIMGetUserTokenResp.ErrCode == constant.LoginLimit {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.LoginLimit, "errMsg": "用户登录被限制"})
		} else {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.GetIMTokenErr, "errMsg": ""})
		}
		return
	}
	CallBackOnUserLogin(userID)
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": openIMGetUserTokenResp.UserToken})
}
