package friend

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/friend"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"

	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// @Summary		添加黑名单
// @Description	添加黑名单
// @Tags			好友相关
// @ID				AddBlack
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.AddBlacklistReq	true	"fromUserID为设置的用户 <br> toUserID为被设置的用户"
// @Produce		json
// @Success		0	{object}	api.AddBlacklistResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/add_black [post]
func AddBlack(c *gin.Context) {
	params := api.AddBlacklistReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.AddBlacklistReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params)
	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, "AddBlacklist args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.AddBlacklist(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddBlacklist failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add blacklist rpc server failed"})
		return
	}
	resp := api.AddBlacklistResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "AddBlacklist api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		批量加好友
// @Description	批量加好友
// @Tags			好友相关
// @ID				ImportFriend
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.ImportFriendReq	true	"fromUserID批量加好友的用户ID<br> friendUserIDList为"
// @Produce		json
// @Success		0	{object}	api.ImportFriendResp	"data列表中对象的result-1为添加该用户失败<br>0为成功"
// @Failure		500	{object}	api.Swagger500Resp		"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp		"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/import_friend [post]
func ImportFriend(c *gin.Context) {
	params := api.ImportFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.ImportFriendReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "ImportFriend args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.ImportFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "ImportFriend failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "ImportFriend failed "})
		return
	}
	resp := api.ImportFriendResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	if resp.ErrCode == 0 {
		for _, v := range RpcResp.UserIDResultList {
			resp.UserIDResultList = append(resp.UserIDResultList, api.UserIDResult{UserID: v.UserID, Result: v.Result})
		}
	}
	if len(resp.UserIDResultList) == 0 {
		resp.UserIDResultList = []api.UserIDResult{}
	}
	log.NewInfo(req.OperationID, "ImportFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		添加好友
// @Description	添加好友
// @Tags			好友相关
// @ID				AddFriend
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.AddFriendReq	true	"reqMsg为申请信息 <br> fromUserID为申请用户 <br> toUserID为被添加用户"
// @Produce		json
// @Success		0	{object}	api.AddFriendResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/add_friend [post]
func AddFriend(c *gin.Context) {
	params := api.AddFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.AddFriendReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	req.ReqMsg = params.ReqMsg

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "AddFriend args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.AddFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddFriend failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call AddFriend rpc server failed"})
		return
	}

	resp := api.AddFriendResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "AddFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// FollowAddFriend
// @Summary		追随好友或者取消追随
// @Description	追随好友或者取消追随
// @Tags			好友相关
// @ID				FollowAddFriend
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.FollowAddFriendReq	true	"reqMsg为申请信息 <br> fromUserID为申请用户 <br> toUserID为被添加用户"
// @Produce		json
// @Success		0	{object}	api.FollowAddFriendResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/follow_add_friend [post]
func FollowAddFriend(c *gin.Context) {
	params := api.FollowAddFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.FollowAddFriendReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	req.Follow = params.Follow

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "FollowAddFriend args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImFriendName,
		req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.FollowAddFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddFriend failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call AddFriend rpc server failed"})
		return
	}

	resp := api.FollowAddFriendResp{
		CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg},
	}
	for _, v := range RpcResp.PublicUserInfo {
		black := open_im_sdk.PublicUserInfo{}
		utils.CopyStructFields(&black, v)
		resp.PublicUserInfo = append(resp.PublicUserInfo, &black)
	}
	resp.Data = jsonData.JsonDataList(resp.PublicUserInfo)
	log.NewInfo(req.CommID.OperationID, "FollowAddFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// FollowFriendList
// @Summary		获取对方好友的关注或者是追谁者
// @Description	获取对方好友的关注或者是追谁者
// @Tags			好友相关
// @ID				FollowFriendList
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetFollowListReq	true	"参数"
// @Produce		json
// @Success		0	{object}	api.GetFollowListRsp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/following_list [post]
func FollowFriendList(c *gin.Context) {
	params := api.GetFollowListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetFriendFollowListReq{
		CommID:   &rpc.CommID{},
		IsFollow: params.IsFollow,
	}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	req.IsFollow = params.IsFollow

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "FollowFriendList args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImFriendName,
		req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.GetFriendFollowList(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "FollowFriendList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call AddFriend rpc server failed"})
		return
	}
	resp := api.GetFollowListRsp{
		CommResp: api.CommResp{
			ErrCode: RpcResp.CommonResp.ErrCode,
			ErrMsg:  RpcResp.CommonResp.ErrMsg,
		},
	}
	FollowListUserInfos := make([]api.FollowListUserInfo, len(RpcResp.PublicUserInfo))
	var wg sync.WaitGroup
	wg.Add(len(RpcResp.PublicUserInfo))
	for i, v := range RpcResp.PublicUserInfo {
		utils.CopyStructFields(&FollowListUserInfos[i], v)
		go func(index int, userId string) {
			defer wg.Done()
			gbroupInfo, err := rocksCache.GetSpaceInfoByUser(userId)
			if err == nil {
				grouInfo := open_im_sdk.GroupInfo{}
				utils.CopyStructFields(&grouInfo, gbroupInfo)
				FollowListUserInfos[index].Group = &grouInfo
			}
		}(i, v.UserID)
	}
	wg.Wait()
	resp.Data = FollowListUserInfos
	log.NewInfo(req.CommID.OperationID, "FollowFriendList api return ", resp)
	c.JSON(http.StatusOK, resp)

}

// @Summary		同意/拒绝好友请求
// @Description	同意/拒绝好友请求
// @Tags			好友相关
// @ID				AddFriendResponse
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.AddFriendResponseReq	true	"fromUserID同意/拒绝的用户ID<br>toUserID为申请用户D<br>handleMsg为处理信息<br>flag为具体操作, 1为同意, 2为拒绝"
// @Produce		json
// @Success		0	{object}	api.AddFriendResponseResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/add_friend_response [post]
func AddFriendResponse(c *gin.Context) {
	params := api.AddFriendResponseReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.AddFriendResponseReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	req.HandleMsg = params.HandleMsg
	req.HandleResult = params.Flag

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	utils.CopyStructFields(req, &params)
	log.NewInfo(req.CommID.OperationID, "AddFriendResponse args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.AddFriendResponse(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddFriendResponse failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add_friend_response rpc server failed"})
		return
	}

	resp := api.AddFriendResponseResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "AddFriendResponse api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		删除好友
// @Description	删除好友
// @Tags			好友相关
// @ID				DeleteFriend
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.DeleteFriendReq	true	"fromUserID为操作用户<br>toUserID为被删除用户"
// @Produce		json
// @Success		0	{object}	api.DeleteFriendResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/delete_friend [post]
func DeleteFriend(c *gin.Context) {
	params := api.DeleteFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.DeleteFriendReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "DeleteFriend args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.DeleteFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "DeleteFriend failed ", err, req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call delete_friend rpc server failed"})
		return
	}

	resp := api.DeleteFriendResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "DeleteFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		获取黑名单列表
// @Description	获取黑名单列表
// @Tags			好友相关
// @ID				GetBlacklist
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.GetBlackListReq	true	"fromUserID要获取黑名单的用户"
// @Produce		json
// @Success		0	{object}	api.GetBlackListResp{data=[]open_im_sdk.PublicUserInfo}
// @Failure		500	{object}	api.Swagger400Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger500Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/get_black_list [post]
func GetBlacklist(c *gin.Context) {
	params := api.GetBlackListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetBlacklistReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params)

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "GetBlacklist args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.GetBlacklist(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetBlacklist failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get blacklist rpc server failed"})
		return
	}

	resp := api.GetBlackListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	for _, v := range RpcResp.BlackUserInfoList {
		black := open_im_sdk.PublicUserInfo{}
		utils.CopyStructFields(&black, v)
		resp.BlackUserInfoList = append(resp.BlackUserInfoList, &black)
	}
	resp.Data = jsonData.JsonDataList(resp.BlackUserInfoList)
	log.NewInfo(req.CommID.OperationID, "GetBlacklist api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		设置好友备注
// @Description	设置好友备注
// @Tags			好友相关
// @ID				SetFriendRemark
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.SetFriendRemarkReq	true	"fromUserID为设置的用户<br> toUserID为被设置的用户<br> remark为好友备注"
// @Produce		json
// @Success		0	{object}	api.SetFriendRemarkResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/set_friend_remark [post]
func SetFriendRemark(c *gin.Context) {
	params := api.SetFriendRemarkReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.SetFriendRemarkReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	req.Remark = params.Remark

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "SetFriendComment args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.SetFriendRemark(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "SetFriendComment failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call set friend comment rpc server failed"})
		return
	}
	resp := api.SetFriendRemarkResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}

	log.NewInfo(req.CommID.OperationID, "SetFriendComment api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		把用户移除黑名单
// @Description	把用户移除黑名单
// @Tags			好友相关
// @ID				RemoveBlack
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.RemoveBlackListReq	true	"fromUserID要获取黑名单的用户"
// @Produce		json
// @Success		0	{object}	api.RemoveBlackListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/remove_black [post]
func RemoveBlack(c *gin.Context) {
	params := api.RemoveBlackListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.RemoveBlacklistReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "RemoveBlacklist args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.RemoveBlacklist(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "RemoveBlacklist failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call remove blacklist rpc server failed"})
		return
	}
	resp := api.RemoveBlackListResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "RemoveBlacklist api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		检查用户之间是否为好友
// @Description	检查用户之间是否为好友
// @Tags			好友相关
// @ID				IsFriend
// @Accept			json
// @Param			token	header	string			true	"im token"
// @Param			req		body	api.IsFriendReq	true	"fromUserID为请求用户<br> toUserID为要检查的用户"
// @Produce		json
// @Success		0	{object}	api.IsFriendResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/is_friend [post]
func IsFriend(c *gin.Context) {
	params := api.IsFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.IsFriendReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "IsFriend args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.IsFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "IsFriend failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add friend rpc server failed"})
		return
	}
	resp := api.IsFriendResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	resp.Response.Friend = RpcResp.Response

	log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		检查用户之间是否为互相关注
// @Description	检查用户之间是否为互相关注
// @Tags			好友相关
// @ID				IsFollowEachOtherFriend
// @Accept			json
// @Param			token	header	string			true	"im token"
// @Param			req		body	api.IsFriendReq	true	"fromUserID为请求用户<br> toUserID为要检查的用户"
// @Produce		json
// @Success		0	{object}	api.IsFriendResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/is_follow_eachother_friend [post]
func IsFollowEachOtherFriend(c *gin.Context) {
	params := api.IsFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.IsFriendReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "IsFollowEachOther args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.IsFollowEachOther(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "IsFollowEachOther failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add friend rpc server failed"})
		return
	}
	resp := api.IsFriendResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	resp.Response.Friend = RpcResp.Response

	log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		获取用户的好友列表
// @Description	获取用户的好友列表
// @Tags			好友相关
// @ID				GetFriendList
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetFriendListReq	true	"fromUserID为要获取好友列表的用户ID"
// @Produce		json
// @Success		0	{object}	api.GetFriendListResp{data=[]open_im_sdk.FriendInfo}
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/get_friend_list [post]
func GetFriendList(c *gin.Context) {
	params := api.GetFriendListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetFriendListReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params)

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "GetFriendList args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.GetFriendList(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetFriendList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get friend list rpc server failed"})
		return
	}

	resp := api.GetFriendListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, FriendInfoList: RpcResp.FriendInfoList}
	resp.Data = jsonData.JsonDataList(resp.FriendInfoList)
	log.NewInfo(req.CommID.OperationID, "GetFriendList api return ", resp)
	c.JSON(http.StatusOK, resp)
	//c.JSON(http.StatusOK, resp)
}

// @Summary		获取好友申请列表
// @Description	删除好友
// @Tags			好友相关
// @ID				GetFriendApplyList
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.GetFriendApplyListReq	true	"fromUserID为要获取申请列表的用户ID"
// @Produce		json
// @Success		0	{object}	api.GetFriendApplyListResp{data=[]open_im_sdk.FriendRequest}
// @Failure		500	{object}	api.Swagger400Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/get_friend_apply_list [post]
func GetFriendApplyList(c *gin.Context) {
	params := api.GetFriendApplyListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetFriendApplyListReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params)

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "GetFriendApplyList args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)

	RpcResp, err := client.GetFriendApplyList(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetFriendApplyList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get friend apply list rpc server failed"})
		return
	}

	resp := api.GetFriendApplyListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, FriendRequestList: RpcResp.FriendRequestList}
	resp.Data = jsonData.JsonDataList(resp.FriendRequestList)
	log.NewInfo(req.CommID.OperationID, "GetFriendApplyList api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary		获取自己的好友申请列表
// @Description	获取自己的好友申请列表
// @Tags			好友相关
// @ID				GetSelfFriendApplyList
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetSelfApplyListReq	true	"fromUserID为自己的用户ID"
// @Produce		json
// @Success		0	{object}	api.GetSelfApplyListResp{data=[]open_im_sdk.FriendRequest}
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/get_self_friend_apply_list [post]
func GetSelfFriendApplyList(c *gin.Context) {
	params := api.GetSelfApplyListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetSelfApplyListReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params)

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "GetSelfApplyList args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.GetSelfApplyList(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetSelfApplyList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get self apply list rpc server failed"})
		return
	}
	resp := api.GetSelfApplyListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, FriendRequestList: RpcResp.FriendRequestList}
	resp.Data = jsonData.JsonDataList(resp.FriendRequestList)
	log.NewInfo(req.CommID.OperationID, "GetSelfApplyList api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// GetUserFollowingList
// @Summary		获取用户的关注列表
// @Description	获取用户的关注列表
// @Tags			好友相关
// @ID				GetUserFollowingList
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetUserFollowingListReq	true	"参数"
// @Produce		json
// @Success		0	{object}	api.GetUserFollowingListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/get_user_following_list [post]
func GetUserFollowingList(c *gin.Context) {
	params := api.GetUserFollowingListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetUserFollowingListReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImFriendName,
		req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.GetUserFollowingList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserFollowingList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call GetUserFollowingList rpc server failed"})
		return
	}
	resp := api.GetFollowListRsp{
		CommResp: api.CommResp{
			ErrCode: RpcResp.CommonResp.ErrCode,
			ErrMsg:  RpcResp.CommonResp.ErrMsg,
		},
	}
	FollowListUserInfos := make([]api.FollowListUserInfo, len(RpcResp.PublicUserInfo))
	var wg sync.WaitGroup
	wg.Add(len(RpcResp.PublicUserInfo))
	for i, v := range RpcResp.PublicUserInfo {
		utils.CopyStructFields(&FollowListUserInfos[i], v)
		go func(index int, userId string) {
			defer wg.Done()
			gbroupInfo, err := rocksCache.GetSpaceInfoByUser(userId)
			if err == nil {
				grouInfo := open_im_sdk.GroupInfo{}
				utils.CopyStructFields(&grouInfo, gbroupInfo)
				FollowListUserInfos[index].Group = &grouInfo
			}
		}(i, v.UserID)
	}
	wg.Wait()
	resp.Data = FollowListUserInfos
	log.NewInfo(req.OperationID, "GetUserFollowingList api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// GetUserFollowedList
// @Summary		获取用户的粉丝列表
// @Description	获取用户的粉丝列表
// @Tags			好友相关
// @ID				GetUserFollowedList
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetUserFollowedListReq	true	"参数"
// @Produce		json
// @Success		0	{object}	api.GetUserFollowedListResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/friend/get_user_followed_list [post]
func GetUserFollowedList(c *gin.Context) {
	params := api.GetUserFollowedListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetUserFollowedListReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImFriendName,
		req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.GetUserFollowedList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserFollowedList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call GetUserFollowedList rpc server failed"})
		return
	}
	resp := api.GetFollowListRsp{
		CommResp: api.CommResp{
			ErrCode: RpcResp.CommonResp.ErrCode,
			ErrMsg:  RpcResp.CommonResp.ErrMsg,
		},
	}
	FollowListUserInfos := make([]api.FollowListUserInfo, len(RpcResp.PublicUserInfo))
	var wg sync.WaitGroup
	wg.Add(len(RpcResp.PublicUserInfo))
	for i, v := range RpcResp.PublicUserInfo {
		utils.CopyStructFields(&FollowListUserInfos[i], v)
		go func(index int, userId string) {
			defer wg.Done()
			gbroupInfo, err := rocksCache.GetSpaceInfoByUser(userId)
			if err == nil {
				grouInfo := open_im_sdk.GroupInfo{}
				utils.CopyStructFields(&grouInfo, gbroupInfo)
				FollowListUserInfos[index].Group = &grouInfo
			}
		}(i, v.UserID)
	}
	wg.Wait()
	resp.Data = FollowListUserInfos
	log.NewInfo(req.OperationID, "GetUserFollowedList api return ", resp)
	c.JSON(http.StatusOK, resp)
}
