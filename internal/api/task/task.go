package task

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbTask "Open_IM/pkg/proto/task"
	pbWeb3 "Open_IM/pkg/proto/web3pub"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetTaskList
// @Summary		获取全部任务
// @Description	获取全部任务
// @Tags			任务
// @ID				GetTaskList
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetTaskListReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.GetUserTaskListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/task/get_task_list [post]
func GetTaskList(c *gin.Context) {
	var (
		req   api.GetTaskListReq
		resp  api.GetTaskListResp
		reqPb pbTask.GetTaskListReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	utils.CopyStructFields(&reqPb, req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.GetTaskList(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "Create Task failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	for _, v := range respPb.Data {
		taskInfo := api.Task{}
		utils.CopyStructFields(&taskInfo, v)
		resp.Data = append(resp.Data, &taskInfo)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", respPb)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp.Data})
}

// GetUserTaskList
// @Summary		获取用户可进行的任务
// @Description	获取用户可进行的任务
// @Tags			任务
// @ID				GetUserTaskList
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetUserTaskListReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.GetUserTaskListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/task/get_user_task_list [post]
func GetUserTaskList(c *gin.Context) {
	var (
		req   api.GetUserTaskListReq
		resp  api.GetUserTaskListResp
		reqPb pbTask.GetUserTaskListReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	utils.CopyStructFields(&reqPb, req)
	// get user id from token
	var ok bool
	var errInfo string
	ok, reqPb.UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	// etcd im task
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// rpc
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.GetUserTaskList(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "Create Task failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", respPb)
	for _, v := range respPb.Data {
		taskInfo := api.UserTask{}
		utils.CopyStructFields(&taskInfo, v)
		resp.Data = append(resp.Data, &taskInfo)
	}
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp.Data})
}

// GetUserClaimTask
// @Summary		获取用户领取的任务
// @Description	获取用户领取的任务
// @Tags			任务
// @ID				GetUserClaimTask
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.GetUserClaimTaskListReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.GetUserClaimTaskListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/task/get_user_claim_task_list [post]
func GetUserClaimTaskList(c *gin.Context) {
	var (
		req   api.GetUserClaimTaskListReq
		resp  api.GetUserClaimTaskListResp
		reqPb pbTask.GetUserClaimTaskListReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	utils.CopyStructFields(&reqPb, req)
	// get user id from token
	var ok bool
	var errInfo string
	ok, reqPb.UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	// etcd im task
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// rpc
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.GetUserClaimTaskList(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetUserClaimTaskList failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	for _, v := range respPb.Data {
		taskInfo := api.UserTask{}
		utils.CopyStructFields(&taskInfo, v)
		resp.Data = append(resp.Data, &taskInfo)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", respPb)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp.Data})
}

// ClaimTaskRewards
// @Summary		领取任务奖励
// @Description	领取任务奖励
// @Tags			任务
// @ID				ClaimTaskRewards
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.ClaimTaskRewardsReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.ClaimTaskRewardsResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/task/claim_task_rewards [post]
func ClaimTaskRewards(c *gin.Context) {
	var (
		req   api.ClaimTaskRewardsReq
		reqPb pbTask.ClaimTaskRewardsReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	utils.CopyStructFields(&reqPb, req)
	// get user id from token
	var ok bool
	var errInfo string
	ok, reqPb.UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	// etcd im task
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// rpc
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.ClaimTaskRewards(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "Create Task failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", respPb)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg})
}

// DailyCheckIn
// @Summary		每日签到
// @Description	每日签到
// @Tags			任务
// @ID				DailyCheckIn
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.DailyCheckInReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.DailyCheckInResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/task/daily_check_in [post]
func DailyCheckIn(c *gin.Context) {
	var (
		req   api.DailyCheckInReq
		reqPb pbTask.DailyCheckInReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	utils.CopyStructFields(&reqPb, req)
	// get user id from token
	var ok bool
	var errInfo string
	ok, reqPb.UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	// etcd im task
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	// rpc
	client := pbTask.NewTaskServiceClient(etcdConn)
	respPb, err := client.DailyCheckIn(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "Create Task failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", respPb)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg})
}

// DailyIsCheckIn
// @Summary		今日是否签到
// @Description	今日是否签到
// @Tags			任务
// @ID				DailyIsCheckIn
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.DailyIsCheckInReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.DailyIsCheckInResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/task/daily_is_check_in [post]
func DailyIsCheckIn(c *gin.Context) {
	var (
		req api.DailyIsCheckInReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	// get user id from token
	var ok bool
	var errInfo string
	ok, userId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	// etcd im task
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImTask, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	isCheckIn, err := db.DB.GetUserIsCheckIn(userId)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CheckTaskIsFinished failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": isCheckIn})
}

// CheckIsHaveNftRecvID
// @Summary		检查是否有设置nft头像
// @Description	检查是否有设置nft头像
// @Tags			任务
// @ID				CheckIsHaveNftRecvID
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.CheckIsHaveNftRecvIDReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.CheckIsHaveNftRecvIDResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/task/check_is_have_nft_recvid [post]
func CheckIsHaveNftRecvID(c *gin.Context) {
	var (
		req    api.CheckIsHaveNftRecvIDReq
		userId string
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	// get user id from token
	var ok bool
	var errInfo string
	ok, userId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	// etcd im task
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImWeb3Js, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbWeb3.NewWeb3PubClient(etcdConn)
	respPb, err := client.CheckIsHaveNftRecvID(context.Background(), &pbWeb3.CheckIsHaveNftRecvIDReq{
		UserId: userId,
	})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CheckTaskIsFinished failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": respPb.HaveNft})
}

// CheckIsHaveOfficialNftRecvID
// @Summary		检查是否设置了官方nft头像
// @Description	检查是否设置了官方nft头像
// @Tags			任务
// @ID				CheckIsHaveOfficialNftRecvID
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.CheckIsHaveOfficialNftRecvIDReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.CheckIsHaveOfficialNftRecvIDResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/task/check_is_have_official_nft_recvid [post]
func CheckIsHaveOfficialNftRecvID(c *gin.Context) {
	var (
		req    api.CheckIsHaveOfficialNftRecvIDReq
		userId string
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	// get user id from token
	var ok bool
	var errInfo string
	ok, userId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	// etcd im task
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImWeb3Js, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbWeb3.NewWeb3PubClient(etcdConn)
	respPb, err := client.CheckIsHaveGuanFangNftRecvID(context.Background(), &pbWeb3.CheckIsHaveGuanFangNftRecvIDReq{
		UserId: userId,
	})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CheckTaskIsFinished failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": respPb.HaveNft})
}

// CheckIsFollowSystemTwitter
// @Summary		检查是否关注了官方twitter
// @Description	检查是否关注了官方twitter
// @Tags			任务
// @ID				CheckIsFollowSystemTwitter
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.CheckIsFollowSystemTwitterReq	true	"参数"
// @Produce		json
// @Success		0	{object}    api.CheckIsFollowSystemTwitterResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/task/check_is_follow_system_twitter [post]
func CheckIsFollowSystemTwitter(c *gin.Context) {
	var (
		req    api.CheckIsFollowSystemTwitterReq
		userId string
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	// get user id from token
	var ok bool
	var errInfo string
	ok, userId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	// etcd im task
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImWeb3Js, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbWeb3.NewWeb3PubClient(etcdConn)
	respPb, err := client.CheckIsFollowSystemTwitter(context.Background(), &pbWeb3.CheckUserIsFollowSystemTwitterReq{
		UserId: userId,
	})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CheckTaskIsFinished failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if respPb.CommonResp.ErrCode != 0 {
		// c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": false})
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": respPb.CommonResp.ErrMsg, "data": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": true})
}
