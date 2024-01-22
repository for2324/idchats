package user

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpcfriend "Open_IM/pkg/proto/friend"
	rpcgroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	rpc "Open_IM/pkg/proto/user"
	"Open_IM/pkg/proto/web3pub"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// ShowUserProfile
// @Summary		针对个人信息分享页面分享
// @Description	针对个人信息分享页面分享
// @Tags		用户相关
// @ID			ShowUserProfile
// @Accept		json
// @Param		req		body	api.GetGlobalUserProfileReq	true	"123123"
// @Param		token	header	string						false	"im token"
// @Produce		json
// @Success		0	{object}	api.GetGlobalUserProfileResp	""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/user/share_user_profile [post]
func ShowUserProfile(c *gin.Context) {
	var (
		req  api.GetGlobalUserProfileReq
		resp api.GetGlobalUserProfileResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req.UserID = strings.ToLower(req.UserID)
	reqRpc := &rpc.GetUserInfoReq{}
	reqRpc.OperationID = req.OperationID
	_, reqRpc.OpUserID, _ = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)

	reqRpc.UserIDList = append(reqRpc.UserIDList, req.UserID)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfoWithoutToken(context.Background(), reqRpc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed1"})
		return
	}
	resultData := make(map[string]interface{}, 0)
	resultData["userInfo"] = RpcResp.UserInfoList
	resultData["userProfile"] = RpcResp.UserProfile
	resultData["twitter"] = RpcResp.Twitter
	resultData["dnsDomain"] = RpcResp.DnsDomain
	resultData["emailAddress"] = RpcResp.EmailAddress
	resultData["dnsDomainVerify"] = RpcResp.DnsDomainVerify
	resultData["userLinkTree"] = RpcResp.LinkTree
	if req.UserID != "" {
		groupInfo, err := rocksCache.GetSpaceInfoByUser(req.UserID)
		if err == nil {
			resultData["groupPath"] = fmt.Sprintf("/%s/%s", groupInfo.GroupID, groupInfo.GroupName)
			resultData["groupID"] = groupInfo.GroupID
		}
	}

	var reqgroupinfo rpcgroup.OperatorCommunityChannelRoleReq
	reqgroupinfo.OpUserID = req.UserID
	reqgroupinfo.OperationID = req.OperationID

	etcdConngroup := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName,
		req.OperationID)
	if etcdConngroup == nil {
		errMsg := reqgroupinfo.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqgroupinfo.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	clientgroup := rpcgroup.NewGroupClient(etcdConngroup)
	RpcRespTagTap, err2 := clientgroup.GetUserRoleTagInfo(context.Background(), &reqgroupinfo)
	if err2 != nil {
		if RpcRespTagTap != nil {
			resultData["userTag"] = RpcRespTagTap.RoleIpfs
		}
	}

	etcdConn2 := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImFriendName,
		req.OperationID)
	if etcdConn2 == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client2 := rpcfriend.NewFriendClient(etcdConn2)
	RpcResp2, err := client2.GetUserFollowedCount(context.Background(), &rpcfriend.GetUerFollowedCountReq{
		OperationID: req.OperationID,
		UserId:      req.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed2" + err.Error()})
		return
	}
	resultData["followedCount"] = RpcResp2.Count
	RpcResp3, err := client2.GetUserFollowingCount(context.Background(), &rpcfriend.GetUerFollowingCountReq{
		OperationID: req.OperationID,
		UserId:      req.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed2" + err.Error()})
		return
	}
	resultData["followingCount"] = RpcResp3.Count
	var ok bool
	var FromUserId string
	ok, FromUserId, _ = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if ok {
		RpcResp4, err := client2.IsFollowUser(context.Background(), &rpcfriend.IsFollowUserReq{
			OperationID: req.OperationID,
			FromUserId:  FromUserId,
			ToUserId:    req.UserID,
		})
		if err == nil {
			resultData["isFollow"] = RpcResp4.IsFollow
		}
	}

	// 临时性返回粉丝和关注列表
	reqrpc := &rpcfriend.GetFriendFollowListReq{
		CommID:   &rpcfriend.CommID{},
		IsFollow: true,
	}
	reqrpc.CommID.OpUserID = ""
	reqrpc.CommID.OperationID = req.OperationID
	reqrpc.CommID.ToUserID = req.UserID
	RpcResp4, err3 := client2.GetFriendFollowList(context.Background(), reqrpc)
	if err3 == nil {
		resultData["follow"] = RpcResp4.PublicUserInfo
	}
	reqrpc.IsFollow = false
	RpcResp5, err2 := client2.GetFriendFollowList(context.Background(), reqrpc)
	if err2 == nil {
		resultData["following"] = RpcResp5.PublicUserInfo
	}

	resp.UserInfoProfile = resultData

	c.JSON(http.StatusOK, resp)
}

// BindShowNft
// @Summary		绑定个性页面需要展示的nft
// @Description	绑定个性页面需要展示的nft
// @Tags			用户相关
// @ID				BindShowNft
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.BindShowNftReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.BindShowNftResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/bind_show_nft [post]
func BindShowNft(c *gin.Context) {
	var (
		req   api.BindShowNftReq
		resp  api.BindShowNftResp
		reqPb rpc.RPCBindShowNftReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.OperationID = req.OperationID
	reqPb.NftInfo = make([]*rpc.NftInfo, len(req.NftListShow))
	var waitSync sync.WaitGroup
	for key2, value := range req.NftListShow {
		key2 = key2
		tempValue := value
		if tempValue == nil {
			continue
		}
		waitSync.Add(1)
		go func(key int, tempValue *api.NftInfo, userid string) {
			defer waitSync.Done()
			reqPb.NftInfo[key] = &rpc.NftInfo{NftChainID: utils.StringToInt32(tempValue.NftChainID),
				NftContract:     tempValue.NftContract,
				TokenID:         tempValue.TokenID,
				NftContractType: tempValue.NftContractType,
			}
			PostCheckData := new(api.RequestImageTokenIdReq)
			PostCheckData.ContractAddress = tempValue.NftContract
			PostCheckData.TokenID = tempValue.TokenID
			PostCheckData.ChainID = tempValue.NftChainID
			PostCheckData.TokenImageModel = tempValue.NftContractType
			PostCheckDataByte, _ := json.Marshal(PostCheckData)
			resultByte, err := utils.HttpPost(config.Config.EnsPostCheck.Url+"/graph/tokenOwnerAddressContractChainID",
				"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, PostCheckDataByte)
			if err == nil {
				var resultData api.RequestTokenIdResp
				json.Unmarshal(resultByte, &resultData)
				if strings.EqualFold(resultData.TokenOwnerAddress, userid) {
					reqPb.NftInfo[key].NftTokenURL = resultData.TokenUrl
					fmt.Println("contract ：", tempValue.NftContract, " url is:", resultData.TokenUrl)
				}
			}

		}(key2, tempValue, reqPb.UserID)
	}
	waitSync.Wait()
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.BindShowNft(context.Background(), &reqPb)
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

// GetShowNft
// @Summary		展示NFT的内容
// @Description	展示NFT的内容
// @Tags			用户相关
// @ID				GetShowNft
// @Accept			json
// @Param			req		body	api.GetGlobalUserProfileReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.GetShowNftResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/get_show_nft [post]
func GetShowNft(c *gin.Context) {
	var (
		req   api.GetGlobalUserProfileReq
		resp  api.GetShowNftResp
		reqPb rpc.RPCBindShowNftReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.UserID = req.UserID
	_, reqPb.OpUserID, _ = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.GetBindShowNft(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg

	for _, value := range respPb.NftInfo {
		resp.NftListShow = append(resp.NftListShow, &api.NftInfo{
			ArticleID:       value.ID,
			TokenID:         value.TokenID,
			NftChainID:      utils.Int32ToString(value.NftChainID),
			NftContract:     value.NftContract,
			NftContractType: value.NftContractType,
			NftTokenURL:     value.NftTokenURL,
			LikesCount:      value.LikesCount,
			IsLikes:         value.IsLikes,
		})
	}
	c.JSON(http.StatusOK, resp)
	return
}

// GetUserSettingPage
// @Summary		V2.1个人设置页面展示
// @Description	V2.1个人设置页面展示
// @Tags			用户相关
// @ID				GetUserSettingPage
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.UserSettingPageReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UserSettingPageResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/get_user_setting_page [post]
func GetUserSettingPage(c *gin.Context) {
	var (
		req   api.UserSettingPageReq
		resp  api.UserSettingPageResp
		reqPb rpc.GetShowUserSettingReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req.UserID = strings.ToLower(req.UserID)
	reqPb.OperationID = req.OperationID
	reqPb.UserID = strings.ToLower(req.UserID)
	_, reqPb.OpUserID, _ = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.RpcUserSettingInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg

	for _, value := range respPb.ShowNftList {
		resp.NftListShow = append(resp.NftListShow, &api.NftInfo{
			ArticleID:       value.ID,
			TokenID:         value.TokenID,
			NftChainID:      utils.Int32ToString(value.NftChainID),
			NftContract:     value.NftContract,
			NftContractType: value.NftContractType,
			NftTokenURL:     value.NftTokenURL,
			LikesCount:      value.LikesCount,
			IsLikes:         value.IsLikes,
		})
	}
	resp.FaceURL = &respPb.FaceURL
	resp.Nickname = &respPb.Nickname
	resp.UserID = &respPb.UserID
	resp.UserIntroduction = &respPb.UserIntroduction
	resp.UserProfile = &respPb.UserProfile
	resp.UserHeadTokenInfo = new(api.NftInfo)
	if respPb.UserHeadTokenInfo != nil {
		utils.CopyStructFields(resp.UserHeadTokenInfo, respPb.UserHeadTokenInfo)
	}
	resp.EmailAddress = &respPb.EmailAddress
	resp.UserTwitter = &respPb.UserTwitter
	resp.DnsDomain = &respPb.DnsDomain
	resp.DnsDomainVerify = &respPb.DnsDomainVerify
	resp.IsShowTwitter = &respPb.IsShowTwitter
	resp.IsShowEmail = &respPb.ShowUserEmail
	resp.OpenAnnouncement = &respPb.OpenAnnouncement
	resp.FollowingCount = respPb.FollowingCount
	resp.FollowsCount = respPb.FollowsCount
	if len(respPb.LinkTree) > 0 {
		resp.LinkTree = new([]*open_im_sdk.LinkTreeMsgReq)
		*resp.LinkTree = make([]*open_im_sdk.LinkTreeMsgReq, len(respPb.LinkTree))
		for key, value := range respPb.LinkTree {
			(*resp.LinkTree)[key] = &open_im_sdk.LinkTreeMsgReq{
				LinkName:    value.LinkName,
				Link:        value.Link,
				FaceUrl:     value.FaceUrl,
				ShowStatus:  value.ShowStatus,
				DefaultIcon: value.DefaultIcon,
				Des:         value.Des,
				Bgc:         value.Bgc,
				DefaultUrl:  value.DefaultUrl,
				Type:        value.Type,
			}
		}
	}

	c.JSON(http.StatusOK, resp)
	return
}

// UpdateUserSettingPage
// @Summary		V2.1更新个人设置信息
// @Description	V2.1更新个人设置信息
// @Tags			用户相关
// @ID				UpdateUserSettingPage
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.UpdateUserSettingPageReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.UpdateUserSettingPageResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/update_user_setting_page [post]
func UpdateUserSettingPage(c *gin.Context) {
	var (
		req   api.UpdateUserSettingPageReq
		resp  api.UpdateUserSettingPageResp
		reqPb rpc.UpdateUserSettingReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	if req.Nickname != nil && *req.Nickname != "" {
		if !strings.HasSuffix(*req.Nickname, ".biu") {
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 502, "errMsg": "nickname only support ens by biu"})
			return
		}
		PostCheckDataByte, _ := json.Marshal(&api.EnsApiReq{
			Address:   reqPb.UserID,
			EnsDomain: *req.Nickname,
		})
		resultByte, err := utils.HttpPost(config.Config.EnsPostCheck.Url+"/graph/bindEnsCheck",
			"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, PostCheckDataByte)
		if err == nil {
			var resultData api.EnsApiResp
			err := json.Unmarshal(resultByte, &resultData)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"errCode": 501, "errMsg": "can't connect ens verifyserver"})
				return
			}
			if resultData.ErrCode != 0 {
				c.JSON(http.StatusBadRequest, gin.H{"errCode": 501, "errMsg": resultData.ErrMsg})
				return
			}
		}
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	if req.IsShowEmail != nil {
		if *req.IsShowEmail {
			reqPb.IsShowUserEmail = &wrapperspb.Int32Value{Value: 1}
		} else {
			reqPb.IsShowUserEmail = &wrapperspb.Int32Value{Value: 0}
		}
	}
	if req.IsShowTwitter != nil {
		if *req.IsShowTwitter {
			reqPb.IsShowTwitter = &wrapperspb.Int32Value{Value: 1}
		} else {
			reqPb.IsShowTwitter = &wrapperspb.Int32Value{Value: 0}
		}
	}
	if req.IsShowBalance != nil {
		if *req.IsShowBalance {
			reqPb.ShowBalance = &wrapperspb.Int32Value{Value: 1}
		} else {
			reqPb.ShowBalance = &wrapperspb.Int32Value{Value: 0}
		}
	}
	if req.OpenAnnouncement != nil {
		if *req.OpenAnnouncement > 0 {
			reqPb.OpenAnnouncement = &wrapperspb.Int32Value{Value: 1}
		} else {
			reqPb.OpenAnnouncement = &wrapperspb.Int32Value{Value: 0}
		}
	}
	if req.Nickname != nil {
		reqPb.Nickname = &wrapperspb.StringValue{Value: *req.Nickname}
	}
	if req.UserProfile != nil {
		reqPb.UserProfile = &wrapperspb.StringValue{Value: *req.UserProfile}
		if dataTemp := strings.Split(*req.UserProfile, "&&"); len(dataTemp) > 5 {
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "customer tag can't more than 5 count"})
			return
		}
	}
	if req.UserIntroduction != nil {
		reqPb.UserIntroduction = &wrapperspb.StringValue{Value: *req.UserIntroduction}
	}
	if req.DnsDomain != nil {
		reqPb.DnsDomain = &wrapperspb.StringValue{Value: *req.DnsDomain}
	}
	if req.UserHeadTokenInfo != nil {
		//检查nft是否是官方头像
		TokenUrl, boolFlag, errordata := CheckParamUserHeadInfo(req.OperationID, reqPb.UserID,
			req.UserHeadTokenInfo.NftContract, req.UserHeadTokenInfo.TokenID,
			req.UserHeadTokenInfo.NftChainID)
		if !boolFlag {
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "head nft not belong yourself," + errordata.Error()})
			return
		}
		if errordata == nil && TokenUrl == "" {
			TokenUrl = req.UserHeadTokenInfo.NftTokenURL
		}

		reqPb.UserHeadTokenInfo = new(rpc.NftInfo)
		reqPb.UserHeadTokenInfo.TokenID = req.UserHeadTokenInfo.TokenID
		reqPb.UserHeadTokenInfo.NftContract = req.UserHeadTokenInfo.NftContract
		if req.UserHeadTokenInfo.NftChainID == "" {
			reqPb.UserHeadTokenInfo.NftChainID = 0
		} else {
			reqPb.UserHeadTokenInfo.NftChainID = utils.StringToInt32(req.UserHeadTokenInfo.NftChainID)
		}
		reqPb.UserHeadTokenInfo.NftTokenURL = TokenUrl
		reqPb.FaceURL = &wrapperspb.StringValue{Value: TokenUrl}
	}
	if req.NftListShow != nil {
		reqPb.ShowNftList = make([]*rpc.NftInfo, len(req.NftListShow))
		var waitSync sync.WaitGroup
		for key2, value := range req.NftListShow {
			tempValue := value
			if tempValue == nil {
				continue
			}
			waitSync.Add(1)
			go func(key int, tempValue *api.NftInfo, userid string) {
				defer waitSync.Done()
				reqPb.ShowNftList[key] = &rpc.NftInfo{NftChainID: utils.StringToInt32(tempValue.NftChainID),
					NftContract:     tempValue.NftContract,
					TokenID:         tempValue.TokenID,
					NftContractType: tempValue.NftContractType,
				}
				PostCheckData := new(api.RequestImageTokenIdReq)
				PostCheckData.ContractAddress = tempValue.NftContract
				PostCheckData.TokenID = tempValue.TokenID
				PostCheckData.ChainID = tempValue.NftChainID
				PostCheckData.TokenImageModel = tempValue.NftContractType
				PostCheckData.OwnerAddress = userid
				PostCheckDataByte, _ := json.Marshal(PostCheckData)
				resultByte, err := utils.HttpPost(config.Config.EnsPostCheck.Url+"/graph/tokenOwnerAddressContractChainID",
					"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, PostCheckDataByte)
				if err == nil {
					var resultData api.RequestTokenIdResp
					fmt.Println("Post data :", tempValue.NftContract, ">> NftChainID", tempValue.NftChainID, "TokenID:>>",
						tempValue.TokenID, "返回结果:", string(resultByte))
					err := json.Unmarshal(resultByte, &resultData)
					if err != nil {
						return
					}
					if strings.EqualFold(resultData.TokenOwnerAddress, userid) {
						reqPb.ShowNftList[key].NftTokenURL = resultData.TokenUrl
						fmt.Println("contract ：", tempValue.NftContract, " url is:", resultData.TokenUrl)
					}
				}

			}(key2, tempValue, reqPb.UserID)
		}
		waitSync.Wait()
		reqPb.ShowNftListCount = &wrapperspb.Int32Value{Value: int32(len(reqPb.ShowNftList))}
	} else {
		reqPb.ShowNftListCount = nil
	}
	if req.LinkTree != nil {
		reqPb.LinkTree = make([]*open_im_sdk.LinkTreeMsgReq, 0)
		for _, value := range *req.LinkTree {
			tempValue := value
			//本来要判断是否是有效的url 地址
			//if utils.IsVaildUrl(tempValue.Link) {
			reqPb.LinkTree = append(reqPb.LinkTree, &open_im_sdk.LinkTreeMsgReq{
				LinkName:    tempValue.LinkName,
				Link:        tempValue.Link,
				FaceUrl:     tempValue.FaceUrl,
				ShowStatus:  tempValue.ShowStatus,
				UserID:      reqPb.UserID,
				DefaultIcon: tempValue.DefaultIcon,
				Des:         tempValue.Des,
				Bgc:         tempValue.Bgc,
				DefaultUrl:  tempValue.DefaultUrl,
				Type:        tempValue.Type,
			})
			//	}
		}
		reqPb.LinkTreeCount = &wrapperspb.Int32Value{Value: int32(len(reqPb.LinkTree))}
	} else {
		reqPb.LinkTreeCount = nil
	}
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.RpcUserSettingUpdate(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	if respPb.CommonResp.ErrCode == 0 {
		go func() {
			dbGroupInfo, err := rocksCache.GetSpaceInfoByUser(reqPb.UserID)
			if err == nil && dbGroupInfo != nil && reqPb.UserHeadTokenInfo != nil {
				etcdConngroup := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
					strings.Join(config.Config.Etcd.EtcdAddr, ","),
					config.Config.RpcRegisterName.OpenImGroupName,
					req.OperationID)
				if etcdConngroup == nil {
					errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
					log.NewError(req.OperationID, errMsg)
					c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
					return
				}
				clientgroup := rpcgroup.NewGroupClient(etcdConngroup)
				reqGroup := &rpcgroup.SetGroupInfoReq{GroupInfoForSet: &open_im_sdk.GroupInfoForSet{}}
				reqGroup.GroupInfoForSet.GroupID = dbGroupInfo.GroupID
				reqGroup.GroupInfoForSet.FaceURL = reqPb.UserHeadTokenInfo.NftTokenURL
				clientgroup.SetGroupInfo(context.Background(), reqGroup)
			}
		}()
	}
	c.JSON(http.StatusOK, resp)
}

// CheckDnsDomain
// @Summary		V2.1检查dnsdomain
// @Description	V2.1检查dnsdomain
// @Tags			用户相关
// @ID				CheckDnsDomain
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.UpdateUserSettingPageReq	true	"只能传递dnsdomain的值得"
// @Produce		json
// @Success		0	{object}	api.UpdateUserSettingPageResp	"返回0 即代表已经记录该域名并认证"
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/user/check_dns_domain [post]
func CheckDnsDomain(c *gin.Context) {
	var (
		req   api.UpdateUserSettingPageReq
		resp  api.UpdateUserSettingPageResp
		reqPB = new(web3pub.CheckDomainHadParseTxtReq)
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	if req.DnsDomain == nil || *req.DnsDomain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "dns domain error"})
		return
	}
	if _, err := url.Parse(*req.DnsDomain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "dns domain error"})
		return
	}

	var ok bool
	var errInfo string
	ok, reqPB.UserId, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPB.OperatorID = req.OperationID
	reqPB.DnsDomain = *req.DnsDomain

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImWeb3Js, reqPB.OperatorID)
	if etcdConn == nil {
		errMsg := reqPB.OperatorID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPB.OperatorID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := web3pub.NewWeb3PubClient(etcdConn)
	respPb, err := client.CheckDnsDomainHadParseBiuBiuTxt(context.Background(), reqPB)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	c.JSON(http.StatusOK, resp)
}
