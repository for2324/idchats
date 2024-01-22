package group

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	// rpcGroup "Open_IM/internal/rpc/group"
)

// CreateCommunity
// @Summary		创建社区
// @Description	创建社区
// @Tags			社区相关
// @ID				CreateCommunity
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.CreateCommunityReq	true	"groupIDList为群ID列表"
// @Produce		json
// @Success		0	{object}	api.CreateCommunityResp{data=open_im_sdk.GroupInfo} ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/group/create_community [post]
func CreateCommunity(c *gin.Context) {
	params := api.CreateCommunityReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.CreateCommunityReq{GroupInfo: &open_im_sdk.GroupInfo{}}
	utils.CopyStructFields(req.GroupInfo, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "CreateCommunity failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	req.OwnerUserID = params.OwnerUserID
	req.OperationID = params.OperationID
	req.GroupInfo.GroupType = constant.WorkingGroup
	if params.OwnerUserID == "" {
		req.OwnerUserID = req.OpUserID
	}
	log.NewInfo(req.OperationID, "CreateCommunity args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.CreateCommunity(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CreateCommunity failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	resp := api.CreateCommunityResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	if RpcResp.ErrCode == 0 {
		utils.CopyStructFields(&resp.GroupInfo, RpcResp.GroupInfo)
		resp.Data = jsonData.JsonDataOne(&resp.GroupInfo)
	}
	log.NewInfo(req.OperationID, "CreateCommunity api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// CreateCommunityChannel
// @Summary		创建社区频道
// @Description	创建社区频道
// @Tags			社区相关
// @ID				CreateCommunityChannel
// @Accept			json
// @Param			token	header	string		true	"im token"
// @Param			req		body	api.UpdateCommunityChannelReq	true	"频道类型 opinfo:add,del,update"
// @Produce		json
// @Success		0	{object}	api.UpdateCommunityChannelResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/group/create_community_channel [post]
func CreateCommunityChannel(c *gin.Context) {
	params := api.UpdateCommunityChannelReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.CommunityChannelReq{
		OperationID: params.OperationID,
		OpInfo:      params.OpInfo,
		ChannelInfo: &open_im_sdk.GroupChannelInfo{
			GroupID: params.GroupID,
		},
	}
	switch params.OpInfo {
	case "add":
		if strings.Trim(params.ChannelName, " ") == "" {
			c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": "pls set channelname"})
			return
		}
		req.ChannelInfo.ChannelName = params.ChannelName
	case "del":
		if params.ChannelId == "" || utils.StringToInt(params.ChannelId) <= 10 {
			c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": "pls set ChannelId"})
			return
		}
		req.ChannelInfo.ChannelID = params.ChannelId
	case "update":
		if params.ChannelId == "" || utils.StringToInt(params.ChannelId) <= 10 || strings.Trim(params.ChannelName, " ") == "" {
			c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": "pls set ChannelId"})
			return
		}
		req.ChannelInfo.ChannelID = params.ChannelId
		req.ChannelInfo.ChannelName = params.ChannelName
	}

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "CreateCommunityChannel failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	req.OwnerUserID = req.OpUserID
	req.OperationID = params.OperationID
	log.NewInfo(req.OperationID, "CreateGroup args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.UpdateCommunityChannel(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CreateGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	resp := api.UpdateCommunityChannelResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode,
		ErrMsg: RpcResp.CommonResp.ErrMsg}}
	if RpcResp.CommonResp.ErrCode == 0 {
		utils.CopyStructFields(&resp.GroupChannelInfo, RpcResp.ChannelList)
		resp.Data = jsonData.JsonDataOne(&resp.GroupChannelInfo)
	}
	log.NewInfo(req.OperationID, "CreateGroup api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// GetHotSpaceBanner
// @Summary		获取banner
// @Description	获取banner
// @Tags			空间相关
// @ID				GetHotSpaceBanner
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.GetHotSpaceReq	true	"123"
// @Produce		json
// @Success		0	{object}	api.GetHotSpaceResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/api_space/get_hot_space_banner [post]
func GetHotSpaceBanner(c *gin.Context) {
	params := api.GetHotSpaceReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetHotCommunityReq{
		OperationID: params.OperationID,
		OpUserID:    "",
	}
	_, req.OpUserID, _ = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	req.OperationID = params.OperationID

	log.NewInfo(req.OperationID, "GetHotCommunity args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetHotCommunityBanner(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CreateGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	resp := api.GetHotSpaceMainInfoResp{
		CommResp: api.CommResp{ErrCode: RpcResp.ErrCode,
			ErrMsg: RpcResp.ErrMsg}}
	if RpcResp.ErrCode == 0 {
		resp.BannerArrayImage = make([]*api.BannerImage, 0)
		for _, value := range RpcResp.HotBannerInfo {
			resp.BannerArrayImage = append(resp.BannerArrayImage, &api.BannerImage{
				BannerImage: value.BannerImage,
				BannerSort:  value.BannerSort,
				BannerUrl:   value.BannerUrl,
			})
		}
	}
	log.NewInfo(req.OperationID, "GetHotCommunityBanner api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// GetHotSpace
// @Summary		获取热门空间
// @Description	获取热门空间
// @Tags			空间相关
// @ID				GetHotSpace
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.GetHotSpaceReq	true	"123"
// @Produce		json
// @Success		0	{object}	api.GetHotSpaceResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/api_space/get_hot_space [post]
func GetHotSpace(c *gin.Context) {
	params := api.GetHotSpaceReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetHotSpaceReq{
		OperationID: params.OperationID,
		PageIndex:   params.PageIndex,
		PageSize:    params.PageSize,
	}
	var ok bool
	var errInfo string
	ok, req.UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetHotSpace failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	var RpcResp *rpc.GetHotSpaceResp
	var err error
	if params.SearchType == "follow" {
		RpcResp, err = client.GetMyFollowingSpace(context.Background(), req)
	} else {
		RpcResp, err = client.GetHotSpace(context.Background(), req)
		// s := rpcGroup.NewGroupServer(0)
		// RpcResp, err := s.GetHotSpace(context.TODO(), req)

	}
	if err != nil {
		log.NewError(req.OperationID, "CreateGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	if RpcResp.CommonResp.ErrCode != 0 {
		c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.CommonResp.ErrCode, "errMsg": RpcResp.CommonResp.ErrMsg})
		return
	}

	log.NewInfo(req.OperationID, "GetHotSpace api return ", RpcResp.UserInfoList)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": RpcResp.UserInfoList})
}

// SearchCommunity
// @Summary		查找社区
// @Description	查找社区
// @Tags			社区相关
// @ID				SearchCommunity
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.GetHotCommunityReq	true	"123123"
// @Produce		json
// @Success		0	{object}	api.SearchCommunityResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/group/search_community [post]
func SearchCommunity(c *gin.Context) {
	params := api.GetHotCommunityReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.SearchCommunityReq{
		OperationID: params.OperationID,
		OpUserID:    "",
	}
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetHotCommunity failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	req.OperationID = params.OperationID
	req.OwnerUserID = req.OpUserID
	req.SearchName = params.SearchTitle

	log.NewInfo(req.OperationID, "GetHotCommunity args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.SearchCommunity(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CreateGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	resp := api.SearchCommunityResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	if RpcResp.ErrCode == 0 {
		resp.GroupInfoArray = make([]*api.GroupInfoHotMessage, 0)
		for _, value := range RpcResp.GroupInfo {
			tempdata := new(api.GroupInfoHotMessage)
			utils.CopyStructFields(tempdata, value)
			resp.GroupInfoArray = append(resp.GroupInfoArray, tempdata)
		}
	}
	log.NewInfo(req.OperationID, "GetHotCommunity api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// JoinCommunity
// @Summary		加入社区
// @Description	加入社区
// @Tags			社区相关
// @ID				JoinCommunity
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.JoinGroupReq	true	"reqMessage为申请进群信息<br>groupID为申请的群ID"
// @Produce		json
// @Success		0	{object}	api.JoinGroupResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/group/join_community [post]
func JoinCommunity(c *gin.Context) {
	params := api.JoinGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.JoinGroupReq{}
	utils.CopyStructFields(req, params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "JoinGroup args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)

	RpcResp, err := client.JoinGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "JoinGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}
	log.NewInfo(req.OperationID, "JoinGroup api return", RpcResp.String())
	c.JSON(http.StatusOK, resp)
}

// CommunityChannel
// @Summary		社区频道列表
// @Description	社区频道列表
// @Tags			社区相关
// @ID				CommunityChannel
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.CommunityChannelAllListReq	true	"群号"
// @Produce		json
// @Success		0	{object}	api.CommunityChannelAllListResp	""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/group/get_community_channel_list [post]
func CommunityChannel(c *gin.Context) {
	params := api.CommunityChannelAllListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.CommunityChannelAllListReq{}
	utils.CopyStructFields(req, params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "CommunityChannel args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)

	RpcResp, err := client.GetCommunityAllChannel(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CommunityChannel failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := &api.CommunityChannelAllListResp{
		CommResp: api.CommResp{
			ErrCode: RpcResp.CommonResp.ErrCode,
			ErrMsg:  RpcResp.CommonResp.ErrMsg,
		},
		GroupChannelInfoList: RpcResp.ChannelList,
		Data:                 nil,
	}
	resp.Data = jsonData.JsonDataList(RpcResp.ChannelList)
	log.NewInfo(req.OperationID, "CommunityChannel api return ", len(RpcResp.ChannelList))
	c.JSON(http.StatusOK, resp)
}

// CommunityChannelStatus
// @Summary		社区频道列表状态
// @Description	社区频道列表状态
// @Tags			社区相关
// @ID				CommunityChannelStatus
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.CommunityChannelStatusReq	true	"群号"
// @Produce		json
// @Success		0	{object}	api.CommunityChannelStatusResp	""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/group/group_channel_status [post]
func CommunityChannelStatus(c *gin.Context) {
	params := api.CommunityChannelStatusReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqpb := &rpc.CommunityChannelInfoReq{}
	utils.CopyStructFields(reqpb, params)

	var ok bool
	var errInfo string
	ok, reqpb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), reqpb.OperationID)
	if !ok {
		errMsg := reqpb.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(reqpb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(reqpb.OperationID, "CommunityChannelStatus args ", reqpb.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, reqpb.OperationID)
	if etcdConn == nil {
		errMsg := reqpb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqpb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)

	RpcResp, err := client.GetCommunityChannelByGroupIDAndChannelID(context.Background(), reqpb)
	if err != nil {
		log.NewError(reqpb.OperationID, "CommunityChannelStatus failed ", err.Error(), reqpb.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	resp := &api.CommunityChannelStatusResp{
		CommResp: api.CommResp{
			ErrCode: RpcResp.CommonResp.ErrCode,
			ErrMsg:  RpcResp.CommonResp.ErrMsg,
		},
	}
	if RpcResp.CommonResp.ErrCode != 0 {
		resp.Data = constant.GroupStatusDismissed
	} else {
		resp.Data = RpcResp.ChannelList.ChannelStatus
	}
	log.NewInfo(params.OperationID, "CommunityChannelStatus api return ", resp.Data)
	c.JSON(http.StatusOK, resp)
}

// GetUserJoinedGroupList
// @Summary		获取指定某个用户加入群列表
// @Description	获取指定某个用户加入群列表
// @Tags			群组相关
// @ID				GetUserJoinedGroupList
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.GetJoinedGroupListReq	true	"fromUserID为要获取的用户ID"
// @Produce		json
// @Success		0	{object}	api.GetJoinedGroupListResp{data=[]open_im_sdk.GroupInfo}
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/group/get_user_joined_group_list [post]
func GetUserJoinedGroupList(c *gin.Context) {
	params := api.GetJoinedGroupListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetJoinedGroupListReq{}
	utils.CopyStructFields(req, params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "GetJoinedGroupList args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetUserJoinedGroupList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetJoinedGroupList failed  ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	GroupListResp := api.GetJoinedGroupListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, GroupInfoList: RpcResp.GroupList}
	GroupListResp.Data = jsonData.JsonDataList(GroupListResp.GroupInfoList)
	log.NewInfo(req.OperationID, "GetJoinedGroupList api return ", GroupListResp)
	c.JSON(http.StatusOK, GroupListResp)
}

// GetGroupAllHistoryMessageList
// @Summary		获取群内聊天信息的历史消息
// @Description	获取群内聊天信息的历史消息
// @Tags			群组相关
// @ID				GetGroupAllHistoryMessageList
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.GetHistoryMessageListParams	true	"fromUserID为要获取的用户ID"
// @Produce		json
// @Success		0	{object}	api.GetHistoryMessageListResp{data=[]api.MsgStruct}
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/group/get_group_all_history_message_list [post]
func GetGroupAllHistoryMessageList(c *gin.Context) {
	params := api.GetHistoryMessageListParams{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetHistoryMessageListParamsReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	var errInfo string
	ok, _, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, "GetGroupAllHistoryMessageList args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetGroupHistoryMessageList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupAllHistoryMessageList failed  ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	groupHistoryMessageListResp := api.GetHistoryMessageListResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode,

		ErrMsg: RpcResp.CommonResp.ErrMsg}}
	groupHistoryMessageListResp.Data = make([]*api.MsgStruct, 0)
	for _, value := range RpcResp.Message {
		apimsgstruct := new(api.MsgStruct)
		utils.CopyStructFields(apimsgstruct, value)
		msgHandleByContentType(apimsgstruct)
		groupHistoryMessageListResp.Data = append(groupHistoryMessageListResp.Data, apimsgstruct)
	}

	log.NewInfo(req.OperationID, "GetGroupAllHistoryMessageList api return ", len(groupHistoryMessageListResp.Data))
	c.JSON(http.StatusOK, groupHistoryMessageListResp)
}

func msgHandleByContentType(msg *api.MsgStruct) (err error) {
	_ = utils.JsonStringToStruct(msg.AttachedInfo, &msg.AttachedInfoElem)
	switch msg.ContentType {
	case constant.Picture:
		msg.PictureElem = new(api.PictureElem)
		err = utils.JsonStringToStruct(msg.Content, msg.PictureElem)
	case constant.Voice:
		msg.SoundElem = new(api.SoundElem)
		err = utils.JsonStringToStruct(msg.Content, msg.SoundElem)
	case constant.Video:
		msg.VideoElem = new(api.VideoElem)
		err = utils.JsonStringToStruct(msg.Content, msg.VideoElem)
	case constant.File:
		msg.FileElem = new(api.FileElem)
		err = utils.JsonStringToStruct(msg.Content, msg.FileElem)
	case constant.AdvancedText:
		msg.MessageEntityElem = new(api.MessageEntityElem)
		err = utils.JsonStringToStruct(msg.Content, msg.MessageEntityElem)
	case constant.Location:
		msg.LocationElem = new(api.LocationElem)
		err = utils.JsonStringToStruct(msg.Content, &msg.LocationElem)
	case constant.Custom:
		msg.CustomElem = new(api.CustomElem)
		err = utils.JsonStringToStruct(msg.Content, &msg.CustomElem)
	case constant.Quote:
		msg.QuoteElem = new(api.QuoteElem)
		err = utils.JsonStringToStruct(msg.Content, msg.QuoteElem)
	case constant.Merger:
		msg.MergeElem = new(api.MergeElem)
		err = utils.JsonStringToStruct(msg.Content, &msg.MergeElem)
	}

	return utils.Wrap(err, "")
}

// CreateRoleTag
// @Summary		创建某个个社区下某个标签
// @Description	创建某个个社区下某个标签
// @Tags			群组相关
// @ID				CreateRoleTag
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.CreateRoleTagReq	true	"fromUserID为要获取的用户ID"
// @Produce		json
// @Success		0	{object}	api.CreateRoleTagResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/group/create_role_tag [post]
func CreateRoleTag(c *gin.Context) {
	params := api.CreateRoleTagReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.CreateCommunityRoleReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	var errInfo string
	//如果用户的token 不存在 那么要怎么操作： 从线上拉去记录来做？
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, "CreateCommunityRole args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.CreateCommunityRole(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CreateCommunityRole failed  ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	GroupListResp := api.CreateRoleTagResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode,
		ErrMsg: RpcResp.CommonResp.ErrMsg}}
	if RpcResp.CommonResp.ErrCode == 0 {
		GroupListResp.GroupTagString = RpcResp.ReBackOrderID
	}
	log.NewInfo(req.OperationID, "CreateRoleTagResp api return ", GroupListResp)
	c.JSON(http.StatusOK, GroupListResp)
}

// GetCommunityRoleTag
// @Summary		查看某个社区下有多少标签
// @Description	查看某个社区下有多少标签
// @Tags			群组相关
// @ID				GetCommunityRoleTag
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetCommunityRoleTagReq	true	"groupid 为查询的groupid"
// @Produce		json
// @Success		0	{object}	api.GetCommunityRoleTagResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/group/get_community_role_tag [post]
func GetCommunityRoleTag(c *gin.Context) {
	params := api.GetCommunityRoleTagReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.GetCommunityRoleReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	var errInfo string
	//如果用户的token 不存在 那么要怎么操作： 从线上拉去记录来做？
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, "CreateCommunityRole args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetCommunityRole(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CreateCommunityRole failed  ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	GroupListResp := api.GetCommunityRoleTagResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode,
		ErrMsg: RpcResp.CommonResp.ErrMsg}}
	if RpcResp.CommonResp.ErrCode == 0 {
		for _, value := range RpcResp.CommunityRoleInfo {
			tempdata := new(api.CommitRoleTagReq)
			utils.CopyStructFields(tempdata, value)
			GroupListResp.CommitRoleTagReq = append(GroupListResp.CommitRoleTagReq, tempdata)
		}
	}
	log.NewInfo(req.OperationID, "GetCommunityRoleTag api return ", GroupListResp)
	c.JSON(http.StatusOK, GroupListResp)
}

// GetCommunityRoleTagDetail
// @Summary		查看某个标签下拥有标签的人
// @Description	查看某个标签下拥有标签的人
// @Tags			群组相关
// @ID				GetCommunityRoleTagDetail
// @Accept			json
// @Param			token	header	string					true	"im token"
// @Param			req		body	api.GetCommunityRoleTagReq	true	"groupid 为查询的groupid"
// @Produce		json
// @Success		0	{object}	api.GetCommunityRoleTagResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/group/get_community_role_tag_detail [post]
func GetCommunityRoleTagDetail(c *gin.Context) {
	params := api.GetCommunityRoleTagReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	if params.GroupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "SpaceID cant't empty"})
		return
	}
	req := &rpc.GetCommunityRoleReqDetail{}
	utils.CopyStructFields(req, params)
	var ok bool
	var errInfo string
	//如果用户的token 不存在 那么要怎么操作： 从线上拉去记录来做？
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, "CreateCommunityRole args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetCommunityRoleDetail(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CreateCommunityRole failed  ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	GroupListResp := api.GetCommunityRoleTagDetailResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode,
		ErrMsg: RpcResp.CommonResp.ErrMsg}}
	if RpcResp.CommonResp.ErrCode == 0 {
		for _, value := range RpcResp.CommunityRoleInfoDetail {
			tempdata := new(api.CommunityRoleUserInfoList)
			utils.CopyStructFields(tempdata, value)
			GroupListResp.MemberList = append(GroupListResp.MemberList, tempdata)
		}
	}
	log.NewInfo(req.OperationID, "GetCommunityRoleTag api return ", GroupListResp)
	c.JSON(http.StatusOK, GroupListResp)
}
