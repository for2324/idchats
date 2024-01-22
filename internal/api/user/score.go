package user

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbScore "Open_IM/pkg/proto/score"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetUserScore
// @Summary		获取用户积分
// @Description	获取用户积分
// @Tags			用户相关
// @ID				GetUserScore
// @Accept			json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.GetUserScoreReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetUserScoreResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/get_user_score [post]
func GetUserScore(c *gin.Context) {
	var (
		req   api.GetUserScoreReq
		resp  api.GetUserScoreResp
		reqPb pbScore.GetUserScoreInfoReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	reqPb.OperationID = req.OperationID
	ok, userId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPb.UserId = userId
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.UserScoreName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbScore.NewScoreServiceClient(etcdConn)
	respPb, err := client.GetUserScoreInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.Data = api.UserScoreInfo{UserId: userId, Score: 0}
	if respPb.CommonResp.ErrCode == 0 {
		resp.Data = api.UserScoreInfo{UserId: respPb.UserScoreInfo.UserId, Score: int64(respPb.UserScoreInfo.Score)}
	}
	c.JSON(http.StatusOK, resp)
}

// GetRewardEventLogs
// @Summary		获取用户积分奖励事件日志
// @Description	获取用户积分奖励事件日志
// @Tags			用户相关
// @ID				GetRewardEventLogs
// @Accept			json
// @Param		token	header	string						true	"im token"
// @Param			req		body	api.GetRewardEventLogsReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetRewardEventLogsResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/get_reward_event_logs [post]
func GetRewardEventLogs(c *gin.Context) {
	var (
		req   api.GetRewardEventLogsReq
		resp  api.GetRewardEventLogsResp
		reqPb pbScore.GetRewardEventLogsReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	reqPb.OperationID = req.OperationID
	ok, userId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPb.UserId = userId
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.UserScoreName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbScore.NewScoreServiceClient(etcdConn)
	respPb, err := client.GetRewardEventLogs(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	utils.CopyStructFields(&resp.Data, &respPb.EventLogs)
	c.JSON(http.StatusOK, resp)
}

// WithdrawScore
// @Summary		提现积分
// @Description	提现积分
// @Tags			用户相关
// @ID				WithdrawScore
// @Accept			json
// @Param		token	header	string						true	"im token"
// @Param			req		body	api.WithdrawScoreReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.WithdrawScoreResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/withdraw_score [post]
func WithdrawScore(c *gin.Context) {
	var (
		req   api.WithdrawScoreReq
		resp  api.WithdrawScoreResp
		reqPb pbScore.WithdrawScoreReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	reqPb.OperationID = req.OperationID
	ok, userId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPb.UserId = userId
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.UserScoreName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbScore.NewScoreServiceClient(etcdConn)
	respPb, err := client.WithdrawScore(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.Data = api.WithdrawScoreRespInfo{Id: respPb.WithdrawId}
	c.JSON(http.StatusOK, resp)
}

// GetWithdrawScoreLogs
// @Summary		获取用户积分提现事件日志
// @Description	获取用户积分提现事件日志
// @Tags			用户相关
// @ID				GetWithdrawScoreLogs
// @Accept			json
// @Param		token	header	string						true	"im token"
// @Param			req		body	api.GetWithdrawScoreLogsReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetWithdrawScoreLogsResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/get_withdraw_score_logs [post]
func GetWithdrawScoreLogs(c *gin.Context) {
	var (
		req   api.GetWithdrawScoreLogsReq
		resp  api.GetWithdrawScoreLogsResp
		reqPb pbScore.GetWithdrawScoreLogsReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	reqPb.OperationID = req.OperationID
	ok, userId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPb.UserId = userId
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.UserScoreName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbScore.NewScoreServiceClient(etcdConn)
	respPb, err := client.GetWithdrawScoreLogs(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	utils.CopyStructFields(&resp.Data, &respPb.EventLogs)
	c.JSON(http.StatusOK, resp)
}
