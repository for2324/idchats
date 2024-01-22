package user

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_msg_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	cacheRpc "Open_IM/pkg/proto/cache"
	pbChat "Open_IM/pkg/proto/msg"
	pbRelay "Open_IM/pkg/proto/relay"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	pbTask "Open_IM/pkg/proto/task"
	rpc "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"gopkg.in/gomail.v2"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	gogpt "github.com/sashabaranov/go-openai"
)

//ShowUserProfile

// AskUserQuestion
// @Summary		向openai提问
// @Description	向openai提问
// @Tags		用户相关
// @ID			AskUserQuestion
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.AskQuestionGptReq	true	"用户信息简介"
// @Produce		json
// @Success		0	{object}	api.AskQuestionGptResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/user/ask_user_question [post]
func AskUserQuestion(c *gin.Context) {
	//绑定号码
	var (
		req  api.AskQuestionGptReq
		resp api.AskQuestionGptResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	ok, userId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserSignReq failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	if req.AskMode != "times" && req.AskMode != "tokens" {
		c.JSON(http.StatusOK, gin.H{"errCode": 501, "errMsg": "param mode error"})
		return
	}
	sesstionType := 1
	groupid := ""
	if req.ChannelID != "" && len(req.RecvID) < 20 {
		sesstionType = 3
		groupid = req.RecvID
	} else {
		req.RecvID = userId
	}
	if groupid != "" && req.ChannelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "don't send to group with no channelid"})
		return
	}
	mutexname := "user_chat_token_ask:" + userId
	mutex := db.DB.Pool.NewMutex(mutexname, redsync.WithExpiry(time.Second*10))
	ctx := context.Background()
	if err := mutex.LockContext(ctx); err != nil {
		resp.ErrCode = 10001
		resp.ErrMsg = "the AI is thinking about your last question" + err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	defer mutex.UnlockContext(ctx)
	dbuserinfo, err := rocksCache.GetUserBaseInfoFromCache(userId)
	if err != nil || len(req.CompletionRequest.Messages) == 0 {
		c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": fmt.Sprintf("please to recharge chat token1 %d", len(req.CompletionRequest.Messages))})
		return
	}
	needChatToken := int64(len(req.CompletionRequest.Messages[0].Content)) * 2
	isUseGroup := false
	if req.AskMode == "tokens" {
		if groupid != "" {
			groupinfo, err := rocksCache.GetGroupInfoFromCache(groupid)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": fmt.Sprintf("不存在这个地址")})
				return
			}
			grouplist, err := rocksCache.GetJoinedGroupIDListFromCache(userId)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": fmt.Sprintf("不存在这个地址")})
				return
			}
			if !utils.IsContain(groupinfo.GroupID, grouplist) {
				c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": fmt.Sprintf("you not" +
					" in this group")})
				return
			}
			//群的余额必须小于这个数字 并且 用户余额小雨这个数字。 那么就报厝 {
			//剩余数量大于50 个才是使用
			if groupinfo.ChatTokenCount >= needChatToken && groupinfo.ChatTokenCount >= 50 {
				isUseGroup = true
			} else if int64(dbuserinfo.ChatTokenCount) <= needChatToken || dbuserinfo.ChatTokenCount <= 5 {
				c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": fmt.Sprintf("please to recharge chat tokens")})
				return
			}
		} else if needChatToken > int64(dbuserinfo.ChatTokenCount) ||
			dbuserinfo.ChatTokenCount <= 5 {
			c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": "please to recharge chat tokens"})
			return
		}
	}
	if req.AskMode == "times" && dbuserinfo.ChatCount <= 0 {
		c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": "please to recharge chat ChatCount"})
		return
	}
	proxyAddress, _ := url.Parse("http://proxy.idchats.com:7890")
	conf := gogpt.DefaultConfig(config.Config.ChatGptToken)
	if config.Config.OpenNetProxy.OpenFlag {
		conf.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyAddress),
			},
		}
	}
	client := gogpt.NewClientWithConfig(conf)
	if req.AskMode == "times" {
		req.CompletionRequest.MaxTokens = 2000
	} else {
		req.CompletionRequest.MaxTokens = int(dbuserinfo.ChatTokenCount)
		if dbuserinfo.ChatTokenCount > uint64(config.Config.ChatGptMaxToken) {
			req.CompletionRequest.MaxTokens = config.Config.ChatGptMaxToken
		}
	}
	respdata, err2 := client.CreateChatCompletion(ctx, req.CompletionRequest)
	if err2 != nil {
		resp.ErrCode = 10002
		resp.ErrMsg = ">the AI is thinking about your last question" + err2.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	var rpcUserChatToken rpc.OperatorUserChatTokenReq

	rpcUserChatToken.OperationID = req.OperationID
	rpcUserChatToken.Value = int64(respdata.Usage.TotalTokens)
	if isUseGroup == true {
		rpcUserChatToken.OpUserID = groupid
		rpcUserChatToken.Operator = "groupsub"
	} else {
		rpcUserChatToken.OpUserID = userId
		rpcUserChatToken.Operator = "sub"
	}
	rpcUserChatToken.TxType = req.AskMode
	rpcUserChatToken.ParamStr = utils.StructToJsonString(respdata.Usage)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	clientclient := rpc.NewUserClient(etcdConn)
	RpcResp, err := clientclient.OperatorUserChatToken(context.Background(), &rpcUserChatToken)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), RpcResp.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp.CompletionResponse = respdata
	resp.UserChatToken = RpcResp.NowChatToken
	resp.UserChatCount = RpcResp.NowChatCount
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	text := "我暂时无法思考出这个问题@IDChatsBD"
	if len(respdata.Choices) > 0 {
		text = respdata.Choices[0].Message.Content
	}
	contentType := constant.Text
	contentbyte := utils.String2bytes(text)
	if req.ClientMsgID != "" {
		getmessagefromchatdb, err := im_mysql_msg_model.GetMessageFromMysqlDbByClientID(req.ClientMsgID)
		if err == nil && getmessagefromchatdb != nil {
			apimsgsutrt := new(api.MsgStruct)
			copier.Copy(apimsgsutrt, getmessagefromchatdb)
			QuoteElemValue := &api.QuoteElem{
				Text:         text,
				QuoteMessage: apimsgsutrt,
			}
			contentType = constant.Quote
			contentbyte = utils.StructToJsonBytes(QuoteElemValue)
		}
	} else {
		log.Info(">>>>>>>>>>>>>>>>>>>>>>>>")
	}

	params := api.SendMsgReqStructReq{
		SenderPlatformID: 5,
		SendID:           "0x8858Af738d3F7c33250d7cFd48c89196eA7Dc728",
		SenderNickName:   "Biubiu_AI",
		SenderFaceURL:    "ipfs://QmdsV7cijcBdyfCpCDAUx1jBMfwgy8qDTaBFL4SoKxdfVH",
		OperationID:      utils.Md5(fmt.Sprintf("%d", time.Now().UnixNano())),
		Data: &api.SendMsgReqDataReq{
			SessionType: int32(sesstionType),
			MsgFrom:     100,
			ContentType: int32(contentType),
			RecvID:      req.RecvID,
			GroupID:     groupid,
			ChannelID:   req.ChannelID,
			Content:     contentbyte,
			Options:     nil,
			ClientMsgID: utils.GetMsgID("0x8858Af738d3F7c33250d7cFd48c89196eA7Dc728"),
			CreateTime:  time.Now().UnixMilli(),
		},
	}
	token := c.Request.Header.Get("token")
	log.NewInfo(params.OperationID, "api call success to sendMsgReq", params)
	pbData := &pbChat.SendMsgReq{
		Token:       token,
		OperationID: params.OperationID,
		MsgData: &open_im_sdk.MsgData{
			SendID:           params.SendID,
			RecvID:           params.Data.RecvID,
			GroupID:          params.Data.GroupID,
			ChannelID:        params.Data.ChannelID,
			ClientMsgID:      params.Data.ClientMsgID,
			SenderPlatformID: params.SenderPlatformID,
			SenderNickname:   params.SenderNickName,
			SenderFaceURL:    params.SenderFaceURL,
			SessionType:      params.Data.SessionType,
			MsgFrom:          params.Data.MsgFrom,
			ContentType:      params.Data.ContentType,
			Content:          params.Data.Content,
			CreateTime:       params.Data.CreateTime,
			Options:          params.Data.Options,
		},
	}
	log.Info(params.OperationID, "", "api SendMsg call start..., [data: %s]", pbData.String())
	etcdConn = getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	clientChat := pbChat.NewMsgClient(etcdConn)
	log.Info(params.OperationID, "", "api SendMsg call, api call rpc...")

	resultchatresponse, err := clientChat.SendMsg(context.Background(), pbData)
	if err != nil {
		log.NewError(params.OperationID, "SendMsg rpc failed, ", params, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "SendMsg rpc failed, " + err.Error()})
		return
	}
	if resultchatresponse != nil {
		stringstr, _ := json.Marshal(resultchatresponse)
		log.Info(string(stringstr))
	}

	c.JSON(http.StatusOK, resp)
	return

}

// ChatTokenHistory
// @Summary		查询用户消耗的chattoken的历史记录
// @Description	查询用户消耗的chattoken的历史记录
// @Tags		用户相关
// @ID			ChatTokenHistory
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param		req		body	api.GetUserChatTokenHistoryReq	true	"用户积分请求"
// @Produce		json
// @Success		0	{object}	api.GetUserChatTokenHistoryRsp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/user/user_chat_token_history [post]
func ChatTokenHistory(c *gin.Context) {
	//绑定号码
	var (
		req  api.GetUserChatTokenHistoryReq
		resp api.GetUserChatTokenHistoryRsp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	if req.FromID != "" && !utils.IsDigit(req.FromID) || (req.PageCount != "" && !utils.IsDigit(req.PageCount)) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "paramid error"})
		return
	}

	ok, operatorUid, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserSignReq failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	dbChatResult, err := imdb.GetChatTokenHistory(req.FromID, operatorUid, req.PageCount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	for _, value := range dbChatResult {
		resp.UserChatTokenHistory = append(resp.UserChatTokenHistory, &api.UserChatTokenHistory{
			Action:     value.TxType,
			Param:      value.ParamStr,
			CreateTime: utils.Int64ToString(value.CreatedTime.Unix()),
			ID:         utils.Int64ToString(value.ID),
			ChainID:    value.ChainID,
			Value:      utils.Int64ToString(getCostValue(value.NowCount, value.NewToken, value.OldToken)),
		})
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}
func getCostValue(chatCount, newToken, oldToken uint64) int64 {
	if newToken == oldToken {
		if chatCount == 0 {
			return -1
		} else {
			return 1
		}
	} else {
		return int64(newToken - oldToken)
	}

}

// TransFerTokenToGroup
// @Summary		个人转入空间Token
// @Description	个人转入空间Token
// @Tags		用户相关
// @ID			TransFerTokenToGroup
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param		req		body	api.TransferTokenToGroupReq	true	"转让积分请求"
// @Produce		json
// @Success		0	{object}	api.TransferTokenToGroupRsp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/user/transfer_token_to_group [post]
func TransferTokenToGroup(c *gin.Context) {
	//绑定号码
	var (
		req  api.TransferTokenToGroupReq
		resp api.TransferTokenToGroupRsp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	if req.UserChatToken <= 0 || req.ToGroupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "paramid error"})
		return
	}

	ok, operatorUid, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserSignReq failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	rpcreq := new(rpc.TransferChatTokenOperatorReq)
	rpcreq.ToGroupID = req.ToGroupID
	rpcreq.OpUserID = operatorUid
	rpcreq.OperationID = req.OperationID
	rpcreq.ChatTokenCount = req.UserChatToken
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.TransferChatTokenFromUserToGroup(context.Background(), rpcreq)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), RpcResp.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp.CommResp.ErrCode = RpcResp.CommonResp.ErrCode
	resp.CommResp.ErrMsg = RpcResp.CommonResp.ErrMsg
	resp.ChatTokenCount = make(map[string]int64, 0)
	resp.ChatTokenCount["groupChatTokenCount"] = RpcResp.GroupChatTokenCount
	resp.ChatTokenCount["userChatTokenCount"] = int64(RpcResp.NowChatToken)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}

// GetUserChatToken
// @Summary		查看用户聊天的usage
// @Description	查看用户聊天的usage
// @Tags		用户相关
// @ID			GetUserChatToken
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.GetUserChatTokenReq	true	"用户积分请求"
// @Produce		json
// @Success		0	{object}	api.GetUserChatTokenRsp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/user/get_user_chat_token [post]
func GetUserChatToken(c *gin.Context) {
	//绑定号码
	var (
		req  api.GetUserChatTokenReq
		resp api.GetUserChatTokenRsp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	var operatorUid string
	operatorUid = ""

	ok, operatorUid, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserChatToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	dbUserData, err := rocksCache.GetUserBaseInfoFromCache(operatorUid)
	if err != nil {
		resp.CommResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommResp.ErrMsg = err.Error()
		return
	}
	resp.ChatTokenCount = dbUserData.ChatTokenCount
	resp.ChatCount = dbUserData.ChatCount
	resp.GlobalMoneyCount = dbUserData.GlobalMoneyCount
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}

// OperatorUserChatToken 操作用户+-积分
func OperatorUserChatToken(c *gin.Context) {
	//绑定号码
	var (
		req              api.UserChatTokenReq
		resp             api.UserChatTokenResp
		rpcUserChatToken rpc.OperatorUserChatTokenReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	switch req.Operator {
	case "add":
		if req.TxID == "" || req.ChainID == "" {
			resp.CommResp.ErrCode = constant.ErrInternal.ErrCode
			resp.CommResp.ErrMsg = "没有交易记录，无法添加chat token"
			c.JSON(http.StatusOK, resp)
			return
		}
		rpcUserChatToken.TxID = req.TxID
		rs := db.DB.Pool
		mutexname := "user_chat_token_recharge:" + req.ChainID + ":" + req.TxID
		mutex := rs.NewMutex(mutexname, redsync.WithExpiry(time.Second*10))
		if err := mutex.LockContext(c); err != nil {
			resp.CommResp.ErrCode = constant.ErrInternal.ErrCode
			resp.CommResp.ErrMsg = "正在充值这笔交易"
			c.JSON(http.StatusOK, resp)
			return
		}
		defer mutex.UnlockContext(c)
	case "del":
		if req.Value <= 0 {
			resp.CommResp.ErrCode = constant.ErrInternal.ErrCode
			resp.CommResp.ErrMsg = "param error"
			c.JSON(http.StatusOK, resp)
			return
		}
	default:
		resp.CommResp.ErrCode = constant.ErrInternal.ErrCode
		resp.CommResp.ErrMsg = "不支持的类型参数"
		c.JSON(http.StatusOK, resp)
		return
	}
	rpcUserChatToken.Operator = req.Operator
	rpcUserChatToken.OperationID = req.OperationID
	rpcUserChatToken.Value = req.Value
	rpcUserChatToken.OpUserID = req.UserID
	rpcUserChatToken.ChainID = req.ChainID
	rpcUserChatToken.TxType = req.TxType
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.OperatorUserChatToken(context.Background(), &rpcUserChatToken)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), RpcResp.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp.CommResp.ErrCode = RpcResp.CommonResp.ErrCode
	resp.CommResp.ErrMsg = RpcResp.CommonResp.ErrMsg
	resp.UserChatToken = RpcResp.NowChatToken
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}
func GetUsersInfoFromCache(c *gin.Context) {
	params := api.GetUsersInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "GetUsersInfoFromCache req: ", params)
	req := &rpc.GetUserInfoReq{}
	err := utils.CopyStructFields(req, &params)
	if err != nil {
		return
	}
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	var publicUserInfoList []*open_im_sdk.PublicUserInfo
	for _, v := range RpcResp.UserInfoList {
		publicUserInfoList = append(publicUserInfoList,
			&open_im_sdk.PublicUserInfo{UserID: v.UserID, Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender, Ex: v.Ex})
	}
	resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserInfoList: publicUserInfoList}
	resp.Data = jsonData.JsonDataList(resp.UserInfoList)
	log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}
func GetEachOtherFriendIdListFromCache(c *gin.Context) {
	var (
		req    api.GetFriendIDListFromCacheReq
		resp   api.GetFriendIDListFromCacheResp
		reqPb  cacheRpc.GetFollowEachOtherFriendIDListFromCacheReq
		respPb *cacheRpc.GetFollowEachOtherFriendIDListFromCacheResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)
	reqPb.OperationID = req.OperationID
	var ok bool
	var errInfo string
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := cacheRpc.NewCacheClient(etcdConn)
	respPb, err := client.GetFollowEachOtherFriendIDListFromCache(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFriendIDListFromCache", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed:" + err.Error()})
		return
	}
	resp.UserIDList = respPb.UserIDList
	resp.CommResp = api.CommResp{ErrMsg: respPb.CommonResp.ErrMsg, ErrCode: respPb.CommonResp.ErrCode}
	c.JSON(http.StatusOK, resp)
}

func GetFriendIDListFromCache(c *gin.Context) {
	var (
		req    api.GetFriendIDListFromCacheReq
		resp   api.GetFriendIDListFromCacheResp
		reqPb  cacheRpc.GetFriendIDListFromCacheReq
		respPb *cacheRpc.GetFriendIDListFromCacheResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)
	reqPb.OperationID = req.OperationID
	var ok bool
	var errInfo string
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := cacheRpc.NewCacheClient(etcdConn)
	respPb, err := client.GetFriendIDListFromCache(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFriendIDListFromCache", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed:" + err.Error()})
		return
	}
	resp.UserIDList = respPb.UserIDList
	resp.CommResp = api.CommResp{ErrMsg: respPb.CommonResp.ErrMsg, ErrCode: respPb.CommonResp.ErrCode}
	c.JSON(http.StatusOK, resp)
}

func GetBlackIDListFromCache(c *gin.Context) {
	var (
		req    api.GetBlackIDListFromCacheReq
		resp   api.GetBlackIDListFromCacheResp
		reqPb  cacheRpc.GetBlackIDListFromCacheReq
		respPb *cacheRpc.GetBlackIDListFromCacheResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.OperationID = req.OperationID
	var ok bool
	var errInfo string
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := cacheRpc.NewCacheClient(etcdConn)
	respPb, err := client.GetBlackIDListFromCache(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFriendIDListFromCache", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed:" + err.Error()})
		return
	}
	resp.UserIDList = respPb.UserIDList
	resp.CommResp = api.CommResp{ErrMsg: respPb.CommonResp.ErrMsg, ErrCode: respPb.CommonResp.ErrCode}
	c.JSON(http.StatusOK, resp)
}

// GetUsersPublicInfo
// @Summary		获取用户信息
// @Description	根据用户列表批量获取用户信息
// @Tags			用户相关
// @ID				GetUsersInfo
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.GetUsersInfoReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetUsersInfoResp{Data=[]open_im_sdk.PublicUserInfo}
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/get_users_info [post]
func GetUsersPublicInfo(c *gin.Context) {
	params := api.GetUsersInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetUserInfoReq{}
	utils.CopyStructFields(req, &params)
	var SelectIndex []string
	for _, value := range params.UserIDList {
		SelectIndex = append(SelectIndex, strings.ToLower(value))
	}
	req.UserIDList = SelectIndex

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(params.OperationID, "GetUserInfo args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	var publicUserInfoList []*open_im_sdk.PublicUserInfo
	for _, v := range RpcResp.UserInfoList {
		publicUserInfoList = append(publicUserInfoList,
			&open_im_sdk.PublicUserInfo{UserID: v.UserID,
				Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender, Ex: v.Ex,
				TokenContractChain: v.TokenContractChain})
	}

	resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserInfoList: publicUserInfoList}
	resp.Data = jsonData.JsonDataList(resp.UserInfoList)
	log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// UpdateUserInfo
// @Summary		修改用户信息
// @Description	修改用户信息 userID faceURL等
// @Tags			用户相关
// @ID				UpdateUserInfo
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.UpdateSelfUserInfoReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UpdateUserInfoResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/update_user_info [post]
func UpdateUserInfo(c *gin.Context) {
	params := api.UpdateSelfUserInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.UpdateUserInfoReq{UserInfo: &open_im_sdk.UserInfo{}}
	utils.CopyStructFields(req.UserInfo, &params)
	fmt.Println("\n", params.UserID, params.FaceURL, "\n")
	req.OperationID = params.OperationID
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return

	}

	log.NewInfo(params.OperationID, "UpdateUserInfo args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.UpdateUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "UpdateUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.UpdateUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "UpdateUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// UpdateUserHead
// @Summary		更改用户头像
// @Description	更改用户头像 userID faceURL等
// @Tags			用户相关
// @ID				UpdateUserHead
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.UpdateSelfUserHeadReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UpdateSelfUserHeadResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/update_user_head [post]
func UpdateUserHead(c *gin.Context) {
	params := api.UpdateSelfUserHeadReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	ok, Uid, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	TokenUrl, boolFlag, err := CheckParamUserHeadInfo(params.OperationID, Uid, params.NftContract, params.TokenID, params.NftChainID)
	if !boolFlag {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if err == nil && TokenUrl == "" {
		TokenUrl = params.NftTokenURL
	}
	req := &rpc.UpdateUserInfoHeadReq{UserInfo: &open_im_sdk.UserInfo{}}
	//utils.CopyStructFields(req.UserInfo, )
	req.OperationID = params.OperationID
	req.OpUserID = Uid
	req.UserInfo.UserID = Uid
	req.UserInfo.FaceURL = TokenUrl
	if params.NftContract != "" && params.NftChainID != "" {
		req.ContractChain = params.NftContract + "&" + params.NftChainID
	}
	req.UserInfo.TokenId = params.TokenID

	log.NewInfo(params.OperationID, "UpdateUserInfo args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.UpdateUserInfoHead(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "UpdateUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.UpdateSelfUserHeadResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "UpdateUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// SetGlobalRecvMessageOpt
// @Summary		全局免打扰设置
// @Description	全局免打扰设置
// @Tags			用户相关
// @ID				SetGlobalRecvMessageOpt
// @Accept			json
// @Param			token	header	string							true	"im token"
// @Param			req		body	api.SetGlobalRecvMessageOptReq	true	"globalRecvMsgOpt为接收全局推送设置0为关闭 1为开启"
// @Produce		json
// @Success		0	{object}	api.SetGlobalRecvMessageOptResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/set_global_msg_recv_opt [post]
func SetGlobalRecvMessageOpt(c *gin.Context) {
	params := api.SetGlobalRecvMessageOptReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.SetGlobalRecvMessageOptReq{}
	utils.CopyStructFields(req, &params)
	req.OperationID = params.OperationID
	var ok bool
	var errInfo string
	ok, req.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, "SetGlobalRecvMessageOpt args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	req.FieldName = "global_recv_msg_opt"
	RpcResp, err := client.RpcUpdateUserFieldData(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "SetGlobalRecvMessageOpt failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.UpdateUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "SetGlobalRecvMessageOpt api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// SetShowBalance
// @Summary		设置显示用户的资产
// @Description	设置显示用户的资产
// @Tags			用户相关
// @ID				SetShowBalance
// @Accept			json
// @Param			token	header	string							true	"im token"
// @Param			req		body	api.SetGlobalRecvMessageOptReq	true	"globalRecvMsgOpt为全局免打扰设置0为关闭 1为开启"
// @Produce		json
// @Success		0	{object}	api.SetGlobalRecvMessageOptResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/set_show_balance [post]
func SetShowBalance(c *gin.Context) {
	params := api.SetGlobalRecvMessageOptReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.SetGlobalRecvMessageOptReq{}
	utils.CopyStructFields(req, &params)
	req.OperationID = params.OperationID
	var ok bool
	var errInfo string
	ok, req.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, "SetShowBalance args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	req.FieldName = "show_balance"
	req.GlobalRecvMsgOpt = *params.ShowBalance
	RpcResp, err := client.RpcUpdateUserFieldData(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "SetShowBalance failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.UpdateUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "SetShowBalance api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// GetSelfUserInfo
// @Summary		获取自己的信息
// @Description	传入ID获取自己的信息
// @Tags			用户相关
// @ID				GetSelfUserInfo
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetSelfUserInfoReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetSelfUserInfoResp{data=open_im_sdk.UserInfo}
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/get_self_user_info [post]
func GetSelfUserInfo(c *gin.Context) {
	params := api.GetSelfUserInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := " BindJSON failed " + err.Error()
		log.NewError("0", errMsg)
		c.JSON(http.StatusOK, gin.H{"errCode": 1001, "errMsg": errMsg})
		return
	}
	req := &rpc.GetSelfUserInfoReq{}

	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusOK, gin.H{"errCode": 1001, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, "GetUserInfo args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	var wg sync.WaitGroup
	var userInfo *open_im_sdk.UserInfo
	var groupInfo *db.Group
	var err error
	wg.Add(2)
	go func() {
		defer wg.Done()
		client := rpc.NewUserClient(etcdConn)
		RpcResp, rpcErr := client.GetSelfUserInfo(context.Background(), req)
		if rpcErr != nil {
			err = errors.New("call  rpc server failed:" + rpcErr.Error())
			log.NewError(req.OperationID, "GetUserInfo failed ", rpcErr.Error(), req.String())
			return
		}
		if RpcResp.CommonResp.ErrCode != 0 {
			err = errors.New(RpcResp.CommonResp.ErrMsg)
			log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
			return
		}
		userInfo = RpcResp.UserInfo
		// userInfoReq := &rpc.GetUserInfoReq{
		// 	OperationID: req.OperationID,
		// 	UserIDList:  []string{req.UserID},
		// 	OpUserID:    req.UserID,
		// }
		// RpcResp, rpcErr := client.GetUserInfo(context.Background(), userInfoReq)
		// if rpcErr != nil {
		// 	err = errors.New("call  rpc server failed:" + rpcErr.Error())
		// 	log.NewError(req.OperationID, "GetUserInfo failed ", rpcErr.Error(), req.String())
		// 	return
		// }
		// if len(RpcResp.UserInfoList) == 1 {
		// 	userInfo = RpcResp.UserInfoList[0]
		// }
	}()
	// get groupInfo
	go func() {
		defer wg.Done()
		rGroupInfo, err := rocksCache.GetSpaceInfoByUser(req.UserID)
		if err != nil {
			log.NewError(req.OperationID, "GetSpaceInfoByUser failed ", err.Error(), req.String())
			return
		}
		groupInfo = rGroupInfo
	}()
	wg.Wait()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrInternal.ErrCode, "errMsg": err.Error()})
		return
	}
	respData := api.ApiSelfUserInfo{
		GroupInfo: &open_im_sdk.GroupInfo{},
	}
	utils.CopyStructFields(&respData, userInfo)
	utils.CopyStructFields(&respData.GroupInfo, groupInfo)
	log.NewInfo(req.OperationID, "GetUserInfo api return ", respData)
	c.JSON(http.StatusOK, api.GetSelfUserInfoResp{CommResp: api.CommResp{}, Data: respData})
}

// GetUsersOnlineStatus
// @Summary		获取用户在线状态
// @Description	获取用户在线状态
// @Tags			用户相关
// @ID				GetUsersOnlineStatus
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.GetUsersOnlineStatusReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetUsersOnlineStatusResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/get_users_online_status [post]
func GetUsersOnlineStatus(c *gin.Context) {
	params := api.GetUsersOnlineStatusReq{}
	if err := c.BindJSON(&params); err != nil {

		c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbRelay.GetUsersOnlineStatusReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	if len(config.Config.Manager.AppManagerUid) == 0 {
		log.NewError(req.OperationID, "Manager == 0")
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "Manager == 0"})
		return
	}
	req.OpUserID = config.Config.Manager.AppManagerUid[0]

	log.NewInfo(params.OperationID, "GetUsersOnlineStatus args ", req.String())
	var wsResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult
	flag := false
	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), params.OperationID)
	for _, v := range grpcCons {
		log.Debug(params.OperationID, "get node ", *v, v.Target())
		client := pbRelay.NewRelayClient(v)
		reply, err := client.GetUsersOnlineStatus(context.Background(), req)
		if err != nil {
			log.NewError(params.OperationID, "GetUsersOnlineStatus rpc  err", req.String(), err.Error())
			continue
		} else {
			if reply.ErrCode == 0 {
				wsResult = append(wsResult, reply.SuccessResult...)
			}
		}
	}
	log.NewInfo(params.OperationID, "call GetUsersOnlineStatus rpc server is success", wsResult)
	//Online data merge of each node
	for _, v1 := range params.UserIDList {
		flag = false
		temp := new(pbRelay.GetUsersOnlineStatusResp_SuccessResult)
		for _, v2 := range wsResult {
			if v2.UserID == v1 {
				flag = true
				temp.UserID = v1
				temp.Status = constant.OnlineStatus
				temp.DetailPlatformStatus = append(temp.DetailPlatformStatus, v2.DetailPlatformStatus...)
			}

		}
		if !flag {
			temp.UserID = v1
			temp.Status = constant.OfflineStatus
		}
		respResult = append(respResult, temp)
	}
	resp := api.GetUsersOnlineStatusResp{CommResp: api.CommResp{ErrCode: 0, ErrMsg: ""}, SuccessResult: respResult}
	if len(respResult) == 0 {
		resp.SuccessResult = []*pbRelay.GetUsersOnlineStatusResp_SuccessResult{}
	}
	log.NewInfo(req.OperationID, "GetUsersOnlineStatus api return", resp)
	c.JSON(http.StatusOK, resp)
}

func GetUsers(c *gin.Context) {
	var (
		req   api.GetUsersReq
		resp  api.GetUsersResp
		reqPb rpc.GetUsersReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, _, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.OperationID = req.OperationID
	reqPb.UserID = req.UserID
	reqPb.UserName = req.UserName
	reqPb.Content = req.Content
	reqPb.Pagination = &open_im_sdk.RequestPagination{ShowNumber: req.ShowNumber, PageNumber: req.PageNumber}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.GetUsers(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	for _, v := range respPb.UserList {
		user := api.CMSUser{}
		utils.CopyStructFields(&user, v.User)
		user.IsBlock = v.IsBlock
		resp.Data.UserList = append(resp.Data.UserList, &user)
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.Data.TotalNum = respPb.TotalNums
	resp.Data.CurrentPage = respPb.Pagination.CurrentPage
	resp.Data.ShowNumber = respPb.Pagination.ShowNumber
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}

// GetThirdInfo
// @Summary		获取用户绑定的第三方信息信息
// @Description	根据用户列表批量获取用户绑定第三方信息 twitter facebook
// @Tags			用户相关
// @ID				GetThirdInfo
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.GetUsersThirdInfoReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetUsersThirdInfoResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/get_third_status [post]
func GetThirdInfo(c *gin.Context) {
	var (
		req   api.GetUsersThirdInfoReq
		resp  api.GetUsersThirdInfoResp
		reqPb rpc.GetUserThirdInfoReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	chainid := c.Request.Header.Get("chainId")
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.OperationID = req.OperationID
	reqPb.UserList = req.UserIDList
	reqPb.ChainID = chainid
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.GetUserThird(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.Data = jsonData.JsonDataList(respPb.UserThirdInfoList)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}

// BindUserEmail
// @Summary		用户绑定邮箱
// @Description	用户绑定邮箱
// @Tags			用户相关
// @ID				BindUserEmail
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.BindUserSelfDomainReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.BindUserSelfDomainResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/bind_user_email [post]
func BindUserEmail(c *gin.Context) {
	var (
		req   api.BindUserSelfDomainReq
		resp  api.BindUserSelfDomainResp
		reqPb rpc.BindUserThirdInfoReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	regster := req.EmailAddress + ":" + reqPb.OpUserID
	//如果找不懂到验证码直接报错
	if !db.DB.ExistVerifyCode(regster) {
		c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "your email verify code is expire"})
		return
	} else {
		codestring, err := db.DB.GetEmailVerifyCode(regster)
		if err != nil || req.EmailVerifyCode != codestring {
			c.JSON(http.StatusOK, gin.H{"errCode": 401, "errMsg": "verify  code error"})
			return
		}

		db.DB.DeleteEmailVerifyCode(regster)
	}

	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.OperationID = req.OperationID
	reqPb.EmailAddress = req.EmailAddress
	reqPb.Action |= 1 << 1
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.BindUserThirdInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}

// BindDnsDomain
// @Summary		绑定第三方信息ens域名信息
// @Description	绑定第三方信息ens域名信息
// @Tags			用户相关
// @ID				BindDnsDomain
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.BindUserSelfDomainReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.BindUserSelfDomainResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/bind_web_domain [post]
func BindDnsDomain(c *gin.Context) {
	var (
		req   api.BindUserSelfDomainReq
		resp  api.BindUserSelfDomainResp
		reqPb rpc.BindUserThirdInfoReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.OperationID = req.OperationID
	reqPb.Domain = req.DnsDomain
	reqPb.Action |= 0x01
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.BindUserThirdInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}

// BindEnsDomain
// @Summary		绑定第三方信息ens域名信息
// @Description	绑定第三方信息ens域名信息
// @Tags			用户相关
// @ID				GetUsersThirdInfoReq
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.BindUserSelfDomainReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.BindUserSelfDomainResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/bind_ens_domain [post]
func BindEnsDomain(c *gin.Context) {
	var (
		req   api.BindUserSelfDomainReq
		resp  api.BindUserSelfDomainResp
		reqPb rpc.BindUserEnsDomainReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	chainid := c.Request.Header.Get("chainId")
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.OperationID = req.OperationID
	reqPb.EnsDomain = req.EnsDomain
	reqPb.ChainID = chainid
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.BindUserEnsDomain(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}

// BindUserInfoTelephoneInfo
// @Summary		绑定用户的电话号码
// @Description	绑定用户的电话号码
// @Tags	用户相关
// @ID		BindUserInfoTelephoneInfo
// @Accept	json
// @Param	token header	string			true	"im token"
// @Param	req	 body	api.BindUserTelephoneReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.BindUserTelephoneResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/updateTelephoneInfo [post]
func BindUserInfoTelephoneInfo(c *gin.Context) {
	//绑定号码
	var (
		req   api.BindUserTelephoneReq
		resp  api.BindUserTelephoneResp
		reqPb rpc.BindUserTelephoneReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	//判断验证码是否等于他提交的验证码

	var account string
	if req.Email != "" {
		account = req.Email
	} else {
		account = req.PhoneNumber
	}
	var accountKey = req.AreaCode + account + reqPb.OpUserID + constant.BindTelePhoneNumber
	accountKeyOldPhone := "OldPhone:" + reqPb.OpUserID + constant.ChangeTelePhoneNumber
	if req.UpdateSecret != "" {
		code, _ := db.DB.GetAccountCode(accountKeyOldPhone)
		if code != req.UpdateSecret {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.ResetPasswordFailed, "errMsg": "Old phone verification code expired!"})
			return
		}
	}
	code, err := db.DB.GetAccountCode(accountKey)
	log.NewInfo(req.OperationID, "redis phone number and verificating Code",
		"key: ", accountKey, "code: ", code, "params: ", req)
	if err != nil {
		log.NewError(req.OperationID, "Verification code expired", accountKey, "err", err.Error())
		data := make(map[string]interface{})
		data["account"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.CodeInvalidOrExpired, "errMsg": "Verification code expired!", "data": data})
		return
	}

	if req.TelephoneCode == "909090" || req.TelephoneCode == code {
		log.Info(req.OperationID, "Verified successfully", account)
		data := make(map[string]interface{})
		data["account"] = account
		data["telephoneCode"] = req.TelephoneCode
		reqPb.OperationID = req.OperationID
		reqPb.Email = req.Email
		reqPb.Telephone = req.PhoneNumber
		if req.UpdateSecret != "" {
			reqPb.IsUpdatePhone = true
		}
		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
		if etcdConn == nil {
			errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
			log.NewError(reqPb.OperationID, errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
			return
		}
		client := rpc.NewUserClient(etcdConn)
		respPb, err := client.BindUserTelephoneRPC(context.Background(), &reqPb)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
			return
		}
		resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
		resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
		db.DB.DelAccountCode(accountKey)
		db.DB.DelAccountCode(accountKeyOldPhone)
		c.JSON(http.StatusOK, resp)
		return
	} else {
		log.Info(req.OperationID, "Verification code error", account, req.TelephoneCode)
		data := make(map[string]interface{})
		data["account"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.CodeInvalidOrExpired, "errMsg": "Verification code error!", "data": data})
	}
	return
}

// BindUserTelephoneCode
// @Summary		绑定用户的电话号码发送验证码
// @Description	绑定用户的电话号码发送验证码
// @Tags			用户相关
// @ID				BindUserTelephoneCode
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.BindUserTelephoneCodeReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.BindUserTelephoneResp "usedfor 一般为3"
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/sendSms [post]
func BindUserTelephoneCode(c *gin.Context) {
	req := api.BindUserTelephoneCodeReq{}
	if err := c.BindJSON(&req); err != nil {
		log.NewError("", "BindJSON failed", "err:", err.Error(), "phoneNumber", req.PhoneNumber, "email", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	operationID := req.OperationID
	if operationID == "" {
		operationID = utils.OperationIDGenerator()
	}
	ok, uid, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.Info(operationID, "SendVerificationCode args: ", "area code: ",
		req.AreaCode, "Phone Number: ", req.PhoneNumber)
	var account string
	if req.Email != "" {
		account = req.Email
	} else {
		account = req.PhoneNumber
	}
	var accountKey = req.AreaCode + account + uid
	if req.UsedFor == 0 {
		req.UsedFor = constant.VerificationCodeForRegister
	}
	switch req.UsedFor {
	case constant.BindTelephoneNumber:
		accountKey = accountKey + constant.BindTelePhoneNumber
		fmt.Println("\n accountKey:", accountKey)
		ok, err := db.DB.JudgeAccountEXISTS(accountKey)
		if err != nil {
			log.NewError(req.OperationID, "Repeat send code", req, accountKey)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.RepeatSendCode, "errMsg": "Repeat send code"})
			return
		}
		if ok {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verification code has been set!", "data": ""})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": "Error UsedFor"})
		return
	}
	rand.Seed(time.Now().UnixNano())
	code := 100000 + rand.Intn(999999)
	log.NewInfo(req.OperationID, req.UsedFor, "begin store redis", accountKey, code)
	levelTime := config.Config.Demo.CodeTTL * 100
	if !config.Config.OpenNetProxy.OpenFlag {
		levelTime = config.Config.Demo.CodeTTL
	}
	err := db.DB.SetAccountCode(accountKey, code, levelTime)
	if err != nil {
		log.NewError(req.OperationID, "set redis error", accountKey, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
		return
	}
	log.NewDebug(req.OperationID, config.Config.Demo)
	if req.Email != "" {
		m := gomail.NewMessage()
		m.SetHeader(`From`, config.Config.Demo.Mail.SenderMail)
		m.SetHeader(`To`, []string{account}...)
		m.SetHeader(`Subject`, config.Config.Demo.Mail.Title)
		m.SetBody(`text/html`, fmt.Sprintf("%d", code))
		if err := gomail.NewDialer(config.Config.Demo.Mail.SmtpAddr, config.Config.Demo.Mail.SmtpPort, config.Config.Demo.Mail.SenderMail, config.Config.Demo.Mail.SenderAuthorizationCode).DialAndSend(m); err != nil {
			log.Error(req.OperationID, "send mail error", account, err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": constant.MailSendCodeErr, "errMsg": ""})
			return
		}
	} else {
		//发送短信
		err := SendYunPian(operationID, req.AreaCode+req.PhoneNumber, config.Config.Demo.Yunpiansms.Templateid, utils.Int32ToString(int32(code)))
		if err != nil {
			log.NewError(req.OperationID, "sendSms error", account, "err", err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "ErrorCode SMS"})
			return
		}
	}
	log.Debug(req.OperationID, "send sms success", code, accountKey)
	data := make(map[string]interface{})
	data["account"] = account
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verification code has been set!", "data": data})
}

// ChangeOldPhoneNumber
// @Summary		修改发送旧的手机号码只提交接口旧可以
// @Description	修改发送旧的手机号码只提交接口旧可以 userfor 请传输4
// @Tags			用户相关
// @ID				ChangeOldPhoneNumber
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.BindUserTelephoneCodeReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.BindUserTelephoneResp "usedfor 一般为3"
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/change_old_phone_number [post]
func ChangeOldPhoneNumber(c *gin.Context) {
	req := api.BindUserTelephoneCodeReq{}
	if err := c.BindJSON(&req); err != nil {
		log.NewError("", "BindJSON failed", "err:", err.Error(), "phoneNumber", req.PhoneNumber, "email", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	operationID := req.OperationID
	if operationID == "" {
		operationID = utils.OperationIDGenerator()
	}
	ok, uid, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	telePhoneNumberUserInfo, _ := imdb.GetUserByUserID(uid)
	req.AreaCode = "86"
	req.PhoneNumber = telePhoneNumberUserInfo.PhoneNumber

	ccountKeyOldPhone := "OldPhone:" + uid + constant.ChangeTelePhoneNumber
	switch req.UsedFor {
	case constant.ChangeTelePhoneCode:
		ok, err := db.DB.JudgeAccountEXISTS(ccountKeyOldPhone)
		if err != nil {
			log.NewError(req.OperationID, "Repeat send code", req, ccountKeyOldPhone)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.RepeatSendCode, "errMsg": "Repeat send code"})
			return
		}
		if ok {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verification code has been set!", "data": ""})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": "Error UsedFor"})
		return
	}
	rand.Seed(time.Now().UnixNano())
	code := 100000 + rand.Intn(999999)
	log.NewInfo(req.OperationID, req.UsedFor, "begin store redis", ccountKeyOldPhone, code)
	levelTime := config.Config.Demo.CodeTTL * 100
	if !config.Config.OpenNetProxy.OpenFlag {
		levelTime = config.Config.Demo.CodeTTL
	}
	err := db.DB.SetAccountCode(ccountKeyOldPhone, code, levelTime)
	if err != nil {
		fmt.Println(req.OperationID, "set redis error", ccountKeyOldPhone, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "insert into redis error"})
		return
	}
	log.NewDebug(req.OperationID, config.Config.Demo)
	err = SendYunPian(operationID, req.AreaCode+req.PhoneNumber, config.Config.Demo.Yunpiansms.Templateid, utils.Int32ToString(int32(code)))
	if err != nil {
		log.NewError(req.OperationID, "sendSms error", telePhoneNumberUserInfo.PhoneNumber, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "ErrorCode SMS:" + err.Error()})
		return
	}

	log.Debug(req.OperationID, "send sms success", code, ccountKeyOldPhone)
	data := make(map[string]interface{})
	data["account"] = telePhoneNumberUserInfo.PhoneNumber
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verification code has been set!", "data": data})
}

// DeletePlatformInfo
// @Summary		删除用户的绑定的信息
// @Description	删除用户的绑定的信息
// @Tags		用户相关
// @ID			DeletePlatformInfo
// @Accept		json
// @Param		token	header	string						true	"im token"
// @Param		req		body	api.DelThirdPlatformReq	true	"请求体，twitter phone dnsDomain ensDomain facebook faceURL"
// @Produce		json
// @Success		0	{object}	api.DelThirdPlatformResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/user/delThirdPlatform [post]
func DeletePlatformInfo(c *gin.Context) {
	//绑定号码
	var (
		req   api.DelThirdPlatformReq
		resp  api.DelThirdPlatformResp
		reqPb rpc.DelThirdPlatformReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	chainid := c.Request.Header.Get("chainId")

	reqPb.OperationID = req.OperationID
	reqPb.PlatformName = req.PlatformName
	reqPb.ChainID = chainid
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.DeletePlatformInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}

// ShowPlatfomrInfo
// @Summary		显示或者隐藏第三方平台信息
// @Description	显示或者隐藏第三方平台信息
// @Tags		用户相关
// @ID			ShowPlatfomrInfo
// @Accept		json
// @Param		token	header	string						true	"im token"
// @Param		req		body	api.ShowThirdPlatformReq	true	"请求体，twitter phone dnsDomain ensDomain facebook "
// @Produce		json
// @Success		0	{object}	api.ShowThirdPlatformResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/user/show_third_platform [post]
func ShowPlatfomrInfo(c *gin.Context) {
	//绑定号码
	var (
		req   api.ShowThirdPlatformReq
		resp  api.ShowThirdPlatformResp
		reqPb rpc.ShowThirdPlatformReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	chainid := c.Request.Header.Get("chainId")

	reqPb.OperationID = req.OperationID
	reqPb.PlatformName = req.PlatformName
	reqPb.ChainID = chainid
	reqPb.ShowFlag = req.ShowFlag
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.ShowPlatformInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}

// 云片短信  	PhoneNumbers:要发送的手机号码，多个号码用逗号隔开
func SendYunPian(operationID, PhoneNumbers, TemplateCode, TemplateParam string) error {
	//请求地址
	smsURL := "https://sms.yunpian.com/v2/sms/tpl_single_send.json"
	apikey := config.Config.Demo.Yunpiansms.Appid
	tplValue := url.Values{"#code#": {TemplateParam}}.Encode()
	if !strings.HasPrefix(PhoneNumbers, "+") {
		PhoneNumbers = "+" + PhoneNumbers
	}
	param := url.Values{"apikey": {apikey}, "mobile": {PhoneNumbers},
		"tpl_id": {TemplateCode}, "tpl_value": {tplValue}}
	log.NewInfo(operationID, "G云片短信参数信息:: ", param)
	////发送请求
	data, err := utils.HttpPostForm(smsURL, param)
	if err != nil {
		fmt.Println("云片短信请求失败,错误信息:\r\n", err.Error())
		return err
	} else {
		type YunPianResp struct {
			Code        int
			Msg         string
			MobilePhone string
		}
		var resultJson YunPianResp
		json.Unmarshal(data, &resultJson)
		Code := resultJson.Code
		Phones := resultJson.MobilePhone
		if Code == 0 {
			fmt.Println("云片短信发送成功:", Phones)
		} else {
			msg := resultJson.Msg
			fmt.Println("云片短信发送失败Code:,msg:, phones:", Code, msg, Phones)
			return errors.New("云片短信发送失败:" + utils.StructToJsonString(resultJson))
		}
	}
	return nil
}

// UpdateUserSign
// @Summary		设置用户的简介
// @Description	设置用户的简介
// @Tags		用户相关
// @ID			UpdateUserSign
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.UpdateUserSignReq	true	"用户信息简介"
// @Produce		json
// @Success		0	{object}	api.UpdateUserSignResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/user/update_user_sign [post]
func UpdateUserSign(c *gin.Context) {
	//绑定号码
	var (
		req  api.UpdateUserSignReq
		resp api.UpdateUserSignResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	var operatorUid string
	ok, operatorUid, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	mapValue := make(map[string]interface{}, 0)
	mapValue["user_profile"] = req.UserProfile
	err := imdb.UpdateUserInfoWithMapping(operatorUid, mapValue)
	if err != nil {
		resp.CommResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommResp.ErrMsg = err.Error()
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.CommResp.ErrCode = 0
	resp.CommResp.ErrMsg = ""
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
	return
}

// GetUserSign
// @Summary		获取个人简介信息
// @Description	获取个人简介信息
// @Tags		用户相关
// @ID			GetUserSign
// @Accept		json
// @Param		token	header	string						true	"im token"
// @Param		req		body	api.GetUserSignReq	true	"123123"
// @Produce		json
// @Success		0	{object}	api.GetUserSignResp	""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/user/get_user_sign [post]
func GetUserSign(c *gin.Context) {
	//绑定号码
	var (
		req  api.GetUserSignReq
		resp api.GetUserSignResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	//var ok bool
	//var errInfo string
	dbuserdata, err := imdb.GetUserByUserID(req.UserId)
	if err != nil {
		resp.CommResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommResp.ErrMsg = err.Error()
		return
	}
	utils.CopyStructFields(resp.UserInfo, dbuserdata)
	resp.Data = jsonData.JsonDataOne(resp.UserInfo)
	resp.Data["userProfile"] = dbuserdata.UserProfile
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
}

func FinishUploadNftHeadTask(OperationID, userId string) {
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, OperationID)
	if etcdConn == nil {
		return
	}
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.FinishUploadNftHeadTask(context.Background(), &pbTask.FinishUploadNftHeadTaskReq{
		OperationID: OperationID,
		UserId:      userId,
	})
	if err != nil {
		log.NewError(OperationID, "FinishUploadNftHeadTask failed ", err.Error(), userId)
		return
	}
	if respPb.CommonResp.ErrCode != 0 {
		log.NewDebug(OperationID, "FinishUploadNftHeadTask failed ", respPb.CommonResp.ErrMsg, userId)
		return
	}
	log.NewInfo(OperationID, "FinishUploadNftHeadTask success ", userId)
}
func CloseUploadNftHeadTask(OperationID, userId string) {
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, OperationID)
	if etcdConn == nil {
		return
	}
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.CloseUploadNftHeadTask(context.Background(), &pbTask.CloseUploadNftHeadTaskReq{
		OperationID: OperationID,
		UserId:      userId,
	})
	if err != nil {
		log.NewError(OperationID, "CloseUploadNftHeadTask failed ", err.Error(), userId)
		return
	}
	if respPb.CommonResp.ErrCode != 0 {
		log.NewDebug(OperationID, "CloseUploadNftHeadTask failed ", respPb.CommonResp.ErrMsg, userId)
		return
	}
	log.NewInfo(OperationID, "CloseUploadNftHeadTask success ", userId)
}

func FinishOfficialNFTHeadTask(OperationID, userId string) {
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, OperationID)
	if etcdConn == nil {
		return
	}
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.FinishOfficialNFTHeadTask(context.Background(), &pbTask.FinishOfficialNFTHeadTaskReq{
		OperationID: OperationID,
		UserId:      userId,
	})
	if err != nil {
		log.NewError(OperationID, "FinishOfficialNFTHeadTask failed ", err.Error(), userId)
		return
	}
	if respPb.CommonResp.ErrCode != 0 {
		log.NewDebug(OperationID, "FinishOfficialNFTHeadTask failed ", respPb.CommonResp.ErrMsg, userId)
		return
	}
	log.NewInfo(OperationID, "FinishOfficialNFTHeadTask success ", userId)
}

func CloseOfficialNFTHeadTask(OperationID, userId string) {
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, OperationID)
	if etcdConn == nil {
		return
	}
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.CloseOfficialNFTHeadTask(context.Background(), &pbTask.CloseOfficialNFTHeadTaskReq{
		OperationID: OperationID,
		UserId:      userId,
	})
	if err != nil {
		log.NewError(OperationID, "CloseOfficialNFTHeadTask failed ", err.Error(), userId)
		return
	}
	if respPb.CommonResp.ErrCode != 0 {
		log.NewDebug(OperationID, "CloseOfficialNFTHeadTask failed ", respPb.CommonResp.ErrMsg, userId)
		return
	}
	log.NewInfo(OperationID, "CloseOfficialNFTHeadTask success ", userId)
}
func CheckParamUserHeadInfo(OperationID, UserID, NftContract, TokenID, NftChainID string) (string, bool, error) {
	officialNftContractList, err := rocksCache.GetOfficialNftContractFromCache()
	log.NewInfo(OperationID, "UpdateUserInfo args<>>>>>>>>>>>>>>>> ", officialNftContractList)
	if err != nil {
		return "", false, errors.New("暂未设置官方nft")
	}
	var RequestTokenIdRespdata api.RequestTokenIdResp
	if NftContract != "" && TokenID != "" && NftChainID != "" {
		if strings.EqualFold(NftContract, "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee") && TokenID == "0" {
			return "", true, nil
		}

		PostCheckData := &api.RequestTokenIdReq{
			ContractAddress: NftContract,
			TokenID:         TokenID,
			ChainID:         NftChainID,
		}
		jsonbyte, _ := json.Marshal(PostCheckData)

		resultByte, err := utils.HttpPost(config.Config.EnsPostCheck.Url+"/graph/requesttokenuri",
			"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, jsonbyte)
		if err != nil {
			return "", false, err
		}
		fmt.Println("Post data :", NftContract, ">> NftChainID", NftChainID, "TokenID:>>", TokenID, "返回结果:", string(resultByte))
		json.Unmarshal(resultByte, &RequestTokenIdRespdata)
		if RequestTokenIdRespdata.ErrCode != 0 || !strings.EqualFold(RequestTokenIdRespdata.TokenOwnerAddress, UserID) {
			return "", false, errors.New(RequestTokenIdRespdata.ErrMsg)
		}
	}

	// 完成上传官方头像任务
	if utils.IsContainEqual(NftContract, officialNftContractList) {
		FinishOfficialNFTHeadTask(OperationID, UserID)
	} else {
		// 取消上传官方头像任务
		CloseOfficialNFTHeadTask(OperationID, UserID)
	}
	// 完成上传头像任务
	if NftContract != "" && NftChainID != "" {
		// 完成上传头像任务
		FinishUploadNftHeadTask(OperationID, UserID)
	} else {
		// 取消上传头像任务
		CloseUploadNftHeadTask(OperationID, UserID)
	}

	return RequestTokenIdRespdata.TokenUrl, true, nil

}
