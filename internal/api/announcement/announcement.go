package announcement

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbGroup "Open_IM/pkg/proto/group"
	pbChat "Open_IM/pkg/proto/msg"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	rpc "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

// SetGlobalAnnouncementMessageOpt
// @Summary		全局推送开关
// @Description	全局推送开关
// @Tags			广播相关
// @ID				SetGlobalAnnouncementMessageOpt
// @Accept			json
// @Param			token	header	string							true	"im token"
// @Param			req		body	api.SetGlobalRecvMessageOptReq	true	"globalRecvMsgOpt为接收全局推送设置0为关闭 1为开启"
// @Produce		json
// @Success		0	{object}	api.SetGlobalRecvMessageOptResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/announce/set_global_announce [post]
func SetGlobalAnnouncementMessageOpt(c *gin.Context) {
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
	log.NewInfo(params.OperationID, "SetGlobalAnnouncementMessageOpt args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	req.FieldName = "open_announcement"
	req.GlobalRecvMsgOpt = *params.OpenAnnouncement
	RpcResp, err := client.RpcUpdateUserFieldData(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "SetGlobalAnnouncementMessageOpt failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.UpdateUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "SetGlobalAnnouncementMessageOpt api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// CreateOrUpdateUserAnnouncementDraft
// @Summary		创建或者更新草稿箱
// @Description	创建或者更新草稿箱
// @Tags		广播相关
// @ID			创建或者更新草稿箱
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.AnnouncementArticleDraftReq	true	"用户信息简介"
// @Produce		json
// @Success		0	{object}	api.AnnouncementArticleDraftResp "更新的情况下ArticleDraftID 不为0 就是更新数据"
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/announce/createorupdate_user_announcement_draft [post]
func CreateOrUpdateUserAnnouncementDraft(c *gin.Context) {
	var (
		req  api.AnnouncementArticleDraftReq
		resp api.AnnouncementArticleDraftResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	if req.GroupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "group is null"})
		return
	}
	var ok bool
	var errInfo string
	ok, req.CreatorUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	cacheGroup, err := rocksCache.GetGroupInfoFromCache(req.GroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "this group not exist"})
		return
	}
	if cacheGroup.CreatorUserID != req.CreatorUserID {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "You are not this Group Creator"})
		return
	}
	var articleDraft db.AnnouncementArticleDraft
	utils.CopyStructFields(&articleDraft, req)
	err = imdb.CreateOrUpdateUserAnnounce(&articleDraft)
	if err != nil {
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteUserAnnouncementDraft
// @Summary		删除草稿箱
// @Description	删除草稿箱
// @Tags		广播相关
// @ID			删除草稿箱
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.AnnouncementArticleDraftReq	true	"除了id和op 我要 其他都不要"
// @Produce		json
// @Success		0	{object}	api.AnnouncementArticleDraftResp "AnnouncementArticleDraftResp 不为0 就是更新数据"
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/announce/delete_user_announcement_draft [post]
func DeleteUserAnnouncementDraft(c *gin.Context) {
	var (
		req  api.AnnouncementArticleDraftReq
		resp api.AnnouncementArticleDraftResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string

	ok, req.CreatorUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	var articleDraft db.AnnouncementArticleDraft
	utils.CopyStructFields(&articleDraft, req)
	err := imdb.DelUserAnnounceDraft(&articleDraft)
	if err != nil {
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
	}
	c.JSON(http.StatusOK, resp)
}

// GetUserAnnouncementDraft
// @Summary		获取草稿箱列表包含详情
// @Description	获取草稿箱列表包含详情
// @Tags		广播相关
// @ID			GetUserAnnouncementDraft
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.AnnouncementArticleDraftReq	true	"groupID op其他都不要"
// @Produce		json
// @Success		0	{object}	api.AnnouncementArticleDraftResp "AnnouncementArticleDraftResp 不为0 就是更新数据"
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/announce/get_user_announcement_draft [post]
func GetUserAnnouncementDraft(c *gin.Context) {
	var (
		req  api.AnnouncementArticleDraftReq
		resp api.AnnouncementArticleDraftResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string

	ok, req.CreatorUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	var err error
	resp.Data, err = imdb.GetTotalUserAnnounceDraft(req.CreatorUserID, req.GroupID)
	if err != nil {
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

// PublishUserAnnouncement
// @Summary		发布广播,包括转发推文
// @Description	发布广播，包括转发推文
// @Tags		广播相关
// @ID			PublishUserAnnouncement
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.CreateAnnouncementArticleReq	true	"除了id和op 我要 其他都不要"
// @Produce		json
// @Success		0	{object}	api.CreateAnnouncementArticleResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/announce/publish_user_announcement [post]
func PublishUserAnnouncement(c *gin.Context) {
	var (
		req  api.CreateAnnouncementArticleReq
		resp api.CreateAnnouncementArticleResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	if (req.ArticleID == nil && req.ArticleType == nil) && req.AnnouncementElem.AnnouncementUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "error AnnouncementUrl is null"})
		return
	}
	if req.IsGlobal == 1 && config.Config.SpaceArticle.PushUsdPrice > 0 {
		_, ok := config.Config.Pay.TnxTypeConfMap[req.TxnType]
		if !ok {
			log.NewError(req.OperationID, "config.Config.Pay.TnxTypeConfMap[req.TxnType] is empty ", req.TxnType)
			c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.ErrChainUp.ErrCode, "errMsg": "不支持的币种"})
			return
		}
	}
	reqPb := new(pbGroup.PublishAnnouncementReq)
	reqPb.AnnouncementElem = new(sdk_ws.AnnouncementMsg)
	utils.CopyStructFields(reqPb, &req)
	utils.CopyStructFields(reqPb.AnnouncementElem, &req.AnnouncementElem)
	reqPb.AnnouncementTitle = req.AnnouncementMsg.Title
	reqPb.AnnouncementSummary = req.AnnouncementMsg.Text
	reqPb.OperationID = req.OperationID
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	if (req.ArticleID == nil && req.ArticleType != nil) || (req.ArticleID != nil && req.ArticleType == nil) {
		c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": "文章的id并且文章的类型必须同时存在"})
		return
	}
	if req.ArticleID == nil && req.ArticleType == nil {
		reqPb.CreatorUserID = reqPb.OpUserID
	} else {
		reqPb.ArticleID = &wrapperspb.UInt64Value{Value: *req.ArticleID}
		reqPb.ArticleType = &wrapperspb.StringValue{Value: *req.ArticleType}
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.PublishAnnounceMoment(context.Background(), reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	c.JSON(http.StatusOK, resp)
	return
}

// DeletePublishUserAnnouncement
// @Summary		删除已经推送出去的广播
// @Description	删除已经推送出去的广播
// @Tags		广播相关
// @ID			DeletePublishUserAnnouncement
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.DeleteAnnouncementArticleReq	true	"除了id和op 我要 其他都不要"
// @Produce		json
// @Success		0	{object}	api.DeleteFromAnnouncementArticleResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/announce/delete_publish_user_announcement [post]
func DeletePublishUserAnnouncement(c *gin.Context) {
	var (
		req  api.DeleteAnnouncementArticleReq
		resp api.DeleteFromAnnouncementArticleResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	var userid string
	ok, userid, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok || req.ArticleID == "0" {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	dbinfodata, err := imdb.GetPublishGroupAnnounceByArticleID(req.ArticleID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	if dbinfodata.ArticleID != 0 && dbinfodata.CreatorUserID != userid {
		c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	err = imdb.DelUserAnnouncement(dbinfodata.ArticleID)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
	}
	return
}

// GetSpacePublishAnnouncementView
// @Summary		获取群空间内推送的广播
// @Description	获取群空间内推送的广播
// @Tags		广播相关
// @ID			GetSpacePublishAnnouncementView
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.GetAnnouncementArticleReq	true	"除了id和op 我要 其他都不要"
// @Produce		json
// @Success		0	{object}	api.GetAnnouncementArticleResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/announce/get_publish_space_announcement_list [post]
func GetSpacePublishAnnouncementView(c *gin.Context) {
	var (
		req   api.GetAnnouncementArticleReq
		resp  api.GetAnnouncementArticleResp
		reqPb pbGroup.GetPublishAnnouncementReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.PageIndex = req.PageIndex
	reqPb.PageSize = req.PageSize
	reqPb.IsGlobal = req.IsGlobal
	reqPb.ArticleID = utils.StringToInt64(req.ArticleID)
	if reqPb.ArticleID < 0 {
		reqPb.ArticleID = 0
	}

	var ok bool
	var errInfo string
	ok, reqPb.CreatorUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetPublishAnnounceMomentList(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	for _, value := range respPb.PublishAnnounce {
		dataTemp := new(api.AnnouncementArticleWithGroupInfo)
		utils.CopyStructFields(dataTemp, value)
		dataTemp.CreatedAt = utils.UnixSecondToTime(value.CreatedAt)
		dataTemp.UpdatedAt = utils.UnixSecondToTime(value.UpdatedAt)
		dataTemp.ArticleID = value.ArticleID
		dataTemp.GroupArticleID = value.GroupArticleID
		dataTemp.GroupName = value.GroupName
		dataTemp.FaceURL = value.FaceUrl
		dataTemp.IsRead = value.IsRead
		dataTemp.IsLikes = value.IsLikes
		resp.Data = append(resp.Data, dataTemp)
	}
	c.JSON(http.StatusOK, resp)
	return
}

// GetSpacePublishAnnouncementViewWithIdo
// @Summary		获取自己空间文章列表包括ido
// @Description	获取自己空间文章列表包括ido
// @Tags		广播相关
// @ID			GetSpacePublishAnnouncementViewWithIdo
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.GetAnnouncementArticleReq	true	"除了id和op 我要 其他都不要"
// @Produce		json
// @Success		0	{object}	api.GetAnnouncementArticleResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/announce/get_publish_space_announcement_ido_list [post]
func GetSpacePublishAnnouncementViewWithIdo(c *gin.Context) {
	var (
		req   api.GetAnnouncementArticleWithIdoReq
		resp  api.GetAnnouncementArticleWithIdoResp
		reqPb pbGroup.GetPublishAnnouncementWithIdoReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.PageIndex = req.PageIndex
	reqPb.PageSize = req.PageSize
	reqPb.IsGlobal = req.IsGlobal
	reqPb.ID = utils.StringToInt64(req.ArticleID)
	if reqPb.ID < 0 {
		reqPb.ID = 0
	}
	reqPb.ArticleType = req.ArticleType

	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetPublishAnnounceMomentListWithIdo(context.Background(), &reqPb)

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if len(respPb.GetPublishAnnounce()) > 0 {
		resp.ApiSpaceArticleListData = new(api.ApiSpaceArticleList)
		resp.ApiSpaceArticleListData.Data = make([]interface{}, len(respPb.PublishAnnounce))
		resp.ApiSpaceArticleListData.TotalCount = respPb.TotalCount
		resp.ApiSpaceArticleListData.CurrentPage = respPb.CurrentPage
		for keyIndex, value := range respPb.PublishAnnounce {
			if value.GenerateArticle.ArticleType == "announce" {
				resp.ApiSpaceArticleListData.Data[keyIndex] = value
			} else if value.GenerateArticle.ArticleType == "ido" {
				//查询数据 并
				jsonByteString := fmt.Sprintf(`{"groupID":"%s","idoID":"%d"}`,
					value.GenerateArticle.OpUserID, value.GenerateArticle.ArticleID)
				fmt.Println(">>>>>>>>>>><<<<<<<<<<<<<<<<", jsonByteString, "  ",
					config.Config.IdoPostCheckUrl+"/idoApi/getIDOByIdoID")
				resultByte, err := utils.HttpPost(config.Config.IdoPostCheckUrl+"/idoApi/getIDOByIdoID",
					"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, utils.String2bytes(jsonByteString))
				if err == nil {
					var TIdoStructData imdb.TIdoStruct
					TIdoStructData.Code = -1
					err = json.Unmarshal(resultByte, &TIdoStructData)
					if err == nil && TIdoStructData.Code == 0 {
						OutIdoStructData := new(IdoArticle)
						_ = utils.CopyStructFields(OutIdoStructData, &TIdoStructData.Data)
						OutIdoStructData.ArticleID = value.GenerateArticle.ArticleID
						OutIdoStructData.ArticleType = "ido"
						OutIdoStructData.IsPin = value.GenerateArticle.IsPin
						OutIdoStructData.CreatorInfo = value.CreatorInfo
						OutIdoStructData.OperatorInfo = value.OperatorInfo
						resp.ApiSpaceArticleListData.Data[keyIndex] = OutIdoStructData
					}
				}
			}
		}
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	c.JSON(http.StatusOK, resp)
	return
}

// GetSpacePersonalPublishAnnouncementViewWithIdo
// @Summary		获取自己收件箱的推文
// @Description	获取自己收件箱的推文
// @Tags		广播相关
// @ID			GetSpacePersonalPublishAnnouncementViewWithIdo
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.GetAnnouncementArticleWithIdoReq	true	"除了id和op 我要 其他都不要"
// @Produce		json
// @Success		0	{object}	api.GetAnnouncementArticleResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/announce/get_personal_publish_space_announcement_ido_list [post]
func GetSpacePersonalPublishAnnouncementViewWithIdo(c *gin.Context) {
	var (
		req   api.GetAnnouncementArticleWithIdoReq
		resp  api.GetAnnouncementArticleWithIdoResp
		reqPb pbGroup.GetPublishAnnouncementWithIdoReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.PageIndex = req.PageIndex
	reqPb.PageSize = req.PageSize
	reqPb.IsGlobal = req.IsGlobal
	reqPb.ID = utils.StringToInt64(req.ArticleID)
	if reqPb.ID < 0 {
		reqPb.ID = 0
	}
	reqPb.ArticleType = req.ArticleType

	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetPersonalPublishAnnounceMomentListWithIdo(context.Background(), &reqPb)

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	if len(respPb.GetPublishAnnounce()) > 0 {
		resp.ApiSpaceArticleListData = new(api.ApiSpaceArticleList)
		resp.ApiSpaceArticleListData.Data = make([]interface{}, len(respPb.PublishAnnounce))
		resp.ApiSpaceArticleListData.TotalCount = respPb.TotalCount
		resp.ApiSpaceArticleListData.CurrentPage = respPb.CurrentPage
		for keyIndex, value := range respPb.PublishAnnounce {
			if value.GenerateArticle.ArticleType == "announce" {
				resp.ApiSpaceArticleListData.Data[keyIndex] = value
			} else if value.GenerateArticle.ArticleType == "ido" {
				//查询数据 并
				jsonByteString := fmt.Sprintf(`{"groupID":"%s","idoID":"%d"}`,
					value.GenerateArticle.CreatorUserId, value.GenerateArticle.ArticleID)
				fmt.Println(">>>>>>>>>>><<<<<<<<<<<<<<<<", jsonByteString, "  ",
					config.Config.IdoPostCheckUrl+"/idoApi/getIDOByIdoID")
				resultByte, err := utils.HttpPost(config.Config.IdoPostCheckUrl+"/idoApi/getIDOByIdoID",
					"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, utils.String2bytes(jsonByteString))
				if err == nil {
					var TIdoStructData imdb.TIdoStruct
					TIdoStructData.Code = -1
					err = json.Unmarshal(resultByte, &TIdoStructData)
					if err == nil && TIdoStructData.Code == 0 {
						OutIdoStructData := new(IdoArticle)
						_ = utils.CopyStructFields(OutIdoStructData, &TIdoStructData.Data)
						OutIdoStructData.ArticleID = value.GenerateArticle.ArticleID
						OutIdoStructData.ArticleType = "ido"
						OutIdoStructData.CreatorInfo = value.CreatorInfo
						OutIdoStructData.OperatorInfo = value.OperatorInfo
						resp.ApiSpaceArticleListData.Data[keyIndex] = OutIdoStructData
					}
				}
			}
		}
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	c.JSON(http.StatusOK, resp)
	return
}

type IdoArticle struct {
	imdb.OutIdoStruct
	OperatorInfo *sdk_ws.PublicUserInfo `json:"operatorInfo"`
	CreatorInfo  *sdk_ws.PublicUserInfo `json:"creatorInfo"`
}

// LikeUnLikeAnnouncement
// @Summary		为推送点赞
// @Description	为推送点赞
// @Tags			广播相关
// @ID				LikeUnLikeAnnouncement
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.LikeActionNftReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.LikeActionNftResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/announce/like_unlike_announcement_article [post]
func LikeUnLikeAnnouncement(c *gin.Context) {
	var (
		req   api.LikeActionNftReq
		resp  api.LikeActionNftResp
		reqPb pbChat.SendLikeMsgReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.MsgData = new(sdk_ws.LikeRewordReq)
	utils.CopyStructFields(reqPb.MsgData, &req)
	reqPb.MsgData.ContentType = "announcement"
	var ok bool
	var errInfo string
	ok, reqPb.MsgData.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImMsgName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbChat.NewMsgClient(etcdConn)
	respPb, err := client.SendLikeAction(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.CommResp.ErrCode = respPb.ErrCode
	resp.CommResp.ErrMsg = respPb.ErrMsg
	c.JSON(http.StatusOK, resp)
	return
}

// ReadPublishAnnouncementView
// @Summary		读取文章详情
// @Description	读取文章详情
// @Tags			广播相关
// @ID				ReadPublishAnnouncementView
// @Accept			json
// @Param			token	header	string						true	"im token"
// @Param			req		body	api.AnnouncementReq	true	"请求体"
// @Produce		json
// @Success		0	{object}	api.BindShowNftResp
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router			/announce/read_publish_space_announcement_list [post]
func ReadPublishAnnouncementView(c *gin.Context) {
	var (
		req  api.AnnouncementReq
		resp = new(api.BindShowNftResp)
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	var UserID string
	ok, UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	var err error
	err = db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		var announcementArticle db.AnnouncementArticle
		err = tx.Table("announcement_article").Where("article_id=? and  status=0", req.ArticleID).
			First(&announcementArticle).Error
		if err != nil {
			return err
		}
		err = tx.Table("announcement_article_logs").Where("user_id =? and article_id=?", UserID, req.ArticleID).
			First(&announcementArticle).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = tx.Table("announcement_article_logs").Create(&db.AnnouncementArticleLog{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				ArticleID: announcementArticle.ArticleID,
				UserID:    UserID,
				IsLikes:   0,
				GroupID:   announcementArticle.GroupID,
				Status:    0,
				IsGlobal:  announcementArticle.IsGlobal,
			}).Error
			return err
		} else {
			return errors.New("已经读取过该文章")
		}
	})
	if err != nil {
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
	} else {
		resp.ErrCode = 0
		resp.ErrMsg = ""
	}

	c.JSON(http.StatusOK, resp)
	return
}

// DelSpacePublishAnnouncementViewWithIdo
// @Summary		删除空间推文列表包括ido不能删除
// @Description	删除空间推文列表包括ido不能删除
// @Tags		广播相关
// @ID			DelSpacePublishAnnouncementViewWithIdo
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param			req		body	api.OperatorSpaceArticleList	true	"除了id和op 我要 其他都不要"
// @Produce		json
// @Success		0	{object}	api.GetAnnouncementArticleResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/announce/del_publish_space_announcement_ido_list [post]
func DelSpacePublishAnnouncementViewWithIdo(c *gin.Context) {
	var (
		req   api.OperatorSpaceArticleList
		resp  api.CommResp
		reqPb pbGroup.DelSpaceArticleID
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.GroupID = req.GroupID
	reqPb.ID = utils.Int64ToString(req.ID)

	var ok bool
	var errInfo string
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.DelSpaceArticleList(context.Background(), &reqPb)

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.ErrCode = respPb.CommonResp.ErrCode
	resp.ErrMsg = respPb.CommonResp.ErrMsg
	c.JSON(http.StatusOK, resp)
	return
}

// PinSpacePublishAnnouncementViewWithIdo
// @Summary		置顶空间某个文章的列表
// @Description	置顶空间某个文章的列表
// @Tags		广播相关
// @ID			PinSpacePublishAnnouncementViewWithIdo
// @Accept		json
// @Param		token	header	string					true	"im token"
// @Param		req		body	api.OperatorSpaceArticleList	true	"ID为id,"
// @Produce		json
// @Success		0	{object}	api.CommResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/announce/pin_publish_space_announcement_ido_list [post]
func PinSpacePublishAnnouncementViewWithIdo(c *gin.Context) {
	var (
		req   api.OperatorSpaceArticleList
		resp  api.CommResp
		reqPb pbGroup.PinSpaceArticleID
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "Bind failed ", err.Error(), req)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.GroupID = req.GroupID
	reqPb.ID = utils.Int64ToString(req.ID)
	reqPb.IsPin = req.IsPin
	var ok bool
	var errInfo string
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.PinSpaceArticleList(context.Background(), &reqPb)

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.ErrCode = respPb.CommonResp.ErrCode
	resp.ErrMsg = respPb.CommonResp.ErrMsg
	c.JSON(http.StatusOK, resp)
	return
}

// GetHotSpaceBannerArticle
// @Summary		获取banner全局广播
// @Description	获取banner全局广播
// @Tags			空间相关
// @ID				GetHotSpaceBannerArticle
// @Accept			json
// @Param			token	header	string				true	"im token"
// @Param			req		body	api.GetHotSpaceReq	true	"123"
// @Produce		json
// @Success		0	{object}	api.GetAnnouncementArticleWithIdoResp ""
// @Failure		500	{object}	api.Swagger500Resp	"errCode为500 一般为服务器内部错误"
// @Failure		400	{object}	api.Swagger400Resp	"errCode为400 一般为参数输入错误, token未带上等"
// @Router		/api_space/get_hot_space_banner_article [post]
func GetHotSpaceBannerArticle(c *gin.Context) {
	params := api.GetHotSpaceReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbGroup.GetHotCommunityReq{
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
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetHotCommunityBannerAnnouncementList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CreateGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	var resp api.GetAnnouncementArticleWithIdoResp
	if len(respPb.GetPublishAnnounce()) > 0 {
		resp.ApiSpaceArticleListData = new(api.ApiSpaceArticleList)
		resp.ApiSpaceArticleListData.Data = make([]interface{}, len(respPb.PublishAnnounce))
		resp.ApiSpaceArticleListData.TotalCount = respPb.TotalCount
		resp.ApiSpaceArticleListData.CurrentPage = respPb.CurrentPage
		for keyIndex, value := range respPb.PublishAnnounce {
			if value.GenerateArticle.ArticleType == "announce" {
				resp.ApiSpaceArticleListData.Data[keyIndex] = value
			} else if value.GenerateArticle.ArticleType == "ido" {
				//查询数据 并
				jsonByteString := fmt.Sprintf(`{"groupID":"%s","idoID":"%d"}`,
					value.GenerateArticle.OpUserID, value.GenerateArticle.ArticleID)
				fmt.Println(">>>>>>>>>>><<<<<<<<<<<<<<<<", jsonByteString, "  ",
					config.Config.IdoPostCheckUrl+"/idoApi/getIDOByIdoID")
				resultByte, err := utils.HttpPost(config.Config.IdoPostCheckUrl+"/idoApi/getIDOByIdoID",
					"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, utils.String2bytes(jsonByteString))
				if err == nil {
					var TIdoStructData imdb.TIdoStruct
					TIdoStructData.Code = -1
					err = json.Unmarshal(resultByte, &TIdoStructData)
					if err == nil && TIdoStructData.Code == 0 {
						OutIdoStructData := new(IdoArticle)
						_ = utils.CopyStructFields(OutIdoStructData, &TIdoStructData.Data)
						OutIdoStructData.ArticleID = value.GenerateArticle.ArticleID
						OutIdoStructData.ArticleType = "ido"
						OutIdoStructData.IsPin = value.GenerateArticle.IsPin
						OutIdoStructData.CreatorInfo = value.CreatorInfo
						OutIdoStructData.OperatorInfo = value.OperatorInfo
						resp.ApiSpaceArticleListData.Data[keyIndex] = OutIdoStructData
					}
				}
			}
		}
	}
	resp.CommResp.ErrCode = respPb.CommonResp.ErrCode
	resp.CommResp.ErrMsg = respPb.CommonResp.ErrMsg
	log.NewInfo(req.OperationID, "GetHotCommunityBanner api return ", resp)
	c.JSON(http.StatusOK, resp)
}
